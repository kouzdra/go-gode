package ide

import (
	"log"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	//"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
	gsci "github.com/kouzdra/go-scintilla/gtk"
	"github.com/kouzdra/go-analyzer/project"
	"os"
	"fmt"
	"path"
)

var _ = log.Printf

type IDE struct {
	Window    *gtk .Window
	Editors   *Editors
	Menubar   *gtk .MenuBar
	StatusBar *gtk .Statusbar
	RED        gsci.Style
	Prj       *project.Project
	View      *gtk.TreeView
	Store     *gtk.TreeStore
}

func NewIDE () *IDE {
	ide := &IDE{}
	
	ide.Window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	ide.Window.SetTitle("Editor")
	ide.Window.Connect("destroy", gtk.MainQuit)

	ide.MakeMenu()

	vbox := gtk.NewVBox(false, 1)
	vbox.PackStart(ide.Menubar, false, false, 0)

	hpaned := gtk.NewHPaned()
	vbox.Add(hpaned)

	ide.MakeTree ()
	swinT := gtk.NewScrolledWindow(nil, nil)
	swinT.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swinT.SetShadowType(gtk.SHADOW_NONE)
	swinT.Add (ide.View)
	hpaned.Add1 (swinT)
	ide.Editors = NewEditors (ide)
	hpaned.Add2 (ide.Editors.Notebook)
	hpaned.SetPosition (256)

	ide.StatusBar = gtk.NewStatusbar()
	context_id := ide.StatusBar.GetContextId("go-gode")
	ide.StatusBar.Push(context_id, "Go Dev.Env.")

	vbox.PackStart(ide.StatusBar, false, false, 0)

	ide.Window.Add(vbox)
	ide.Window.SetSizeRequest(1024, 800)
	ide.Window.ShowAll()

	return ide
}

func  (ide *IDE) fillTree (top bool, dirs [] project.Dir, iter *gtk.TreeIter) {
	model := ide.View.GetModel ()
	//model.Object.Ref ()
	ide.View.SetModel (nil)
	
	defer func () {
		ide.View.SetModel (model)
		//model.Object.Unref ()
	} ()
	subs := make ([]func (), 0, len(dirs))
	for _, dir := range dirs {
		var subIter gtk.TreeIter
		ide.Store.Append(&subIter, iter)
		name := dir.Path
		if !top { name = path.Base (name) }
		ide.Store.Set(&subIter, gtk.NewImage().RenderIcon(gtk.STOCK_DIRECTORY,
			gtk.ICON_SIZE_SMALL_TOOLBAR, "").GPixbuf, name)
		subs = append (subs, func () { ide.fillTree (false, dir.Sub, &subIter) } )
		if pkg := ide.Prj.Pkgs [dir.Path]; pkg != nil {
			for sPath, _ := range pkg.Srcs {
				var srcIter gtk.TreeIter
				ide.Store.Append(&srcIter, &subIter)
				ide.Store.Set(&srcIter, gtk.NewImage().RenderIcon(gtk.STOCK_FILE,
					gtk.ICON_SIZE_SMALL_TOOLBAR, "").GPixbuf, sPath)				
			}
		}
	}
	for _, sub := range subs {
		sub ()
	}
}

func (ide *IDE) LoadView () {
 	log.Printf ("Tree view loading ...")
	ide.fillTree (true, ide.Prj.Tree, nil)
 	log.Printf ("... OK")
}

func (ide *IDE) MakeTree () {
	ide.Store = gtk.NewTreeStore(gdkpixbuf.GetType(), glib.G_TYPE_STRING)
	ide.View  = gtk.NewTreeView()
	model := ide.Store.ToTreeModel()
	ide.View.SetModel(model)
	ide.View.AppendColumn(gtk.NewTreeViewColumnWithAttributes("pixbuf", gtk.NewCellRendererPixbuf(), "pixbuf", 0))
	ide.View.AppendColumn(gtk.NewTreeViewColumnWithAttributes("text"  , gtk.NewCellRendererText  (), "text"  , 1))
	ide.View.SetHeadersVisible (true)

	ide.View.Connect("row_activated", func() {
		var path *gtk.TreePath
		var column *gtk.TreeViewColumn
		ide.View.GetCursor(&path, &column)
		mes := "TreePath is: " + path.String()
		dialog := gtk.NewMessageDialog(
			ide.View.GetTopLevelAsWindow(),
			gtk.DIALOG_MODAL,
			gtk.MESSAGE_INFO,
			gtk.BUTTONS_OK,
			mes)
		dialog.SetTitle("TreePath")
		dialog.Response(func() {
			dialog.Destroy()
		})
		dialog.Run()
	})
}

	
func (ide *IDE) LoadProject () {
 	log.Printf ("Project loading ...")
	ide.Prj = project.NewProject ()
	ide.Prj.SetRoot (os.ExpandEnv("$GOROOT"))
	ide.Prj.SetPath (os.ExpandEnv("$GOPATH"))
	ide.Prj.Load ()
	log.Printf ("Project loaded: #Dirs: %d", len (ide.Prj.Dirs))
}

func (ide *IDE) MakeMenu () {

	ide.Menubar = gtk.NewMenuBar ()
	addCascade := func (label string, fill func (*gtk.Menu)) {
		cascademenu := gtk.NewMenuItemWithMnemonic(label)
		ide.Menubar.Append(cascademenu)
		submenu := gtk.NewMenu()
		cascademenu.SetSubmenu(submenu)
		fill (submenu)
	}
	makeItem := func (label string, action func ()) *gtk.MenuItem {
		menuitem := gtk.NewMenuItemWithMnemonic(label)
		menuitem.Connect("activate", action)
		return menuitem
	}
	
	addCascade ("_File", func (submenu *gtk.Menu) {
		submenu.Append(makeItem ("E_xit", gtk.MainQuit))
		//submenu.Append(makeItem ("_Open", ide.Editor.LoadFileFromDialog))
	})

	addCascade ("_Edit", func (submenu *gtk.Menu) {
		submenu.Append(makeItem ("_Complete", ide.Complete))
		//submenu.Append(makeItem ("_Open", ide.Editor.LoadFileFromDialog))
	})

	addCascade ("_View", func (submenu *gtk.Menu) {
		submenu.Append(makeItem ("_Font", func () {
			fsd := gtk.NewFontSelectionDialog("Font")
			fsd.Response(func() {
				fmt.Println(fsd.GetFontName())
				fsd.Destroy()
			})
			fsd.SetTransientFor(ide.Window)
			fsd.Run()
		}))
	})

	addCascade ("_Help", func (submenu *gtk.Menu) {
		submenu.Append(makeItem ("_About", func () {
			dialog := gtk.NewAboutDialog()
			dialog.SetName("Go-Gtk Demo!")
			dialog.SetProgramName("demo")
			dialog.SetAuthors([]string{"Kouzdra"})
			//dir, _ := filepath.Split(os.Args[0])
			//imagefile := filepath.Join(dir, "../../data/mattn-logo.png")
			//pixbuf, _ := gdkpixbuf.NewPixbufFromFile(imagefile)
			//dialog.SetLogo(pixbuf)
			dialog.SetLicense("The library is available under the same terms and conditions as the Go, the BSD style license, and the LGPL (Lesser GNU Public License). The idea is that if you can use Go (and Gtk) in a project, you should also be able to use go-gtk.")
			dialog.SetWrapLicense(true)
			dialog.Run()
			dialog.Destroy()
		}))
	})
}


//-------------------------------------------------

func (ide *IDE) Complete () {
	if page := ide.Editors.GetCurrent (); page != nil {
		log.Printf ("Complete [%s]\n", page.Editor.FName)
		pos := page.Editor.Sci.GetCurrentPos ()
		if src := page.Editor.Src; src == nil {
			log.Printf ("Complete [%s] at %d, no SRC found\n", page.Editor.FName, pos)
		} else {
			log.Printf ("Complete [%s|%s::%s] at %d\n", page.Editor.FName, src.Dir, src.Name, pos)
			if compl := ide.Prj.Complete (src, int (pos)); compl == nil {
				log.Printf ("    -- No completion context found\n")
			} else {
				log.Printf("  [%s/%s] (%d/%d) #%d\n", compl.Pref, compl.Name, compl.Pos, compl.End, len (compl.Choices))
				for i, c := range compl.Choices {
					log.Printf ("    %d) %s(%s) AKA [%s] pos=%d end=%d\n", i, c.Kind, c.Name, c.Full, c.Pos, c.End)
				}
			}
		}
	}
}

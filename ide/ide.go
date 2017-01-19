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
	"sort"
)

var _ = log.Printf

type IDE struct {
	Window    *gtk .Window
	Editor    *Editor
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

	ide.Editor = NewEditor(ide)
	ide.MakeMenu()

	vbox := gtk.NewVBox(false, 1)
	vbox.PackStart(ide.Menubar, false, false, 0)

	hpaned := gtk.NewHPaned()
	vbox.Add(hpaned)

	swinE := gtk.NewScrolledWindow(nil, nil)
	swinE.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swinE.SetShadowType(gtk.SHADOW_IN)

	ide.MakeTree ()
	swinT := gtk.NewScrolledWindow(nil, nil)
	swinT.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swinT.SetShadowType(gtk.SHADOW_NONE)
	notebook := gtk.NewNotebook()
	hpaned.Add1 (swinT)
	hpaned.Add2 (notebook)
	hpaned.SetPosition (256)

	swinT.Add (ide.View)
	swinE.Add(ide.Editor.Sci)

	notebook.AppendPage(swinE, gtk.NewLabel("fileName"))
	
	ide.StatusBar = gtk.NewStatusbar()
	context_id := ide.StatusBar.GetContextId("go-gode")
	ide.StatusBar.Push(context_id, "Go Dev.Env.")

	vbox.PackStart(ide.StatusBar, false, false, 0)

	
	ide.Window.Add(vbox)
	ide.Window.SetSizeRequest(1024, 800)
	ide.Window.ShowAll()

	return ide
}

type Pkgs []*project.Pkg
func (p Pkgs) Len  () int { return len (p) }
func (p Pkgs) Swap (i, j int) { p [i], p [j] = p [j], p [i] }
func (p Pkgs) Less (i, j int) bool { return p [i].Dir < p [j].Dir }

func (ide *IDE) LoadView () {
	var iter0 gtk.TreeIter
	ide.Store.Append(&iter0, nil)
	ide.Store.Set(&iter0, gtk.NewImage().RenderIcon(gtk.STOCK_FLOPPY,
		gtk.ICON_SIZE_SMALL_TOOLBAR, "").GPixbuf, "GO.PATH")

	pkgs := make ([]*project.Pkg, 0, len (ide.Prj.Pkgs))
	for _, pkg := range ide.Prj.Pkgs {
		pkgs = append (pkgs, pkg)
	}
	sort.Sort (Pkgs (pkgs))

	for _, pkg := range pkgs {
		//log.Printf ("path=%s #srcs=%d\n", pkg.Dir + "+" + pkg.Name, len (pkg.Srcs))
		var iter1 gtk.TreeIter
		ide.Store.Append(&iter1, &iter0)
		ide.Store.Set(&iter1, gtk.NewImage().RenderIcon(gtk.STOCK_DIRECTORY,
			gtk.ICON_SIZE_SMALL_TOOLBAR, "").GPixbuf, pkg.Dir)
		for sPath, _ := range pkg.Srcs {
			//log.Printf ("   src=%s\n", sPath)
			var iter2 gtk.TreeIter
			ide.Store.Append(&iter2, &iter1)
			ide.Store.Set(&iter2, gtk.NewImage().RenderIcon(gtk.STOCK_FILE,
				gtk.ICON_SIZE_SMALL_TOOLBAR, "").GPixbuf, sPath)
		}
	}
	ide.Store.Append(&iter0, nil)
	ide.Store.Set(&iter0, gtk.NewImage().RenderIcon(gtk.STOCK_FLOPPY,
		gtk.ICON_SIZE_SMALL_TOOLBAR, "").GPixbuf, "GO.ROOT")
}

func (ide *IDE) MakeTree () {
	ide.Store = gtk.NewTreeStore(gdkpixbuf.GetType(), glib.G_TYPE_STRING)
	ide.View  = gtk.NewTreeView()
	model := ide.Store.ToTreeModel()
	ide.View.SetModel(model)
	ide.View.AppendColumn(gtk.NewTreeViewColumnWithAttributes("pixbuf", gtk.NewCellRendererPixbuf(), "pixbuf", 0))
	ide.View.AppendColumn(gtk.NewTreeViewColumnWithAttributes("text"  , gtk.NewCellRendererText  (), "text"  , 1))
	ide.View.SetHeadersVisible (false)

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
		submenu.Append(makeItem ("_Open", ide.Editor.LoadFileFromDialog))
	})

	addCascade ("_View", func (submenu *gtk.Menu) {
		submenu.Append(makeItem ("_Font", func () {
			fsd := gtk.NewFontSelectionDialog("Font")
			fsd.Response(func() {
				fmt.Println(fsd.GetFontName())
				ide.Editor.Sci.Styling.SetFont(ide.RED, fsd.GetFontName())
				ide.Editor.Sci.Styling.SetUnderline (ide.RED, true);
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

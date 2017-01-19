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
)

var _ = log.Printf

type IDE struct {
	Window  *gtk .Window
	Editor  *Editor
	Menubar *gtk .MenuBar
	RED      gsci.Style
	Prj     *project.Project
	View    *gtk.TreeView
	Store   *gtk.TreeStore
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

	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetShadowType(gtk.SHADOW_IN)

	ide.MakeTree ()
	hpaned.Add1 (ide.View)
	hpaned.Add2 (swin)
	hpaned.SetPosition (256)

	swin.Add(ide.Editor.Sci)

	ide.Window.Add(vbox)
	ide.Window.SetSizeRequest(1024, 800)
	ide.Window.ShowAll()

	return ide
}

func (ide *IDE) LoadView () {
	var iter0 gtk.TreeIter
	ide.Store.Append(&iter0, nil)
	ide.Store.Set(&iter0, gtk.NewImage().RenderIcon(gtk.STOCK_FLOPPY,
		gtk.ICON_SIZE_SMALL_TOOLBAR, "").GPixbuf, "GO.PATH")
	for path, pkg := range ide.Prj.Pkgs {
		log.Printf ("path=%s #srcs=%d\n", path, len (pkg.Srcs))
		
		var iter1 gtk.TreeIter
		ide.Store.Append(&iter1, &iter0)
		ide.Store.Set(&iter1, gtk.NewImage().RenderIcon(gtk.STOCK_DIRECTORY,
			gtk.ICON_SIZE_SMALL_TOOLBAR, "").GPixbuf, path)
		for sPath, _ := range pkg.Srcs {
			log.Printf ("   src=%s\n", sPath)
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

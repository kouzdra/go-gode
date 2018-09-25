package ide

import "log"
import "strings"
import "path/filepath"
import "github.com/mattn/go-gtk/gdkpixbuf"
import "github.com/mattn/go-gtk/glib"
import "github.com/mattn/go-gtk/gtk"
import gsci "github.com/kouzdra/go-scintilla/gtk"
import project "github.com/kouzdra/go-analyzer/gproject"
import "github.com/kouzdra/go-gode/options"
import "github.com/kouzdra/go-gode/icons"

var _ = log.Printf

type IDE struct {
	Window      *gtk.Window
	Editors     *Editors
	ErrorsView  *gtk.TreeView
	ErrorsStore *gtk.TreeStore
	Menubar     *gtk.MenuBar
	StatusBar   *gtk.Statusbar
	RED         gsci.Style
	Project     project.Project
	View        *gtk.TreeView
	Store       *gtk.TreeStore
	Accel       *gtk.AccelGroup
	Options     *options.Options
	Icons       *icons.Icons
}

func NewIDE () *IDE {
	ide := &IDE{}

	ide.Options = options.New ()
	ide.Icons = icons.New ()
	
	ide.Window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	ide.Window.SetTitle("Editor")
	ide.Window.Connect("destroy", gtk.MainQuit)

	ide.Accel = gtk.NewAccelGroup ()
	ide.Window.AddAccelGroup (ide.Accel)
	
	ide.MakeMenu()

	vbox := gtk.NewVBox(false, 1)
	vbox.PackStart(ide.Menubar, false, false, 0)

	hpaned := gtk.NewHPaned()
	vbox.Add(hpaned)

	ide.MakeTree ()
	swinT := gtk.NewScrolledWindow(nil, nil)
	swinT.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swinT.SetShadowType(gtk.SHADOW_NONE)
	swinT.AddWithViewPort (ide.View)
	hpaned.Add1 (swinT)

	vpaned := gtk.NewVPaned()
	ide.Editors = NewEditors (ide)
	vpaned.Add1 (ide.Editors.Notebook)
	ide.MakeErrorsList ()
	vpaned.Add1 (ide.ErrorsView)

	hpaned.Add2 (vpaned)
	hpaned.SetPosition (256)

	ide.StatusBar = gtk.NewStatusbar()
	context_id := ide.StatusBar.GetContextId("go-gode")
	ide.StatusBar.Push(context_id, "Go Dev.Env.")

	vbox.PackStart(ide.StatusBar, false, false, 0)

	ide.Window.Add(vbox)
	ide.Window.SetSizeRequest(1000, 600)
	ide.Window.SetResizable(true)
	ide.Window.Maximize()
	ide.Window.ShowAll()

	return ide
}


func (ide *IDE) MakeErrorsList () {
	ide.ErrorsStore = gtk.NewTreeStore(gdkpixbuf.GetType(), glib.G_TYPE_STRING, glib.G_TYPE_STRING)
	ide.ErrorsView  = gtk.NewTreeView()
	model := ide.Store.ToTreeModel()
	ide.ErrorsView.SetModel(model)
	ide.ErrorsView.AppendColumn(gtk.NewTreeViewColumnWithAttributes("pixbuf", gtk.NewCellRendererPixbuf(), "pixbuf", COL_ICON))
	ide.ErrorsView.AppendColumn(gtk.NewTreeViewColumnWithAttributes("text"  , gtk.NewCellRendererText  (), "text"  , COL_FNAME))
	msgCol := gtk.NewTreeViewColumnWithAttributes("text"  , gtk.NewCellRendererText  (), "text"  , COL_FPATH)
	ide.View.AppendColumn(msgCol)
	ide.View.SetHeadersVisible (false)
}

func (ide *IDE) ReadablePath (fname string) (string, bool) {
	path := filepath.Clean(filepath.Dir(fname))
	goroot := filepath.Clean(ide.Project.GetRoot().Name)
	gopath := filepath.Clean(ide.Project.GetPath().Name)
	//log.Printf("1: path=%s, ROOT=%s PATH=%s", path, goroot, gopath)
	if strings.HasPrefix (path, goroot) {
		return filepath.Join ("$GOROOT", strings.TrimPrefix (path, goroot)), true
	}
	if strings.HasPrefix (path, gopath) {
		return filepath.Join ("$GOPATH", strings.TrimPrefix (path, gopath)), false
	}
	return filepath.Clean (path), false
}

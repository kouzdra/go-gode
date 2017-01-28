package ide

import "log"
import "github.com/mattn/go-gtk/gtk"
import gsci "github.com/kouzdra/go-scintilla/gtk"
import "github.com/kouzdra/go-analyzer/project"
import "github.com/kouzdra/go-gode/options"

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
	Accel     *gtk.AccelGroup
	Options   *options.Options
}

func NewIDE () *IDE {
	ide := &IDE{}

	ide.Options = options.New ()
	
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

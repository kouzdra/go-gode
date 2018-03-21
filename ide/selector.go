package ide

import "log"
import "github.com/mattn/go-gtk/gtk"
//import "github.com/mattn/go-gtk/gdk"
////import gsci "github.com/kouzdra/go-scintilla/gtk"
//import "github.com/kouzdra/go-analyzer/project"
//import "github.com/kouzdra/go-gode/options"
//import "github.com/kouzdra/go-gode/icons"

var _ = log.Printf

type Selector struct {
	ide       *IDE
	Dialog    *gtk.Dialog
	Entry     *gtk.Label
	/*Editors   *Editors
	Menubar   *gtk .MenuBar
	StatusBar *gtk .Statusbar
	RED        gsci.Style
	Prj       *project.Project
	View      *gtk.TreeView
	Store     *gtk.TreeStore*/
	Accel     *gtk.AccelGroup
	/*Options   *options.Options
	Icons     *icons.Icons*/
}

func (ide *IDE) NewSelector () *Selector {
	selector := &Selector{}

	selector.ide = ide
	selector.Dialog = gtk.NewDialog()
	//ide.Window.SetTitle("Editor")
	selector.Dialog.Connect("close", func () { } )
	selector.Dialog.Response(func() {
		log.Println("Selector closed")
		selector.Dialog.Destroy()
	})

	//selector.Accel = gtk.NewAccelGroup ()
	//selector.Window.AddAccelGroup (selector.Accel)
	//selector.Window.AddAccelerator ("destroy", selector.Accel, gdk.KEY_Escape, 0, 0)
	vbox := gtk.NewVBox(false, 1)
	//selector.Entry = gtk.NewLabel()
	//selector.Entry.SetText("Hello-2")
	vbox.Add(gtk.NewLabel("Hell-2"))
	vbox.Add(gtk.NewLabel("Hell-2"))
	//vbox.PackStart(ide.Menubar, false, false, 0)

	//hpaned := gtk.NewHPaned()
	//vbox.Add(hpaned)

	//selector.Window.MakeTree ()
	//win := gtk.NewWindow(nil, nil)
	//winT.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	//winT.SetShadowType(gtk.SHADOW_NONE)
	//winT.AddWithViewPort (ide.View)
	//hpaned.Add1 (swinT)
	//ide.Editors = NewEditors (ide)
	//hpaned.Add2 (ide.Editors.Notebook)
	//hpaned.SetPosition (256)

	//ide.StatusBar = gtk.NewStatusbar()
	//context_id := ide.StatusBar.GetContextId("go-gode")
	//ide.StatusBar.Push(context_id, "Go Dev.Env.")

	//vbox.PackStart(ide.StatusBar, false, false, 0)*/

	selector.Dialog.Add(vbox)
	//selector.Window.SetModal(true)
	selector.Dialog.SetDecorated(false)
	selector.Dialog.SetSizeRequest(400, 300)
	selector.Dialog.Run()

	return selector
}


func (ide *IDE) Select () {
	ide.NewSelector ()	
}

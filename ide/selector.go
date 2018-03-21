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
	Entry     *gtk.Entry
	View      *gtk.TreeView
	Store     *gtk.TreeStore
	Accel     *gtk.AccelGroup
}

func (ide *IDE) NewSelector () *Selector {
	selector := &Selector{}

	selector.ide = ide
	selector.Dialog = gtk.NewDialog()
	//selector.Dialog.Connect("close", func () { } )
	selector.Dialog.Response(func() {
		log.Println("Selector closed")
		selector.Dialog.Destroy()
	})


	//selector.Accel = gtk.NewAccelGroup ()
	//selector.Window.AddAccelGroup (selector.Accel)
	//selector.Window.AddAccelerator ("destroy", selector.Accel, gdk.KEY_Escape, 0, 0)
	vbox := selector.Dialog.GetContentArea()
	selector.Entry = gtk.NewEntry()
	selector.Entry.SetText("prefix")
	vbox.PackStart(selector.Entry, false, false, 0)
	
	l2 := gtk.NewLabel("Hell-2")
	vbox.PackStart(l2, true, true, 0)
	//l2.Show()

	//vbox.PackStart(ide.Menubar, false, false, 0)

	//hpaned := gtk.NewHPaned()
	//vbox.Add(hpaned)

	selector.ide.MakeTree ()
	swinT := gtk.NewScrolledWindow(nil, nil)
	swinT.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swinT.SetShadowType(gtk.SHADOW_NONE)
	swinT.AddWithViewPort (ide.View)
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

	//selector.Dialog.Add(vbox)
	//selector.Window.SetModal(true)
	selector.Dialog.SetDecorated(false)
	selector.Dialog.SetSizeRequest(400, 300)
	selector.Dialog.ShowAll()
	selector.Dialog.Run()

	return selector
}


func (ide *IDE) Select () {
	ide.NewSelector ()	
}

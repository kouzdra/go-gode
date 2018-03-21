package ide

import "log"
import "github.com/mattn/go-gtk/gtk"
import "github.com/mattn/go-gtk/gdkpixbuf"
import 	"github.com/mattn/go-gtk/glib"

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


	vbox := selector.Dialog.GetContentArea()
	selector.Entry = gtk.NewEntry()
	selector.Entry.SetText("prefix")
	vbox.PackStart(selector.Entry, false, false, 0)
	
	selector.Store = gtk.NewTreeStore(gdkpixbuf.GetType(), glib.G_TYPE_STRING, glib.G_TYPE_STRING)
	selector.View  = gtk.NewTreeView()
	model := selector.Store.ToTreeModel()
	selector.View.SetModel(model)
	selector.View.AppendColumn(gtk.NewTreeViewColumnWithAttributes("pixbuf", gtk.NewCellRendererPixbuf(), "pixbuf", COL_ICON))
	selector.View.AppendColumn(gtk.NewTreeViewColumnWithAttributes("text"  , gtk.NewCellRendererText  (), "text"  , COL_FNAME))
	pathCol := gtk.NewTreeViewColumnWithAttributes("text"  , gtk.NewCellRendererText  (), "text"  , COL_FPATH)
	pathCol.SetVisible (false)
	selector.View.AppendColumn(pathCol)
	selector.View.SetHeadersVisible (false)

	selector.View.Connect("row_activated", func() {
	})

	vbox.PackStart(selector.View, true, true, 0)

	selector.ide.MakeTree ()
	swinT := gtk.NewScrolledWindow(nil, nil)
	swinT.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swinT.SetShadowType(gtk.SHADOW_NONE)
	swinT.AddWithViewPort (selector.View)

	selector.Dialog.SetDecorated(false)
	selector.Dialog.SetSizeRequest(400, 300)
	selector.Dialog.ShowAll()
	selector.Dialog.Run()

	return selector
}


func (ide *IDE) Select () {
	ide.NewSelector ()	
}

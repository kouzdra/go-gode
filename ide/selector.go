package ide

import "log"
import "strings"
import "github.com/kouzdra/go-gode/icons"
import "github.com/mattn/go-gtk/gtk"
import "github.com/mattn/go-gtk/gdk"
import "github.com/mattn/go-gtk/gdkpixbuf"
import "github.com/mattn/go-gtk/glib"

var _ = log.Printf

type Selector struct {
	ide       *IDE
	Dialog    *gtk.Dialog
	Entry     *gtk.Entry
	View      *gtk.TreeView
	Store     *gtk.TreeStore
	Accel     *gtk.AccelGroup
	Elems     []SelElem
}

type SelElem struct {
	Icon icons.Icon
	Name string
	Loc  Loc
}

func (ide *IDE) NewSelector () *Selector {
	selector := &Selector{}

	selector.ide = ide
	selector.Elems = nil
	selector.Dialog = gtk.NewDialog()
	//selector.Dialog.Connect("close", func () { } )
	selector.Dialog.Response(func() {
		log.Println("Selector closed")
		selector.Dialog.Destroy()
	})

	vbox := selector.Dialog.GetContentArea()
	selector.Entry = gtk.NewEntry()
	vbox.PackStart(selector.Entry, false, false, 0)
	selector.Entry.Connect("changed", func () {
		log.Printf("##CHAGNED: %s", selector.Entry.GetText())
		selector.Reset()
	})
	//selector.Entry.Connect("key-pressed-event", selector.keyPressed)


	selector.View, selector.Store = createList ()
	selector.View.Connect("row_activated", func() {
		log.Printf("## Row activated")
	})
	vbox.PackStart(selector.View, true, true, 0)
	//selector.View.Connect("key-pressed-event", selector.keyPressed)

	swinT := gtk.NewScrolledWindow(nil, nil)
	swinT.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swinT.SetShadowType(gtk.SHADOW_NONE)
	swinT.AddWithViewPort (selector.View)

	selector.Dialog.SetDecorated(false)
	selector.Dialog.SetSizeRequest(400, 300)

	return selector
}

func (selector *Selector) Run () {
	selector.Dialog.ShowAll()
	selector.Dialog.Run()
}

func (selector *Selector) Set (Elems []SelElem) {
	selector.Elems = Elems
}

func (selector *Selector) Reset () {
	prefix := selector.Entry.GetText()
	selector.Store.Clear ()
	if len (prefix) != 0 {
		for _, elem := range selector.Elems {
			if (strings.HasPrefix (elem.Name, prefix)) {
				var iter gtk.TreeIter
				selector.Store.Append(&iter, nil)
				selector.Store.Set(&iter, elem.Icon.GPixbuf, elem.Name, elem.Loc.FName)
			}
		}
	}
}

func createList () (*gtk.TreeView, *gtk.TreeStore) {
	store := gtk.NewTreeStore(gdkpixbuf.GetType(), glib.G_TYPE_STRING, glib.G_TYPE_STRING)
	view  := gtk.NewTreeView()
	model := store.ToTreeModel()
	view.SetModel(model)
	view.AppendColumn(gtk.NewTreeViewColumnWithAttributes("pixbuf", gtk.NewCellRendererPixbuf(), "pixbuf", COL_ICON))
	view.AppendColumn(gtk.NewTreeViewColumnWithAttributes("text"  , gtk.NewCellRendererText  (), "text"  , COL_FNAME))
	pathCol := gtk.NewTreeViewColumnWithAttributes("text"  , gtk.NewCellRendererText  (), "text"  , COL_FPATH)
	pathCol.SetVisible (false)
	view.AppendColumn(pathCol)
	view.SetHeadersVisible (false)
	return view, store
}

func (selector *Selector) keyPressed (key *gdk.EventKey) {
	log.Printf ("## KEY PRESSED: %s", key.String)
}

func (ide *IDE) Select () {
	selector := ide.NewSelector()
	mk := func (n string) SelElem { return SelElem{Icon:ide.Icons.Dir, Name:n, Loc:Loc{n, 1, 1} } }
	elems := []SelElem{mk("AAAAA"), mk("AABB") }
	selector.Set(elems)
	selector.Run ()
}

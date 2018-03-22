package ide

import "log"
import "sort"
import "strings"
import "strconv"
import "github.com/kouzdra/go-gode/icons"
import "github.com/mattn/go-gtk/gtk"
//import "github.com/mattn/go-gtk/gdk"
import "github.com/mattn/go-gtk/gdkpixbuf"
import "github.com/mattn/go-gtk/glib"

var _ = log.Printf

type Selector struct {
	ide       *IDE
	Dialog    *gtk.Dialog
	Entry     *gtk.Entry
	System    *gtk.CheckButton
	View      *gtk.TreeView
	Store     *gtk.TreeStore
	Accel     *gtk.AccelGroup
	Elems     []SelElem
	Result    *SelElem
}

type SelElem struct {
	Icon icons.Icon
	Name string
	Loc  Loc
}

func (ide *IDE) NewSelector () *Selector {
	selector := &Selector{}

	selector.ide = ide
	selector.Elems  = nil
	selector.Result = nil
	selector.Dialog = gtk.NewDialog()
	//selector.Dialog.Connect("close", func () { } )
	selector.Dialog.Response(func() {
		selector.Dialog.Destroy()
	})

	vbox := selector.Dialog.GetContentArea()

	hbox := gtk.NewHBox (false, 0)
	selector.Entry = gtk.NewEntry()
	hbox.PackStart(selector.Entry, true, true, 0)
	selector.Entry.Connect("changed", func () {
		selector.Reset()
	})

	selector.System = gtk.NewCheckButtonWithLabel("System")
	hbox.PackStart(selector.System, false, false, 0)
	selector.System.Connect("toggled", func () {
		selector.Reset()
	})
	
	vbox.PackStart(hbox, false, false, 0)
	
	selector.View, selector.Store = createList ()
	selector.View.Connect("row_activated", func() {
		var path *gtk.TreePath
		var column *gtk.TreeViewColumn
		var iter gtk.TreeIter
		selector.View.GetCursor(&path, &column)
		log.Printf (">> TreePath is: %s", path.String())
		model := selector.View.GetModel()
		if model.GetIter (&iter, path) {
			var val glib.GValue
			model.GetValue (&iter, COL_NO, &val)
			if no, err := strconv.Atoi (val.GetString ()); err == nil {
				log.Printf (">> ROW NO: %d, loc=%s", no, selector.Elems[no].Loc)
				selector.Result = &selector.Elems[no]
				selector.Dialog.Response(gtk.RESPONSE_ACCEPT)
			} else {
				log.Printf ("Invalid NO")
			}
		} else {
			log.Printf ("Invalid path")
		}
	})
	vbox.PackStart(selector.View, true, true, 0)

	swinT := gtk.NewScrolledWindow(nil, nil)
	swinT.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swinT.SetShadowType(gtk.SHADOW_NONE)
	swinT.AddWithViewPort (selector.View)

	selector.Dialog.SetDecorated(false)
	selector.Dialog.SetSizeRequest(700, 350)

	return selector
}

func (selector *Selector) Run () gtk.ResponseType {
	selector.Dialog.ShowAll()
	return selector.Dialog.Run()
}

func (selector *Selector) Set (Elems []SelElem) {
	selector.Elems = Elems
	selector.Reset()
}

func (selector *Selector) Reset () {
	prefix := selector.Entry.GetText()
	system := selector.System.GetActive()
	selector.Store.Clear ()
	counter := 0
	for i, elem := range selector.Elems {
		if strings.HasPrefix (elem.Name, prefix) {
			readable, sys := selector.ide.ReadablePath (elem.Loc.FName)
			if system || !sys {
				var iter gtk.TreeIter
				selector.Store.Append(&iter, nil)
				selector.Store.Set(&iter, elem.Icon.GPixbuf, elem.Name, readable, strconv.Itoa(i))
				counter ++
				if counter >= 1000 {
					return;
				}
			}
		}
	}
}

func createList () (*gtk.TreeView, *gtk.TreeStore) {
	store := gtk.NewTreeStore(gdkpixbuf.GetType(), glib.G_TYPE_STRING, glib.G_TYPE_STRING, glib.G_TYPE_STRING)
	view  := gtk.NewTreeView()
	model := store.ToTreeModel()
	view.SetModel(model)
	view.AppendColumn(gtk.NewTreeViewColumnWithAttributes("pixbuf", gtk.NewCellRendererPixbuf(), "pixbuf", COL_ICON))
	view.AppendColumn(gtk.NewTreeViewColumnWithAttributes("text"  , gtk.NewCellRendererText  (), "text"  , COL_FNAME))
	pathCol := gtk.NewTreeViewColumnWithAttributes("text"  , gtk.NewCellRendererText  (), "text"  , COL_FPATH)
	view.AppendColumn(pathCol)
	numCol := gtk.NewTreeViewColumnWithAttributes("text"   , gtk.NewCellRendererText  (), "text"  , COL_NO)
	view.AppendColumn(numCol)
	numCol.SetVisible(false)
	view.SetHeadersVisible(false)
	return view, store
}

type SortSelElems []SelElem
func (elems SortSelElems) Len () int { return len (elems) }
func (elems SortSelElems) Swap (i, j int) { elems[i], elems[j] = elems[j], elems[i] }
type SortElems struct { SortSelElems }
func (e SortElems) Less (i, j int) bool { return e.SortSelElems[i].Name < e.SortSelElems[j].Name }

func (ide *IDE) Select (elems []SelElem) *SelElem {
	selector := ide.NewSelector()
	sort.Sort(SortElems{SortSelElems(elems)})
	selector.Set(elems)
	selector.Run()
	return selector.Result
}

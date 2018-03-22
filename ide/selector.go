package ide

import "log"
import "sort"
import "strings"
import "strconv"
import "path/filepath"
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
		//log.Println("Selector closed")
		selector.Dialog.Destroy()
	})

	vbox := selector.Dialog.GetContentArea()

	hbox := gtk.NewHBox (false, 0)
	selector.Entry = gtk.NewEntry()
	hbox.PackStart(selector.Entry, true, true, 0)
	selector.Entry.Connect("changed", func () {
		//log.Printf("##CHAGNED: %s", selector.Entry.GetText())
		selector.Reset()
	})

	selector.System = gtk.NewCheckButtonWithLabel("system:")
	hbox.PackStart(selector.System, false, false, 0)
	selector.System.Connect("toggled", func () {
		selector.Reset()
	})
	
	vbox.PackStart(hbox, false, false, 0)
	
	selector.View, selector.Store = createList ()
	selector.View.Connect("row_activated", func() {
		log.Printf("## Row activated")
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

func (selector *Selector) Run () {
	selector.Dialog.ShowAll()
	selector.Dialog.Run()
}

func (selector *Selector) Set (Elems []SelElem) {
	selector.Elems = Elems
	selector.Reset()
}

func (ide *IDE) ReadablePath (fname string) (string, bool) {
	path := filepath.Clean(filepath.Dir(fname))
	goroot := filepath.Clean(ide.Prj.Context.GOROOT)
	gopath := filepath.Clean(ide.Prj.Context.GOPATH)
	//log.Printf("1: path=%s, ROOT=%s PATH=%s", path, goroot, gopath)
	if strings.HasPrefix (path, goroot) {
		return filepath.Join ("$GOROOT", strings.TrimPrefix (path, goroot)), true
	}
	if strings.HasPrefix (path, gopath) {
		return filepath.Join ("$GOPATH", strings.TrimPrefix (path, gopath)), false
	}
	return filepath.Clean (path), false
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

func (ide *IDE) Select (elems []SelElem) {
	selector := ide.NewSelector()
	sort.Sort(SortElems{SortSelElems(elems)})
	selector.Set(elems)
	selector.Run ()
}

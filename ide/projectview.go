package ide

import "log"
import "github.com/mattn/go-gtk/gdkpixbuf"
import 	"github.com/mattn/go-gtk/glib"
import 	"github.com/mattn/go-gtk/gtk"
import 	"github.com/kouzdra/go-analyzer/project"
import 	"os"
import 	"path"

var _ = log.Printf

func  (ide *IDE) fillTree2 (dirs [] project.Dir, iter *gtk.TreeIter) {
	subs := make ([]func (), 0, len(dirs))
	for _, dir := range dirs {
		var subIter gtk.TreeIter
		log.Printf (">> add dir: [%s]\n", dir.Path)
		ide.Store.Append(&subIter, iter)
		name := dir.Path
		if iter != nil {
			name = path.Base (name)
		}
		fPath := dir.Path
		ide.Store.Set(&subIter, ide.Icons.Dir.GPixbuf, name, fPath)
		subs = append (subs, func () { ide.fillTree2 (dir.Sub, &subIter) } )
		if pkg := ide.Prj.Pkgs [dir.Path]; pkg != nil {
			for sPath, src := range pkg.Srcs {
				log.Printf ("   >> add src: [%s | %s]\n", sPath, src.Dir)
				var srcIter gtk.TreeIter
				ide.Store.Append(&srcIter, &subIter)
				ide.Store.Set(&srcIter, ide.Icons.File.GPixbuf, sPath, src.Dir)
			}
		}
	}
	for _, sub := range subs {
		sub ()
	}
}

func  (ide *IDE) fillTree (dirs [] project.Dir) {
	model := ide.View.GetModel ()
	//model.Object.Ref ()
	ide.View.SetModel (nil)
	
	defer func () {
		ide.View.SetModel (model)
		//model.Object.Unref ()
	} ()
	ide.fillTree2 (dirs, nil)
}

func (ide *IDE) LoadView () {
 	log.Printf ("Tree view loading ...")
	ide.fillTree (ide.Prj.Tree)
 	log.Printf ("... OK")
}

const (
	COL_ICON  = 0
	COL_FNAME = 1
	COL_FPATH = 2
)

func (ide *IDE) MakeTree () {
	ide.Store = gtk.NewTreeStore(gdkpixbuf.GetType(), glib.G_TYPE_STRING, glib.G_TYPE_STRING)
	ide.View  = gtk.NewTreeView()
	model := ide.Store.ToTreeModel()
	ide.View.SetModel(model)
	ide.View.AppendColumn(gtk.NewTreeViewColumnWithAttributes("pixbuf", gtk.NewCellRendererPixbuf(), "pixbuf", COL_ICON))
	ide.View.AppendColumn(gtk.NewTreeViewColumnWithAttributes("text"  , gtk.NewCellRendererText  (), "text"  , COL_FNAME))
	pathCol := gtk.NewTreeViewColumnWithAttributes("text"  , gtk.NewCellRendererText  (), "text"  , COL_FPATH)
	pathCol.SetVisible (false)
	ide.View.AppendColumn(pathCol)
	ide.View.SetHeadersVisible (false)

	ide.View.Connect("row_activated", func() {
		var path *gtk.TreePath
		var column *gtk.TreeViewColumn
		var iter gtk.TreeIter
		ide.View.GetCursor(&path, &column)
		mes := "TreePath is: " + path.String()
		model := ide.View.GetModel ()
		if model.GetIter (&iter, path) {
			var val glib.GValue
			model.GetValue (&iter, COL_FNAME, &val)
			mes += ": [" + val.GetString () + "]"
			model.GetValue (&iter, COL_FPATH, &val)
			mes += " (" + val.GetString () + ")"
		} else {
			mes += ": Invalid path"
		}
		dialog := gtk.NewMessageDialog(
			ide.View.GetTopLevelAsWindow(),
			gtk.DIALOG_MODAL,
			gtk.MESSAGE_INFO,
			gtk.BUTTONS_OK,
			mes)
		dialog.SetTitle("TreePath")
		dialog.Response(dialog.Destroy)
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

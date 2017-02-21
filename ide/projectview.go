package ide

import "log"
import "github.com/mattn/go-gtk/gdkpixbuf"
import 	"github.com/mattn/go-gtk/glib"
import 	"github.com/mattn/go-gtk/gtk"
import 	"github.com/kouzdra/go-analyzer/project"
import 	"os"
import 	"path"

var _ = log.Printf

func  (ide *IDE) fillTree (top bool, dirs [] project.Dir, iter *gtk.TreeIter) {
	model := ide.View.GetModel ()
	//model.Object.Ref ()
	ide.View.SetModel (nil)
	
	defer func () {
		ide.View.SetModel (model)
		//model.Object.Unref ()
	} ()
	subs := make ([]func (), 0, len(dirs))
	for _, dir := range dirs {
		var subIter gtk.TreeIter
		ide.Store.Append(&subIter, iter)
		name := dir.Path
		if !top { name = path.Base (name) }
		ide.Store.Set(&subIter, gtk.NewImage().RenderIcon(gtk.STOCK_DIRECTORY,
			gtk.ICON_SIZE_SMALL_TOOLBAR, "").GPixbuf, name)
		subs = append (subs, func () { ide.fillTree (false, dir.Sub, &subIter) } )
		if pkg := ide.Prj.Pkgs [dir.Path]; pkg != nil {
			for sPath, _ := range pkg.Srcs {
				var srcIter gtk.TreeIter
				ide.Store.Append(&srcIter, &subIter)
				ide.Store.Set(&srcIter, gtk.NewImage().RenderIcon(gtk.STOCK_FILE,
					gtk.ICON_SIZE_SMALL_TOOLBAR, "").GPixbuf, sPath)				
			}
		}
	}
	for _, sub := range subs {
		sub ()
	}
}

func (ide *IDE) LoadView () {
 	log.Printf ("Tree view loading ...")
	ide.fillTree (true, ide.Prj.Tree, nil)
 	log.Printf ("... OK")
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
		var iter gtk.TreeIter
		ide.View.GetCursor(&path, &column)
		mes := "TreePath is: " + path.String()
		model := ide.View.GetModel ()
		if model.GetIter (&iter, path) {
			var val glib.GValue
			model.GetValue (&iter, 1, &val)
			mes += ": " + val.GetString ()
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

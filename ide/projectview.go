package ide

import "log"
import "github.com/mattn/go-gtk/gdkpixbuf"
import 	"github.com/mattn/go-gtk/glib"
import 	"github.com/mattn/go-gtk/gtk"
import 	project "github.com/kouzdra/go-analyzer/project/golang"
import 	"os"
import 	"path"

var _ = log.Printf

/*func  (ide *IDE) printTree (dirs [] project.Dir, ind string) {
	for _, dir := range dirs {
		log.Printf ("%s >> dir: [%s]\n", ind, dir.Path)
		//fPath := dir.Path
		ide.printTree (dir.Sub, ind + "  ")
		if pkg := ide.Prj.Pkgs [dir.Path]; pkg != nil {
			for sPath, src := range pkg.Srcs {
				log.Printf ("%s ---: [%s | %s]\n", ind, sPath, src.Dir)
			}
		}
	}
}*/

func  (ide *IDE) fillTree2 (dirs [] project.Dir, iter *gtk.TreeIter, ind string) {
	subs := make ([]func (), 0, len(dirs))
	for _, dir := range dirs {
		var subIter gtk.TreeIter
		//log.Printf ("%s >> add dir: [%s]\n", ind, dir.Path)
		ide.Store.Append(&subIter, iter)
		name := dir.GetPath ().Name
		if iter != nil {
			name = path.Base (name)
		}
		fPath := dir.GetPath ().Name
		readable, _ := ide.ReadablePath (name)
		if iter != nil {
			readable = name
		}
		ide.Store.Set(&subIter, ide.Icons.Dir.GPixbuf, readable, fPath)
		sDirs := dir.GetSub()
		subs = append (subs, func () { ide.fillTree2 (sDirs, &subIter, ind + "  ") } )
		if pkg := ide.Project.GetPackages() [dir.GetPath ()]; pkg != nil {
			for sPath, src := range pkg.GetSrcs () {
				//log.Printf ("%s   >> add src: [%s | %s]\n", ind, sPath.Name, src.Dir)
				var srcIter gtk.TreeIter
				ide.Store.Append(&srcIter, &subIter)
				ide.Store.Set(&srcIter, ide.Icons.File.GPixbuf, sPath.Name, src.GetDir().Name)
			}
		}
	}
	for _, sub := range subs {
		sub ()
	}
}

func  (ide *IDE) fillTree (dirs [] project.Dir) {
	//ide.printTree (dirs, "+++ ")
	model := ide.View.GetModel ()
	//model.Object.Ref ()
	ide.View.SetModel (nil)
	
	defer func () {
		ide.View.SetModel (model)
		//model.Object.Unref ()
	} ()
	ide.fillTree2 (dirs, nil, "## ")
}

func (ide *IDE) LoadView () {
 	log.Printf ("Tree view loading ...")
	ide.fillTree (ide.Project.GetTree())
 	log.Printf ("... OK")
}

const (
	COL_ICON  = 0
	COL_FNAME = 1
	COL_FPATH = 2
	COL_NO    = 3
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
			fName := val.GetString ()
			mes += ": [" + val.GetString () + "]"
			model.GetValue (&iter, COL_FPATH, &val)
			mes += " (" + val.GetString () + ")"
			fPath := val.GetString ()
			ide.Editors.OpenFile (fPath + "/" + fName)
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
	ide.Project = project.NewProject ()
	ide.Project.SetRoot (os.ExpandEnv("$GOROOT"))
	ide.Project.SetPath (os.ExpandEnv("$GOPATH"))
	ide.Project.Load ()
	log.Printf ("Project loaded: #Dirs: %d", len (ide.Project.GetDirs ()))
}

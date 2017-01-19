package ide

import "fmt"
import "log"
import "io/ioutil"
//import "github.com/mattn/go-gtk/gdk"
import "github.com/mattn/go-gtk/gtk"
import gsci "github.com/kouzdra/go-scintilla/gtk"
import consts "github.com/kouzdra/go-scintilla/gtk/consts"
import "github.com/kouzdra/go-analyzer/project"
//import "github.com/kouzdra/go-analyzer/analyzer"
import "github.com/kouzdra/go-gode/faces"

type Editor struct {
	ide *IDE
	Src *project.Src
	Sci *gsci.Scintilla
	FName string
	lockCount int
}

func NewEditor (ide *IDE) *Editor {
	sci := gsci.NewScintilla ()
	faces.Init (sci)
	e := &Editor{ide:ide, Src:nil, Sci:sci, FName:"", lockCount: 0}
	sci.Handlers.OnModify = e.OnModify
	log.Printf ("Editor created\n")
	return e
}

func (e *Editor) DoLock (actions func ()) {
	e.lockCount ++
	defer func () { e.lockCount -- } ()
	actions ()
}

func (e *Editor) OnModify (notificationType uint, pos gsci.Pos, length uint, linesAdded int, text string,
	line uint, foldLevelNow uint, foldLevelPrev uint) {
		if e.lockCount == 0 && e.Src != nil {
			if ((notificationType & consts.SC_MOD_INSERTTEXT) != 0) {
				log.Printf ("SCI INSERT: %x %d #%d (%s) lines=%d\n", notificationType, pos, length, text, linesAdded);
				e.Src.Changed (int (pos), int (pos), text)
				e.Fontify ()
			} else if ((notificationType & consts.SC_MOD_DELETETEXT) != 0) {
				log.Printf ("SCI DELETE: %x %d #%d (%s) lines=%d\n", notificationType, pos, length, text, linesAdded);
				e.Src.Changed (int (pos), int (pos)+int (length), "")
				e.Fontify ()
			}
		}
}

func (e *Editor) LoadFile (fName string) error {
	e.FName = fName
	if src, err := e.ide.Prj.GetSrc (fName); err == nil {
		e.DoLock (func () {
			e.Src = src
			text := src.Text ()
			log.Printf ("Editor load-0 text#=%d\n", len (text))
			log.Printf ("Editor load-1 text#=%d\n", len (e.Src.Text ()))
			e.Sci.SetText (text)
			log.Printf ("Editor load-2 text#=%d\n", len (e.Src.Text ()))
			e.Src.SetText (text) // to block INSERT MESSAGE
			log.Printf ("Editor load-3 text#=%d\n", len (e.Src.Text ()))
		})
		return nil
	} else {
		text, err := ioutil.ReadFile (fName)
		e.Src = nil
		if err == nil {
			e.Sci.SetText (string (text))
		}
		return err
	}
}


func (e *Editor) Fontify () {
	if src, err := e.ide.Prj.GetSrc (e.FName); err == nil {
		es, f := e.ide.Prj.Analyze (src, 0)
		log.Printf ("Fontify  %s", e.FName)
		for _, m := range f.Markers {
			//log.Printf ("  %s at %d:%d\n", m.Color, m.Beg, m.End)
			bg, en := gsci.Pos (m.Beg), gsci.Pos (m.End)
			if f := faces.Faces [m.Color]; f != nil {
				e.Sci.Styling.Range (f.Style, bg, en)
			}
		}
		for _, err := range es.Errors {
			log.Printf ("  %s %s at %d:%d\n", err.Lvl, err.Msg, err.Beg, err.End)
		}
		
	} else {
		log.Printf ("anal.err %s", err)
	}
}

func (e *Editor) LoadFileFromDialog () {
	filechooserdialog := gtk.NewFileChooserDialog(
		"Choose File...",
		e.ide.Window,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		gtk.STOCK_OK,
		gtk.RESPONSE_ACCEPT)
	filter := gtk.NewFileFilter()
	filter.AddPattern("*.go")
	filechooserdialog.AddFilter(filter)
	filechooserdialog.Response(func() {
		fname := filechooserdialog.GetFilename()
		fmt.Println(fname)
		if err := e.LoadFile (fname); err != nil {
			gtk.NewMessageDialog (e.ide.Window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE,
				"error loading file `%s': %s", fname, err)
		} else {
			e.Fontify ()
		}
		filechooserdialog.Destroy()
	})
	filechooserdialog.Run()
}

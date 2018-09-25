package ide

import "fmt"
import "log"
import "io/ioutil"
import "path/filepath"
import "github.com/mattn/go-gtk/gdk"
import "github.com/mattn/go-gtk/gtk"
import gsci "github.com/kouzdra/go-scintilla/gtk"
import consts "github.com/kouzdra/go-scintilla/gtk/consts"
import "github.com/kouzdra/go-analyzer/project/iface"
//import "github.com/kouzdra/go-analyzer/analyzer"
import "github.com/kouzdra/go-gode/faces"

type Indic uint
const INDIC_ERROR Indic = consts.INDIC_CONTAINER

type Editor struct {
	IDE *IDE
	Src iface.Source
	Sci *gsci.Scintilla
	FName string
	LockCount int

	Win    *gtk.ScrolledWindow
	Label  *gtk.Label
	NbIdx   int
}

func (eds *Editors) New (fName string) *Editor {
	sci := gsci.NewScintilla ()
	faces.Init (sci)
	e := &Editor{IDE:eds.IDE, Src:nil, Sci:sci, FName:"", LockCount: 0}
	e.Sci.SetPhasesDraw (consts.SC_PHASES_MULTIPLE)
	e.Sci.AutoCSetDropRestOfWord (true)
	e.Sci.AutoCSetSeparator ('/')
	e.Sci.AutoCSetTypeSeparator (',')

	e.InitIndic ()
	sci.Handlers.OnModify = e.OnModify

	e.Win = gtk.NewScrolledWindow(nil, nil)
	e.Win.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	e.Win.SetShadowType(gtk.SHADOW_IN)
	e.Win.Add(sci)

	e.Label = gtk.NewLabel(filepath.Base (fName))
	e.NbIdx = eds.Notebook.AppendPage(e.Win, e.Label)

	eds.Editors [fName] = e
	
	log.Printf ("Editor created\n")
	e.Win.ShowAll ()

	return e
}

func (e *Editor) Goto (line, col int) {
	pos := e.Sci.FindColumn (line, col)
	e.Sci.GotoPos (pos)
}

func (e *Editor) Close () {
	e.IDE.Editors.Close (e)
	// destory scintilla
}

func (e *Editor) InitIndic () {
	e.Sci.IndicSetStyle (uint (INDIC_ERROR), consts.INDIC_COMPOSITIONTHICK/*SQUIGGLEPIXMAP*/)
	c := gdk.NewColor  ("red")
	mk := func (u uint16) uint32 { return uint32 (u >> 8) }
	e.Sci.IndicSetFg (uint (INDIC_ERROR),
		gsci.Color ((mk (c.Red ()) << 0) | (mk (c.Green ()) << 8) | (mk (c.Blue ()) << 16)))
}

func (e *Editor) GetIdentifier () int {
	return e.Sci.GetIdentifier ()
}

func (e *Editor) DoLock (actions func ()) {
	e.LockCount ++
	defer func () { e.LockCount -- } ()
	actions ()
}

func (e *Editor) OnModify (notificationType uint, pos gsci.Pos, length uint, linesAdded int, text string,
	line uint, foldLevelNow uint, foldLevelPrev uint) {
		if e.LockCount == 0 && e.Src != nil {
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
	if src, err := e.IDE.Project.GetSrc (fName); err == nil {
		e.DoLock (func () {
			e.Src = src
			text := src.GetText ()
			e.Sci.SetText (text)
			e.Src.SetText (text) // to block INSERT MESSAGE
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
	if src, err := e.IDE.Project.GetSrc (e.FName); err == nil {
		es, f := e.IDE.Project.Analyze (src, 0)
		log.Printf ("Fontify  %s", e.FName)
		e.Sci.StyleClear ()
		e.Sci.IndicClear (uint (INDIC_ERROR))
		for _, m := range f.Markers {
			//log.Printf ("  %s at %d:%d\n", m.Color, m.Beg, m.End)
			bg, en := gsci.Pos (m.Beg), gsci.Pos (m.End)
			if f := faces.Faces [m.Color]; f != nil {
				e.Sci.StyleRange (f.Style, bg, en)
			}
		}
		for _, err := range es.Errors {
			log.Printf ("  %s %s at %d:%d\n", err.Lvl, err.Msg, err.Beg, err.End)
			e.Sci.IndicSetRange (uint (INDIC_ERROR), gsci.Pos (err.Beg), gsci.Pos (err.End))
		}
		
	} else {
		log.Printf ("anal.err %s", err)
	}
}

func (e *Editor) LoadFileFromDialog () {
	filechooserdialog := gtk.NewFileChooserDialog(
		"Choose File...",
		e.IDE.Window,
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
			gtk.NewMessageDialog (e.IDE.Window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE,
				"error loading file `%s': %s", fname, err)
		} else {
			e.Fontify ()
		}
		filechooserdialog.Destroy()
	})
	filechooserdialog.Run()
}

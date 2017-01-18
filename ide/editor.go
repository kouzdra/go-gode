package ide

import "fmt"
import "log"
import "io/ioutil"
//import "github.com/mattn/go-gtk/gdk"
import "github.com/mattn/go-gtk/gtk"
import gsci "github.com/kouzdra/go-scintilla/gtk"
import "github.com/kouzdra/go-analyzer/analyzer"
import "github.com/kouzdra/go-gode/faces"

type Editor struct {
	ide *IDE
	Sci *gsci.Scintilla
	FName string
}

func NewEditor (ide *IDE) *Editor {
	sci := gsci.NewScintilla ()
	faces.Init (sci)
	return &Editor{ide, sci, ""}
}

func (e *Editor) LoadFile (fName string) error {
	text, err := ioutil.ReadFile (fName)
	if err == nil {
		e.FName = fName
		e.Sci.SetText (string (text))
	}
	return err
}


func (e *Editor) Fontify () {
	if src, err := e.ide.Prj.GetSrc (e.FName); err == nil {
		_, f := e.ide.Prj.Analyze (src, 0)
		for _, m := range f.Markers {
			log.Printf ("  %s at %d:%d\n", m.Color, m.Beg, m.End)
			bg, en := gsci.Pos (m.Beg), gsci.Pos (m.End)
			switch m.Color {
			case analyzer.Operator : e.Sci.Styling.Range (faces.Operator .Style, bg, en)
			case analyzer.Separator: e.Sci.Styling.Range (faces.Separator.Style, bg, en)
			case analyzer.Keyword  : e.Sci.Styling.Range (faces.Keyword  .Style, bg, en)
			case analyzer.Error    : e.Sci.Styling.Range (faces.Error    .Style, bg, en)
			}
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

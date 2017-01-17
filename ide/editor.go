package ide

import "fmt"
import "log"
import "io/ioutil"
//import "github.com/mattn/go-gtk/gdk"
import "github.com/mattn/go-gtk/gtk"
import gsci "github.com/kouzdra/go-scintilla/gtk"
import "github.com/kouzdra/go-analyzer/analyzer"

type Editor struct {
	ide *IDE
	Sci *gsci.Scintilla
}

func NewEditor (ide *IDE) *Editor {
	sci := gsci.NewScintilla ()
	return &Editor{ide, sci}
}

func (e *Editor) LoadFile (fname string) error {
	text, err := ioutil.ReadFile (fname)
	if err == nil {
		e.Sci.SetText (string (text))
	}
	return err
}

func (e *Editor) LoadFileFromDialog () {
	filechooserdialog := gtk.NewFileChooserDialog(
		"Choose File...",
		e.ide.window,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		gtk.STOCK_OK,
		gtk.RESPONSE_ACCEPT)
	filter := gtk.NewFileFilter()
	filter.AddPattern("*.go")
	filechooserdialog.AddFilter(filter)
	filechooserdialog.Response(func() {
		fname := filechooserdialog.GetFilename()
		fmt.Println(fname)
		if err := e.ide.editor.LoadFile (fname); err != nil {
			gtk.NewMessageDialog (e.ide.window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE,
				"error loading file `%s': %s", fname, err)
		} else {
			if src, err := e.ide.prj.GetSrc (fname); err == nil {
				_, f := e.ide.prj.Analyze (src, 0)
				for _, m := range f.Markers {
					//log.Printf ("  %s at %d:%d\n", m.Color, m.Beg, m.End)
					
					switch m.Color {
					case analyzer.Keyword:
						e.ide.editor.Sci.Styling.Start (gsci.Pos (m.Beg))
						e.ide.editor.Sci.Styling.Set (uint (m.End-m.Beg), e.ide.RED)
					}
				}
				
			} else {
				log.Printf ("anal.err %s", err)
			}
		}
		filechooserdialog.Destroy()
	})
	filechooserdialog.Run()
}

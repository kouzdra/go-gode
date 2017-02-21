package ide

import "log"
import "github.com/mattn/go-gtk/gtk"
//import gsci   "github.com/kouzdra/go-scintilla/gtk"
//import consts "github.com/kouzdra/go-scintilla/gtk/consts"
//import "github.com/kouzdra/go-gode/faces"

type Editors struct {
	IDE *IDE
	Editors map[string] *Editor
	Notebook *gtk.Notebook
}

type Loc struct {
	FName string
	Line, Col int
}

func NewEditors (ide *IDE) *Editors {
	return &Editors{IDE:ide,
		Editors: make (map [string]*Editor),
		Notebook: gtk.NewNotebook()}
}

func (eds *Editors) OpenLoc (loc Loc) *Editor {
	ed := eds.Editors [loc.FName]
	if ed == nil {
		ed = eds.New (loc.FName)
		ed.LoadFile (loc.FName)
		ed.Fontify ()
	}
	eds.SetCurrent (ed)
	ed.Goto (loc.Line, loc.Col)
	return ed
}

func (eds *Editors) OpenFile (fName string) *Editor {
	log.Printf ("-->> File [%s] opened\n", fName)
	return eds.OpenLoc (Loc{fName, 0, 0})
}

func (eds *Editors) GetCurrent () *Editor {
	no := eds.Notebook.GetCurrentPage ()
	for _, ed := range eds.Editors {
		if ed.NbIdx == no {
			return ed
		}
	}
	panic ("no current page in page table")
}

func (eds *Editors) SetCurrent (ed *Editor) {
	eds.Notebook.SetCurrentPage (ed.NbIdx)
}



func (eds *Editors) Close (e *Editor) {
	//e.ide.Editors.Eds [e.GetIdentifier ()] = nil
	// destory scintilla
}


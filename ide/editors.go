package ide

//import "log"
import "github.com/mattn/go-gtk/gtk"
//import gsci   "github.com/kouzdra/go-scintilla/gtk"
//import consts "github.com/kouzdra/go-scintilla/gtk/consts"
//import "github.com/kouzdra/go-gode/faces"

type Editors struct {
	ide *IDE
	Pages map[int]   *Page
	Fn2Id map[string] int
	Notebook *gtk.Notebook
	
}

func NewEditors (ide *IDE) *Editors {
	return &Editors{ide:ide,
		Pages: make (map [int   ]*Page),
		Fn2Id: make (map [string] int ),
		Notebook: gtk.NewNotebook()}
}

type Page struct {
	Win    *gtk.ScrolledWindow
	Label  *gtk.Label
	Editor *Editor
	NbIdx   int
}

func (eds *Editors) GetCurrent () *Page {
	no := eds.Notebook.GetCurrentPage ()
	for _, page := range eds.Pages {
		if page.NbIdx == no {
			return page
		}
	}
	panic ("no current page in page table")
}

func (eds *Editors) Close (e *Editor) {
	//e.ide.Editors.Eds [e.GetIdentifier ()] = nil
	// destory scintilla
}


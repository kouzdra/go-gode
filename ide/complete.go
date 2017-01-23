package ide

import (
	"log"
	//"github.com/mattn/go-gtk/gdkpixbuf"
	//"github.com/mattn/go-gtk/glib"
	//"github.com/mattn/go-gtk/gdk"
	//"github.com/mattn/go-gtk/gtk"
	//gsci "github.com/kouzdra/go-scintilla/gtk"
	//"github.com/kouzdra/go-analyzer/project"
	//"os"
	//"fmt"
	//"path"
)

var _ = log.Printf



//-------------------------------------------------

func (ide *IDE) Complete () {
	if page := ide.Editors.GetCurrent (); page != nil {
		log.Printf ("Complete [%s]\n", page.Editor.FName)
		pos := page.Editor.Sci.GetCurrentPos () + 1
		if src := page.Editor.Src; src == nil {
			log.Printf ("Complete [%s] at %d, no SRC found\n", page.Editor.FName, pos)
		} else {
			log.Printf ("Complete [%s|%s::%s] at %d\n", page.Editor.FName, src.Dir, src.Name, pos)
			if compl := ide.Prj.Complete (src, int (pos)); compl == nil {
				log.Printf ("    -- No completion context found\n")
			} else {
				log.Printf("  [%s/%s] (%d/%d) #%d\n", compl.Pref, compl.Name, compl.Pos, compl.End, len (compl.Choices))
				list := ""
				sep := ""
				for i, c := range compl.Choices {
					log.Printf ("    %d) %s(%s) AKA [%s] pos=%d end=%d\n", i, c.Kind, c.Name, c.Full, c.Pos, c.End)
					list += sep + c.Name
					sep = "/"
				}
				page.Editor.Sci.AutoCShow (uint (int (compl.Pos) - int (pos)), list)
				
			}
		}
	}
}
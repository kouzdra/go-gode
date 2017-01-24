package ide

import "log"

func (ide *IDE) TechShowSelection () {
	log.Printf (">> Selected: [%s]\n", ide.Editors.GetCurrent ().Editor.Sci.GetSelText ())
}

package ide

import "log"
import project "github.com/kouzdra/go-analyzer/gproject"

func (ide *IDE) getFiles (elems []SelElem, dirs []project.Dir) []SelElem {
	for _, dir := range dirs {
		if pkg := ide.Prj.GetPackages() [dir.GetPath()]; pkg != nil {
			for sPath, src := range pkg.GetSrcs () {
				elems = append(elems, SelElem{
					ide.Icons.File, sPath.Name, Loc{src.FName (), 0, 0}})
			}
		}
		elems = ide.getFiles(elems, dir.Sub)
	}
	return elems
}

func (ide *IDE) SelectFiles () {
	elems := ide.getFiles(make ([]SelElem, 0, 100), ide.Prj.Tree)
	res := ide.Select(elems)
	if res != nil {
		log.Printf(">> OPEN LOC")
		ide.Editors.OpenLoc(res.Loc)
	}
}

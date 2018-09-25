package ide

import "log"
import "github.com/kouzdra/go-analyzer/project/iface"

func (ide *IDE) getFiles (elems []SelElem, dirs [] iface.Dir) []SelElem {
	for _, dir := range dirs {
		if pkg := ide.Project.GetPackages() [dir.GetPath()]; pkg != nil {
			for sPath, src := range pkg.GetSrcs () {
				elems = append(elems, SelElem{
					ide.Icons.File, sPath.Name, Loc{iface.FName (src), 0, 0}})
			}
		}
		elems = ide.getFiles(elems, dir.GetSub())
	}
	return elems
}

func (ide *IDE) SelectFiles () {
	elems := ide.getFiles(make ([]SelElem, 0, 100), ide.Project.GetTree())
	res := ide.Select(elems)
	if res != nil {
		log.Printf(">> OPEN LOC")
		ide.Editors.OpenLoc(res.Loc)
	}
}

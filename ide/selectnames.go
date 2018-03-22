package ide

func (ide *IDE) SelectNames () {
	mk := func (n string) SelElem { return SelElem{Icon:ide.Icons.Dir, Name:n, Loc:Loc{n, 1, 1} } }
	elems := []SelElem{ mk("AAAAA"), mk("AABB") }
	ide.Select(elems)
}

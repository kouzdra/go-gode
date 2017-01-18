package faces

import "github.com/mattn/go-gtk/gdk"
//import "github.com/mattn/go-gtk/gtk"
import gsci "github.com/kouzdra/go-scintilla/gtk"

type Face struct {
	Id string
	Style gsci.Style
	Nm string
	Flags uint
}

var Faces = make (map [string] *Face)

const (
	Bold      = 1 
	Italic    = 2 
	Underline = 4
)

func def (id string, nm string, flags uint) *Face {
	f := &Face{Id: id, Style: gsci.Style (0), Nm: nm, Flags: flags}
	Faces [id] = f
	return f
}

const DefaultFont = "Serif"

var (
	Operator   = def ("Operator" , "green"  , 0)
	Separator  = def ("Separator", "royal blue", 0) 
	Keyword    = def ("Keyword"  , "royal blue", Underline|Bold)
	VarRef     = def ("Var"      , "magenta", Bold)
	VarDef     = def ("VarDef"   , "magenta", Bold|Italic)
	TypRef     = def ("Type"     , "light sky blue", Bold)
	TypDef     = def ("TypeDef"  , "light sky blue", Bold|Italic)
	FunRef     = def ("Meth"     , "olive drab", Bold)
	FunDef     = def ("MethDef"  , "olive drab", Bold|Italic)
	ConRef     = def ("Meth"     , "lawn green", Bold)
	ConDef     = def ("MethDef"  , "lawn green", Bold|Italic)
	Comment    = def ("Comment"  , "gray", 0)
	Token      = def ("Token"    , "dark goldenrod", 0)
	Error      = def ("Error"  , "red", Underline)
	String     = def ("String" , "medium blue", 0)
	Char       = def ("Char"   , "medium blue", 0)
	Number     = def ("Number" , "blue", 0)
)

func Init (sci * gsci.Scintilla) {
	style := gsci.Style (64)
	for _, f := range Faces {
		f.Style = style
		style = style + 1
		s := sci.Styling
		c := gdk.NewColor (f.Nm)
		mk := func (u uint16) uint32 { return uint32 (u >> 8) }
		s.SetFont (f.Style, DefaultFont)
		s.SetFg (f.Style, gsci.Color ((mk (c.Red ()) << 0) | (mk (c.Green ()) << 8) | (mk (c.Blue ()) << 16)));
		s.SetUnderline (f.Style, (f.Flags & Underline) != 0);
		s.SetItalic    (f.Style, (f.Flags & Italic   ) != 0);
		s.SetBold      (f.Style, (f.Flags & Bold     ) != 0);
	}
}

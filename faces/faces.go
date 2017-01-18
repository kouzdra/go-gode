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

var (
	Operator   = def ("Operator" , "green"  , Underline)
	Separator  = def ("Separator", "magenta", Underline) 
	Keyword    = def ("Keyword"  , "DarkSlateGray", Underline|Bold)
	VarRef     = def ("Var"      , "cyan", Underline|Bold)
	VarDef     = def ("VarDef"   , "cyan", Underline|Bold|Italic)
	//Comment = "Comment"
	//Token   = "Token"
	//Error  = &Face{Style: gsci.Style (67), Nm: "red", Bd: true, Ul: true}
	/*String  = "String"
	Char    = "Char"
	Number  = "Number"

	VarRef  = "Var"
	VarDef  = "VarDef"
	ConRef  = "Con"
	ConDef  = "ConDef"
	TypRef  = "Type"
	TypDef  = "TypeDef"
	FunRef  = "Meth"
	FunDef  = "MethDef"*/
)

func Init (sci * gsci.Scintilla) {
	style := gsci.Style (64)
	for _, f := range Faces {
		f.Style = style
		style = style + 1
		s := sci.Styling
		c := gdk.NewColor (f.Nm)
		s.SetFg (f.Style, gsci.Color ((byte (c.Red ()) << 0) | (byte (c.Green ()) << 8) | (byte (c.Blue ()) << 16)));
		s.SetUnderline (f.Style, (f.Flags & Underline) != 0);
		s.SetItalic    (f.Style, (f.Flags & Italic   ) != 0);
		s.SetBold      (f.Style, (f.Flags & Bold     ) != 0);
	}
}

package faces

import "github.com/mattn/go-gtk/gdk"
//import "github.com/mattn/go-gtk/gtk"
import gsci "github.com/kouzdra/go-scintilla/gtk"

type Face struct {
	Style gsci.Style
	Nm string
	It bool
	Bd bool
	Ul bool
}

var (
	Operator   = &Face{Style: gsci.Style (64), Nm: "green", Bd: false, Ul: true}
	Separator  = &Face{Style: gsci.Style (65), Nm: "blue" , Bd: false, Ul: true}
	Keyword    = &Face{Style: gsci.Style (66), Nm: "DarkSlateGray", Bd: true, Ul: true}
	//Comment = "Comment"
	//Token   = "Token"
	Error  = &Face{Style: gsci.Style (67), Nm: "red", Bd: true, Ul: true}
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

func (f *Face) Init (sci * gsci.Scintilla) {
	s := sci.Styling
	c := gdk.NewColor (f.Nm)
	s.SetFg (f.Style, gsci.Color (c.Red () | (c.Green () << 8) | (c.Blue () << 16)));
	s.SetUnderline (f.Style, f.Ul);
}

func Init (sci * gsci.Scintilla) {
	Operator .Init (sci)
	Separator.Init (sci)
	Keyword  .Init (sci)
	Error    .Init (sci)
}

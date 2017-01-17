package main

import (
	"log"
	_ "github.com/mattn/go-gtk/gdkpixbuf"
	_ "github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	gsci "github.com/kouzdra/go-scintilla/gtk"
	"os"
	"io/ioutil"
	"fmt"
)

var _ = log.Printf

type Editor struct {
	Sci *gsci.Scintilla
}

func NewEditor () *Editor {
	sci := gsci.NewScintilla ()
	return &Editor{sci}
}

func (e *Editor) LoadFile (fname string) error {
	text, err := ioutil.ReadFile (fname)
	if err == nil {
		e.Sci.SetText (string (text))
	}
	return err
}


type IDE struct {
	window  *gtk .Window
	editor  *Editor
	menubar *gtk .MenuBar
	RED      gsci.Style
}

func NewIDE () *IDE {
	ide := &IDE{}
	
	ide.window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	ide.window.SetTitle("Editor")
	ide.window.Connect("destroy", gtk.MainQuit)

	vbox := gtk.NewVBox(false, 1)
	ide.editor = NewEditor()
	ide.MakeMenu()
	vbox.PackStart(ide.menubar, false, false, 0)

	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetShadowType(gtk.SHADOW_IN)
	vbox.Add(swin)

	swin.Add(ide.editor.Sci)

	ide.window.Add(vbox)
	ide.window.SetSizeRequest(400, 300)
	ide.window.ShowAll()

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Editor")
	window.Connect("destroy", gtk.MainQuit)

	return ide
}


func main() {
	gtk.Init(&os.Args)
	ide := NewIDE ()
	ide.PreloadTest ()
	gtk.Main()
}

func (ide *IDE) PreloadTest () {
	ide.RED = gsci.Style (1)

	s := ide.editor.Sci.Styling
	s.ResetDefault()
	s.SetFg (ide.RED, gsci.Color (0x0000FF));
	//s.SetBg (ide.RED, gsci.``Color (0x808080));
	s.SetUnderline (ide.RED, false);
	//s.SetFont (ide.RED, "Sans Bold Italic 10")

	ide.editor.Sci.SetText(`#include <iostream>
template<class T>
struct foo_base {
  T operator+(T const &rhs) const {
    T tmp(static_cast<T const &>(*this));
    tmp += rhs;
    return tmp;
  }
};
`)
	s.GetAt (5)
	s.Start (3)
	s.Set (10, ide.RED)
	log.Printf ("AT=%d\n", s.GetAt (5))
	log.Printf ("AT=%d\n", ide.editor.Sci.GetCharAt (5))
	s.GetEnd ()

}

func (ide *IDE) MakeMenu () {

	ide.menubar = gtk.NewMenuBar ()
	addCascade := func (label string, fill func (*gtk.Menu)) {
		cascademenu := gtk.NewMenuItemWithMnemonic(label)
		ide.menubar.Append(cascademenu)
		submenu := gtk.NewMenu()
		cascademenu.SetSubmenu(submenu)
		fill (submenu)
	}
	makeItem := func (label string, action func ()) *gtk.MenuItem {
		menuitem := gtk.NewMenuItemWithMnemonic(label)
		menuitem.Connect("activate", action)
		return menuitem
	}
	
	addCascade ("_File", func (submenu *gtk.Menu) {
		submenu.Append(makeItem ("E_xit", gtk.MainQuit))
		submenu.Append(makeItem ("_Open", func() {
			filechooserdialog := gtk.NewFileChooserDialog(
				"Choose File...",
				ide.window,
				gtk.FILE_CHOOSER_ACTION_OPEN,
				gtk.STOCK_OK,
				gtk.RESPONSE_ACCEPT)
			filter := gtk.NewFileFilter()
			filter.AddPattern("*.go")
			filechooserdialog.AddFilter(filter)
			filechooserdialog.Response(func() {
				fname := filechooserdialog.GetFilename()
				fmt.Println(fname)
				if err := ide.editor.LoadFile (fname); err != nil {
					gtk.NewMessageDialog (ide.window, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE,
						"error loading file `%s': %s", fname, err)
				}
				filechooserdialog.Destroy()
			})
			filechooserdialog.Run()
		}))
	})

	addCascade ("_View", func (submenu *gtk.Menu) {
		submenu.Append(makeItem ("_Font", func () {
			fsd := gtk.NewFontSelectionDialog("Font")
			fsd.Response(func() {
				fmt.Println(fsd.GetFontName())
				ide.editor.Sci.Styling.SetFont(ide.RED, fsd.GetFontName())
				ide.editor.Sci.Styling.SetUnderline (ide.RED, true);
				fsd.Destroy()
			})
			fsd.SetTransientFor(ide.window)
			fsd.Run()
		}))
	})

	addCascade ("_Help", func (submenu *gtk.Menu) {
		submenu.Append(makeItem ("_About", func () {
			dialog := gtk.NewAboutDialog()
			dialog.SetName("Go-Gtk Demo!")
			dialog.SetProgramName("demo")
			dialog.SetAuthors([]string{"Kouzdra"})
			//dir, _ := filepath.Split(os.Args[0])
			//imagefile := filepath.Join(dir, "../../data/mattn-logo.png")
			//pixbuf, _ := gdkpixbuf.NewPixbufFromFile(imagefile)
			//dialog.SetLogo(pixbuf)
			dialog.SetLicense("The library is available under the same terms and conditions as the Go, the BSD style license, and the LGPL (Lesser GNU Public License). The idea is that if you can use Go (and Gtk) in a project, you should also be able to use go-gtk.")
			dialog.SetWrapLicense(true)
			dialog.Run()
			dialog.Destroy()
		}))
	})
}

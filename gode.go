package main

import (
	"log"
	"github.com/mattn/go-gtk/gtk"
	gsci "github.com/kouzdra/go-scintilla/gtkscintilla"
	"os"
)

var _ = log.Printf

func main() {
	gtk.Init(&os.Args)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Editor")
	window.Connect("destroy", gtk.MainQuit)

	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetShadowType(gtk.SHADOW_IN)

	sci := gsci.NewScintilla()
	swin.Add(sci)

	window.Add(swin)
	window.SetSizeRequest(400, 300)
	window.ShowAll()


	RED := uint (1)
	sci.StyleResetDefault()
	sci.StyleSetFg (RED, 0x0000FF);
	//sci.StyleSetBg (RED, 0x808080);
	sci.StyleSetUnderline (RED, false);

	sci.SetText(`#include <iostream>
template<class T>
struct foo_base {
  T operator+(T const &rhs) const {
    T tmp(static_cast<T const &>(*this));
    tmp += rhs;
    return tmp;
  }
};
`)
	sci.GetStyleAt (5)
	sci.StartStyling (3)
	sci.SetStyling (10, RED)
	log.Printf ("AT=%d\n", sci.GetStyleAt (5))
	log.Printf ("AT=%d\n", sci.GetCharAt (5))
	sci.GetEndStyled ()


	gtk.Main()
}

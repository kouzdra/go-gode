package main

//# cgo LDFLAGS='-lgmodule-2.0 --static'

import (
	"github.com/mattn/go-gtk/gtk"
	gsci "github.com/kouzdra/go-scintilla/gtkscintilla"
	"os"
)

func main() {
	gtk.Init(&os.Args)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Editor")
	window.Connect("destroy", gtk.MainQuit)

	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetShadowType(gtk.SHADOW_IN)
	sourceview := gsci.NewScintilla()
	sourceview.InsertText(0, `#include <iostream>
template<class T>
struct foo_base {
  T operator+(T const &rhs) const {
    T tmp(static_cast<T const &>(*this));
    tmp += rhs;
    return tmp;
  }
};
`)

	swin.Add(sourceview)

	window.Add(swin)
	window.SetSizeRequest(400, 300)
	window.ShowAll()

	gtk.Main()
}

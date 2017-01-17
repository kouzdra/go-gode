package main

import "os"
import "github.com/mattn/go-gtk/gtk"
import "github.com/kouzdra/go-gode/ide"

func main() {
	gtk.Init(&os.Args)
	ide := ide.NewIDE ()
	ide.PreloadTest ()
	ide.LoadProject ()
	ide.Editor.LoadFile (os.ExpandEnv ("$GOPATH/src/github.com/kouzdra/go-gode/gode.go"))
	ide.Editor.Fontify ()
	gtk.Main()
}

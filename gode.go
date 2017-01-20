package main

import "os"
import "github.com/mattn/go-gtk/gtk"
import "github.com/kouzdra/go-gode/ide"

func main() {
	gtk.Init(&os.Args)
	ide := ide.NewIDE ()
	ide.LoadProject ()
	ide.LoadView ()
	for _, ed := range ide.Editors.Eds {
		ed.LoadFile (os.ExpandEnv ("$GOPATH/src/github.com/kouzdra/go-gode/gode.go"))
		ed.Fontify ()
	}
	gtk.Main()
}

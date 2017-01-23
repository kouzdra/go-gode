package main

import "os"
import "github.com/mattn/go-gtk/gtk"
import "github.com/kouzdra/go-gode/ide"

func main() {
	gtk.Init(&os.Args)
	ide := ide.NewIDE ()
	ide.LoadProject ()
	{
		ed := ide.Editors.New ("File 1")
		ed.LoadFile (os.ExpandEnv ("$GOPATH/src/github.com/kouzdra/go-gode/gode.go"))
		ed.Fontify ()
	}
	{
		ed := ide.Editors.New ("File 2")
		ed.LoadFile (os.ExpandEnv ("$GOPATH/src/github.com/kouzdra/go-gode/ide/ide.go"))
		ed.Fontify ()
	}
	go ide.LoadView ()

	gtk.Main ()
}

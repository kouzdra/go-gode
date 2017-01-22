package main

import "os"
import "github.com/mattn/go-gtk/gtk"
import "github.com/kouzdra/go-gode/ide"

func main() {
	gtk.Init(&os.Args)
	ide := ide.NewIDE ()
	ide.LoadProject ()
	ide.LoadView ()

	//ide.Editors.New ("File 1")
	//ide.Editors.New ("File 2")
	
	ed1 := ide.Editors.New ("File 1")
	ed2 := ide.Editors.New ("File 2")
	
	ed1.LoadFile (os.ExpandEnv ("$GOPATH/src/github.com/kouzdra/go-gode/gode.go"))
	ed1.Fontify ()
	//ed2.LoadFile (os.ExpandEnv ("$GOPATH/src/github.com/kouzdra/go-gode/ide/editors.go"))
	ed2.LoadFile (os.ExpandEnv ("$GOPATH/src/github.com/kouzdra/go-gode/ide/ide.go"))
	ed2.Fontify ()

	gtk.Main()
}

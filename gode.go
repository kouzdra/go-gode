package main

import "os"
import "github.com/mattn/go-gtk/gtk"
import "github.com/kouzdra/go-gode/ide"

func main() {
	gtk.Init(&os.Args)
	de := ide.NewIDE ()
	de.LoadProject ()
	de.Editors.OpenLoc (ide.Loc {os.ExpandEnv ("$GOPATH/src/github.com/kouzdra/go-gode/gode.go"), 2, 4})
	de.Editors.OpenLoc (ide.Loc {os.ExpandEnv ("$GOPATH/src/github.com/kouzdra/go-gode/ide/ide.go"), 4, 1})
	go de.LoadView ()

	gtk.Main ()
}

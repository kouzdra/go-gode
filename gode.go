package main

import "os"
import "github.com/mattn/go-gtk/gtk"
import "github.com/kouzdra/go-gode/ide"

func main() {
	gtk.Init(&os.Args)
	ide := ide.NewIDE ()
	ide.PreloadTest ()
	ide.LoadProject ()
	gtk.Main()
}

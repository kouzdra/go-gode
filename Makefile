all: gode
	./$<

gode: gode.go ide/ide.go ide/editor.go ide/editors.go ide/menu.go ide/projectview.go ide/build.go faces/faces.go options/options.go
	go build $<


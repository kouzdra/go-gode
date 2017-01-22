all: gode
	./$<

gode: gode.go ide/ide.go ide/editor.go ide/editors.go faces/faces.go
	go build $<


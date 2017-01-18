all: gode
	./$<

gode: gode.go ide/ide.go ide/editor.go faces/faces.go
	go build $<


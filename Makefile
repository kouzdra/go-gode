all: gode
	./$<

gode: gode.go ide/ide.go ide/editor.go
	go build $<


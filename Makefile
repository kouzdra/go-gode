all: gode
	./$<

clean:
	go clean gode.go

gode: gode.go ide/ide.go ide/editor.go ide/complete.go ide/editors.go ide/menu.go ide/projectview.go ide/build.go faces/faces.go options/options.go
	go build $<


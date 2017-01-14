all: gode
	./$<

gode: gode.go
	go build $^


templ:
	templ generate

build:
	go build -o ./tmp/main .

serve:
	./tmp/main

dev: templ build
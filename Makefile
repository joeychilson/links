templ:
	templ generate

css:
	tailwindcss -c ./assets/tailwind.config.js -i ./assets/css/app.css -o ./static/css/app.css

js:
	bun build ./assets/js/app.js --outdir ./static/js

build:
	go build -o ./tmp/main .

serve:
	./tmp/main

dev: templ css js build
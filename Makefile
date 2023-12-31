templ:
	templ generate

css:
	tailwindcss -c ./assets/tailwind.config.js -i ./assets/css/app.css -o ./static/css/app.css --minify

js:
	bun build ./assets/js/app.js --outdir ./static/js --minify

build:
	go build -o ./tmp/main .

serve:
	./tmp/main

dev: templ css js build
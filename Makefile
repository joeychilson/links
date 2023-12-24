dev:
	templ generate && tailwindcss -c ./assets/tailwind.config.js -i ./assets/css/app.css -o ./static/css/app.css && go build -o ./tmp/main .
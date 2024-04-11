build-run: build
	./bin/server

run: 
	go run .

build: tidy
	go build -o bin/server . && cp -r static/ bin/ 

tidy: 
	go mod tidy

tailwind:
	npx tailwindcss -i ./static/css/styles.css -o ./static/css/output.css --watch
build-run: build
	./bin/server

run: 
	go run .

build: tidy
	go build -o bin/server . && cp -r public/ bin/ 

tidy: 
	go mod tidy && go mod download

check-npm:
	which npm || echo -e "\n>>> \u001b[41m NPM is not installed plaese install nodejs and npm \u001b[0m <<<\n"

init-node: check-npm
	npm init

install-tailwind: check-npm
	[ -f tailwind.config.js  ] || npn 

tailwind-watch:
	npx tailwindcss -i ./public/css/base.css -o ./public/css/styles.css --watch

tailwind:
	npx tailwindcss -i ./public/css/base.css -o ./public/css/styles.css

build-docker: which-docker
	docker build .

which-docker:
	which docker

install-air:
	which air || go install github.com/cosmtrek/air@latest

air: install-air
	([ -f .air.toml ] && air -c .air.toml) || air init || air

APP_NAME = htmx
APP_VERSION := $(shell git describe --abbrev=0 || echo v0.1.0)
GIT_HEAD := $(shell git rev-parse --verify HEAD)
BUILD_AT := $(shell date --rfc-3339 'seconds' -u)

build-run: build
	./bin/$(APP_NAME)

run: 
	go run . serve

build: tidy
	go build -ldflags="-w -s -X 'github.com/OZahed/go-htmx/cmd.APP_VERSION=$(APP_VERSION)' -X 'github.com/OZahed/go-htmx/cmd.APP_NAME=$(APP_NAME)' -X 'github.com/OZahed/go-htmx/cmd.GIT_HEAD=$(GIT_HEAD)' -X 'github.com/OZahed/go-htmx/cmd.BUILD_AT=$(BUILD_AT)'" \
		-o bin/$(APP_NAME) . && cp -r public/ bin/ 

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

test: 
	go mod tidy && go test -v --race -p 1 ./internal/... 

APP_NAME = go-htmx
APP_VERSION := $(shell git describe --abbrev=0 || echo v0.1.0)
GIT_HEAD := $(shell git rev-parse --verify HEAD)
BUILD_AT := $(shell date --rfc-3339 'seconds' -u)
GOOS := linux 
GOARCH := amd64

build-run: build
	./bin/$(APP_NAME)

run: build
	./bin/$(APP_NAME) serve

build: tidy
	 go build -ldflags="-w -s -X 'github.com/OZahed/go-htmx/cmd.APP_VERSION=$(APP_VERSION)' -X 'github.com/OZahed/go-htmx/cmd.APP_NAME=$(APP_NAME)' -X 'github.com/OZahed/go-htmx/cmd.GIT_HEAD=$(GIT_HEAD)' -X 'github.com/OZahed/go-htmx/cmd.BUILD_AT=$(BUILD_AT)'" \
		-o bin/$(APP_NAME) . && cp -r public/ bin/ 

build-linux: tidy
	GOGC=off CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-w -s -X 'github.com/OZahed/go-htmx/cmd.APP_VERSION=$(APP_VERSION)' -X 'github.com/OZahed/go-htmx/cmd.APP_NAME=$(APP_NAME)' -X 'github.com/OZahed/go-htmx/cmd.GIT_HEAD=$(GIT_HEAD)' -X 'github.com/OZahed/go-htmx/cmd.BUILD_AT=$(BUILD_AT)'" \
		-o bin/$(APP_NAME) . && cp -r public/ bin/ 

tidy: 
	go mod tidy && go mod download

lint:
	golangci-lint run -c .golangci.yml ./...

tailwind:
	[ ! -f ./tailwindcss ] && echo "\nplease take a look at tailwindcss installation process on officail docs: \n\n\t check standalone CLI process on:\n\t\033[0;33m https://tailwindcss.com/blog/standalone-cli \n\t or put the CLI in ./public directory \n\n\033[0;0m" || echo Tailwind exists

tailwind-watch:
	./tailwindcss -i ./public/css/base.css -o ./public/css/styles.css --watch

tailwind-minify:
	./tailwindcss -i ./public/css/base.css -o ./public/css/styles.css --minify

docker-build: which-docker
	docker build -t go-htmx:latest  .

which-docker:
	which docker

install-air:
	which air || go install github.com/cosmtrek/air@latest

air: install-air
	([ -f .air.toml ] && air -c .air.toml) || air init || air

test: 
	go mod tidy && go test -v --race -p 1 ./internal/... 

ssl-keys: 
	which openssl && openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 7 -out server.crt || echo "\n\t openssl not found\n" 

cleanup:
	rm -rf ./bin
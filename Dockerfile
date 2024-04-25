FROM golang:1.22-bookworm as build


# Add make, ca-certificate, git, tailwindcss
RUN apt-get -y update && apt-get install -y  \ 
  git \
  make \
  ca-certificates \
  openssl && rm -rf /var/cache/apt/archives /var/lib/apt/lists/* 

WORKDIR /app
COPY . . 

RUN go get github.com/a-h/templ
RUN make ssl-keys
RUN make build-linux

FROM alpine:latest AS runtime

# add ca-certificate
USER root
RUN apk --no-cache add ca-certificates \
  && update-ca-certificates

WORKDIR /go-htmx
COPY --from=build /app/bin . 
COPY --from=build /app/templates .

EXPOSE 8080

ENV TEMP_DIR="/go-htmx/templates"
ENV TEMP_ROOT_NAME="Layout"
ENV NAME="go-htmx"

CMD [ "/go-htmx/go-htmx", "serve" ]
FROM golang:latest

LABEL maintainer="Kuramshin Rustam <kuramshin.py@yandex.ru>"

WORKDIR /go/src/go-app

COPY . .

RUN apt-get update

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 8080

CMD [ "go", "run", "http-server.go" ]

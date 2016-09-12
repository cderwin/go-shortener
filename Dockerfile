FROM golang:latest

RUN go get -u github.com/tools/godep
RUN go get -u github.com/golang/lint/golint

WORKDIR /go/src/app
COPY . /go/src/app

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]
EXPOSE 8080

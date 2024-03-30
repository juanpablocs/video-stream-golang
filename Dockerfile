FROM golang:1.22.0

WORKDIR /go/src/app
COPY ./ .

RUN mkdir -p /go/src/app/videos

RUN go get -d -v ./...
RUN go build -o appmain cmd/main.go

CMD ["./appmain"]
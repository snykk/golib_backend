FROM golang:1.19.2-alpine3.16

WORKDIR /myApp
COPY . .
RUN go build -o server ./cmd/api/main.go
RUN cp server /

WORKDIR /
RUN rm -rf ./myApp

CMD ["./server"]
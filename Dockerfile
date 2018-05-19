FROM golang:alpine
RUN apk update && apk add git
RUN go get github.com/ecc1/spi
RUN go get github.com/gorilla/websocket
COPY build/ build/
COPY Server.go Unicorn.go RealUnicorn.go ./
RUN go build -o ./unicornpaint Server.go Unicorn.go RealUnicorn.go
CMD ./unicornpaint

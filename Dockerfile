FROM golang:alpine
RUN apk update && apk add git
RUN go get github.com/ecc1/spi
RUN go get github.com/gorilla/websocket
RUN go get github.com/MFAshby/unicornpaint/unicorn
COPY build/ build/
COPY Server.go ./
RUN go build -o ./unicornpaint Server.go
CMD ./unicornpaint

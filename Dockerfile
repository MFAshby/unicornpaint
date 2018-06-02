FROM golang:alpine
# Add tools for downloading & building
RUN apk add --update git nodejs nodejs-npm

# Copy everything across
COPY . ./

# Build the website
RUN npm install && npm run-script build

# Retrieve server dependencies
RUN go get github.com/ecc1/spi github.com/gorilla/websocket github.com/MFAshby/unicornpaint/unicorn

# Build server
RUN go build -o ./unicornpaint Server2.go

# Get rid of stuff we don't need for runtime
RUN apk del git nodejs nodejs-npm

# Run server!
CMD ./unicornpaint

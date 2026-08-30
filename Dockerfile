FROM golang:alpine AS builder

# We assume only git is needed for all dependencies.
# openssl is already built-in.
RUN apk add -U --no-cache git

RUN adduser -D server
USER server
WORKDIR /home/server

# Cache pulled dependencies if not updated.
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy necessary parts of the Mail-Go source into builder's source
COPY *.go ./
COPY webpanel webpanel
COPY middleware middleware
COPY mii mii

# Build to name "app".
RUN go build -o app .

# Runner
FROM alpine:latest

RUN adduser -D server
USER server
WORKDIR /home/server

# Copy executable
COPY --from=builder /home/server/app .
COPY assets assets
COPY templates templates

CMD ["./app"]
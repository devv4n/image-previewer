FROM golang:1.23-alpine AS builder

RUN apk update && apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY image-previewer .

RUN make build

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/image-previewer .

RUN mkdir -p /app/previews && chmod -R 777 /app/previews

USER nobody
CMD ["./image-previewer"]

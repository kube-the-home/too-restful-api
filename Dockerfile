FROM golang:1.22.6 as builder

ENV CGO_ENABLED=0

COPY . /app

WORKDIR /app

RUN go build -ldflags '-s -w' -trimpath -o 'bin/' ./...

FROM alpine:3.20.2

COPY --from=builder /app/bin/too-restful-api /usr/local/bin/too-restful-api

RUN apk add --no-cache "dumb-init"

ENTRYPOINT ["dumb-init"]
CMD ["too-restful-api"]

FROM golang:1.22.2 as builder

ENV CGO_ENABLED=0

COPY . /app

WORKDIR /app

RUN go build -ldflags '-s -w' -trimpath -o 'bin/' ./...

FROM alpine:3.19.1

COPY --from=builder /app/bin/too-restful-api /usr/local/bin/too-restful-api

RUN apk add --no-cache "dumb-init"

ENTRYPOINT ["dumb-init"]
CMD ["too-restful-api"]

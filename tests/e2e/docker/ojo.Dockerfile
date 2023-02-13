# Builder
FROM golang:1.19-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

# Runner
FROM alpine
RUN apk add bash

COPY --from=builder /app/build/ojod /bin/ojod

EXPOSE 26656 26657 1317 9090
ENTRYPOINT ["ojod"]

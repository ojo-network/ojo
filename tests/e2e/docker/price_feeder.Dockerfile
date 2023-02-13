# Builder
FROM golang:1.19-alpine AS builder

RUN apk add --no-cache \
    ca-certificates \
    build-base \
    linux-headers

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

# Runner
FROM alpine
RUN apk add bash

COPY --from=builder /app/build/price-feeder /bin/price-feeder

EXPOSE 7171
ENTRYPOINT ["price-feeder"]

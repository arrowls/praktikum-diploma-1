FROM golang:1.23-alpine AS builder
WORKDIR /go/src/app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY ../.. .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/app ./cmd/gophermart/main.go

FROM golang:1.23-alpine AS dev
WORKDIR /go/src/app

COPY --from=builder /go/bin/air /go/bin/air
COPY --from=builder /go/src/app /go/src/app

COPY .air.toml ./

CMD ["air", "-c", "air.toml"]
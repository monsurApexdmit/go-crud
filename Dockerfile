FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
# Copy go.sum only if it exists
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o server main.go


FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/server /app/server

EXPOSE 8004

CMD ["./server"]

FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download || true

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/kryvora-node ./cmd/kryvora-node

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /opt/kryvora
COPY --from=builder /app/kryvora-node /usr/local/bin/kryvora-node
COPY config.example.yaml /etc/kryvora/config.yaml

EXPOSE 4177
VOLUME ["/var/lib/kryvora"]

ENTRYPOINT ["/usr/local/bin/kryvora-node"]
CMD ["--config", "/etc/kryvora/config.yaml"]

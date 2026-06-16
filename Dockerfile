FROM golang:1.26.3-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# -trimpath buat ngilangin path di binary
# -ldflags="-s -w" buat ngilangin debug info biar ukuran binary lebih kecil
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o /out/tickitz ./cmd

FROM alpine:3.24.0

WORKDIR /app

RUN apk add --no-cache tzdata \
  && addgroup -S app \
  && adduser -S -G app app \
  && mkdir -p /app/public/img/movies \
  && mkdir -p /app/public/img/profile \
  && chown -R app:app /app

COPY --from=builder /out/tickitz /app/tickitz
COPY --from=builder /src/public /app/public

RUN chown -R app:app /app

USER app

ENTRYPOINT ["/app/tickitz"]

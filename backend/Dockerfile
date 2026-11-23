FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/calculator-api ./cmd/api

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app
COPY --from=builder --chown=app:app /out/calculator-api /usr/local/bin/calculator-api

USER app
EXPOSE 8080

ENTRYPOINT ["calculator-api"]

FROM golang:1.25.1-alpine3.21 as builder

RUN apk --no-cache add gcc shadow musl-dev && apk add tzdata

RUN adduser -D appuser && usermod -u 1000 appuser && groupmod -g 1000 appuser
#USER appuser

WORKDIR /app

COPY environment/dev/.air.toml /tmp/.air.toml
COPY . .

RUN go mod download \
    && go install github.com/go-delve/delve/cmd/dlv@v1.24.0 \
    && go install github.com/cosmtrek/air@v1.40.4

CMD ["air", "-c", "/tmp/.air.toml"]

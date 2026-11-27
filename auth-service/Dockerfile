FROM golang:1.27.4-alpine AS builder

WORKDIR /app

COPY auth-service/go.mod auth-service/go.sum ./auth-service/
COPY contracts/go.mod contracts/go.sum ./contracts/

RUN go work init ./auth-service ./contracts

RUN cd auth-service && go mod download
RUN cd contracts && go mod download


COPY auth-service/internal ./auth-service/internal
COPY auth-service/cmd ./auth-service/cmd
COPY auth-service/jwtauth ./auth-service/jwtauth
COPY contracts/proto ./contracts/proto

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/server ./auth-service/cmd

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app app

COPY --from=builder --chown=app:app /app/server ./server

EXPOSE 8080

USER app
ENTRYPOINT ["./server"]

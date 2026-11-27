FROM golang:1.27.4-alpine AS builder

WORKDIR /app

COPY electrohub/go.mod electrohub/go.sum ./electrohub/
COPY contracts/go.mod contracts/go.sum ./contracts/

RUN go work init ./electrohub ./contracts

RUN cd electrohub && go mod download
RUN cd contracts && go mod download


COPY electrohub/internal ./electrohub/internal
COPY electrohub/cmd ./electrohub/cmd
COPY contracts/proto ./contracts/proto

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/server ./electrohub/cmd

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S -G app app

COPY --from=builder --chown=app:app /app/server ./server

EXPOSE 8080

USER app
ENTRYPOINT ["./server"]

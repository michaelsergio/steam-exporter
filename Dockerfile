FROM golang:alpine as builder

WORKDIR /app

COPY owned_games.go owned_games.go
COPY server.go server.go
COPY go.mod go.mod
COPY go.sum go.sum

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -o steam_exporter .

FROM alpine:3.20

RUN apk update && apk add ca-certificates

COPY --from=builder /app/steam_exporter steam_exporter

ENV STEAM_KEY ""
ENV STEAM_USER_ID ""

ENTRYPOINT ./steam_exporter ${STEAM_USER_ID}

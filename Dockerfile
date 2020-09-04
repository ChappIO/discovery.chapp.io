FROM golang:1.15 as build

WORKDIR /app
COPY cmd cmd
COPY internal internal
COPY go.mod go.mod
COPY go.sum go.sum
RUN CGO_ENABLED=0 go build -o discovery_server ./cmd/discovery_server

FROM scratch

WORKDIR /app
COPY --from=build /app/discovery_server .

ENTRYPOINT ["/app/discovery_server"]

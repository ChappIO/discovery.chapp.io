FROM golang:1.15 as build

WORKDIR /app
COPY main.go .
RUN CGO_ENABLED=0 go build -o discovery_server main.go

FROM scratch

WORKDIR /app
COPY --from=build /app/discovery_server .

ENTRYPOINT ["/app/discovery_server"]

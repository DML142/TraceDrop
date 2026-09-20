FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/tracedrop-api ./cmd/api

FROM alpine:3.20
RUN adduser -D -u 10001 tracedrop
USER tracedrop
COPY --from=build /out/tracedrop-api /tracedrop-api
EXPOSE 8080
ENTRYPOINT ["/tracedrop-api"]

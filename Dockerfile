# Build stage
FROM golang:1.25.1 AS build

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o ngo-site

# Runtime stage
FROM ubuntu:22.04

WORKDIR /app

COPY --from=build /app/ngo-site .
COPY --from=build /app/static ./static
COPY --from=build /app/templates ./templates

EXPOSE 8080

CMD ["./ngo-site"]
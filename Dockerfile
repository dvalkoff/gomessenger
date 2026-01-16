FROM golang:1.25.5 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o messenger-app ./cmd/api/main.go

FROM alpine:latest
COPY --from=build /app/messenger-app /usr/local/bin/messenger-app
EXPOSE 8080
CMD [ "messenger-app" ]

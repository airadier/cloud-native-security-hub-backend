FROM golang:1.13 as builder
WORKDIR /cloud-native-visiblity-hub-backend
COPY go.mod go.sum ./
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s" -o server cmd/server/main.go
RUN strip server

FROM alpine:3.13.9
COPY --from=builder /cloud-native-visiblity-hub-backend/server /bin/server
EXPOSE 8080
CMD ["/bin/server"]

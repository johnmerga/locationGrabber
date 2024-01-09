FROM golang:alpine3.19 as dev
COPY . /go/src/johnmerga/locationGrabber
WORKDIR /go/src/johnmerga/locationGrabber
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/cosmtrek/air@latest
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /server
CMD ["air", "-c", ".air.toml"]
# CMD ["go","run","./main.go"]

## Deploy
FROM alpine:latest as prod
RUN mkdir /data
COPY --from=dev /server ./
COPY ./gkey.json ./
CMD ["./server"]

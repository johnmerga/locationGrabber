FROM golang:alpine3.19 as dev
# binary will be $(go env GOPATH)/bin/air
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# or install it into ./bin/
# RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
WORKDIR /go/src/johnmerga/locationGrabber
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /server
CMD ["air", "-c", ".air.toml"]
# CMD ["go","run","./main.go"]

## Deploy
FROM alpine:latest as prod
# timezone
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip
COPY --from=dev /server ./
COPY ./gkey.json ./
CMD ["./server"]

## for future reference
FROM golang:alpine3.19 as dev
# air library for hot reloading in dev
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
# RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
WORKDIR /go/src/johnmerga/locationGrabber
COPY go.mod ./
COPY go.sum ./
# RUN chown -R botuser:botgroup ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o /server
CMD ["air", "-c", ".air.toml"]
# CMD ["go","run","./main.go"]

## Deploy
FROM scratch as prod
# timezone
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip
# since we are using scratch we need to add certs and user from alpine
COPY --from=dev /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=dev /server ./
CMD ["./server"]

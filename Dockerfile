FROM golang:1.23-alpine


# CHANGE_ME: should be path to your app
WORKDIR /go/src/github.com/IONOS-Forecast/gocast-development-leon
COPY . .

RUN go get -v ./...
RUN go build -o weather .

ENTRYPOINT ["./entrypoint.sh"]

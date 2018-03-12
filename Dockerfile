FROM golang:latest as builder
ADD . /go/src/github.com/vladimir-chernykh/call-tracker-backend
WORKDIR /go/src/github.com/vladimir-chernykh/call-tracker-backend
RUN go build github.com/vladimir-chernykh/call-tracker-backend/cmd/calltracker
CMD ["./calltracker"]
EXPOSE 80
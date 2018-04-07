FROM golang:latest as builder
ADD . /go/src/github.com/vladimir-chernykh/call-tracker-backend
WORKDIR /go/src/github.com/vladimir-chernykh/call-tracker-backend
RUN apt-get update && apt-get -y install ffmpeg  && go build -o main github.com/vladimir-chernykh/call-tracker-backend/cmd/calltracker
CMD ["./main"]
EXPOSE 80
FROM golang

COPY . /go/src/github.com/TinyKitten/TimelineServer
WORKDIR /go/src/github.com/TinyKitten/TimelineServer

RUN go get ./
RUN go build
	
EXPOSE 4000

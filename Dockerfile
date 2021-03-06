FROM golang:1.11
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app
EXPOSE 8080
RUN go get "github.com/gorilla/mux"
RUN go get "github.com/gomarkdown/markdown"
RUN go build
CMD ["./app"]

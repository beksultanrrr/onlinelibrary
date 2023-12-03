FROM golang:1.21.3

RUN go version
ENV GOPATH=/

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build -o onlinelibrary ./cmd/api

CMD ["./onlinelibrary"]
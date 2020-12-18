
#Golang Image
FROM golang


ARG key
ARG secret

ENV KEY=$key
ENV SECRET=$secret

WORKDIR /app

COPY ./backend /app

RUN go run *.go

# RUN go get -u github.com/golang/protobuf/protoc-gen-go
# RUN go get -u google.golang.org/grpc
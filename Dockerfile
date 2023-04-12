
FROM golang:alpine as build-env

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev
RUN mkdir /calcLab2
RUN mkdir -p /calcLab2/calcLab/server
RUN mkdir -p /calcLab2/grpc_api
RUN mkdir -p /calcLab2/calcLab_server

WORKDIR /calcLab2

#COPY ./grpc_api/calcLab.pb.go /calcLab2/grpc_api/calcLab.pb.go
#COPY ./grpc_api/calcLab_grpc.pb.go /calcLab2/grpc_api/calcLab_grpc.pb.go
#COPY ./main.go .
#COPY ./calcLab/server/server.go /calcLab2/calcLab/server/server.go
COPY ./ ./
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go build -o ./main .

RUN chmod "+x" ./main

EXPOSE 8083:8083

ENTRYPOINT [ "./main"]

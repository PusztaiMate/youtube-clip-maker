# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /clipper-server cmd/main.go

##
## Deploy
##
FROM ubuntu:21.04

ENV TZ=Europe/Budapest
ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /

RUN apt-get update && apt-get upgrade -y && apt-get install youtube-dl ffmpeg -y 

COPY --from=build /clipper-server /clipper-server

EXPOSE 8000

ENTRYPOINT ["/clipper-server"]
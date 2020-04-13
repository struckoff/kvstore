# pre-build stage for dependencies
FROM golang:alpine AS repo
WORKDIR /app
ENV GO111MODULE=on
COPY . /app
RUN go mod download
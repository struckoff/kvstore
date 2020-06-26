# build the kvstore binary
FROM golang:latest AS kvstore-builder
ADD . $GOPATH/src/github.com/struckoff/kvstore/
COPY store/cmd/config-docker.json /config-docker.json
WORKDIR $GOPATH/src/github.com/struckoff/kvstore/store/cmd
RUN go get ./
RUN CGO_ENABLED=0 go build -o /kvstore

# small packaging for kvstore
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
RUN mkdir /data
COPY --from=kvstore-builder /config-docker.json .
COPY --from=kvstore-builder /kvstore .
ENTRYPOINT /kvstore -c /config-docker.json
LABEL Name=kvstore Version=0.0.2
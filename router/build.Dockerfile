# build the kvrouter binary
FROM golang:latest AS kvrouter-builder
ADD . $GOPATH/src/github.com/struckoff/kvstore/
COPY router/cmd/config-docker.json /config-docker.json
WORKDIR $GOPATH/src/github.com/struckoff/kvstore/router/cmd
RUN go get ./
RUN CGO_ENABLED=0 go build -o /kvrouter

# small packaging for kvrouter
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
RUN mkdir /data
COPY --from=kvrouter-builder /config-docker.json .
COPY --from=kvrouter-builder /kvrouter .
ENTRYPOINT /kvrouter -c /config-docker.json
LABEL Name=kvrouter Version=0.0.2
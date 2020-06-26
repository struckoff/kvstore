FROM golang:latest
COPY .. $GOPATH/src/github.com/struckoff/kvstore/
COPY cmd/config-docker.json /config.json
WORKDIR $GOPATH/src/github.com/struckoff/kvstore/store/cmd
RUN go get ./
RUN CGO_ENABLED=0 go build -gcflags "all=-N -l" -o /kvstore
WORKDIR /
RUN go get github.com/go-delve/delve/cmd/dlv
EXPOSE 40000
ENTRYPOINT dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec ./kvstore
LABEL Name=kvstore-debug Version=0.0.1

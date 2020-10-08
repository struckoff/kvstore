FROM golang:latest
ADD . $GOPATH/src/github.com/struckoff/kvstore/
COPY router/cmd/config-docker.json /config.json
WORKDIR $GOPATH/src/github.com/struckoff/kvstore/router/cmd
RUN go get ./
RUN CGO_ENABLED=0 go build -gcflags "all=-N -l" -o /kvrouter
WORKDIR /
RUN go get github.com/go-delve/delve/cmd/dlv
EXPOSE 40000
ENTRYPOINT dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec ./kvrouter
LABEL Name=kvrouter-debug Version=0.0.2

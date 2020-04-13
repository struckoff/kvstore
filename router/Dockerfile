FROM golang
ENV GO111MODULE=on

RUN mkdir /app
RUN mkdir -p /var/lib/kvrouter
COPY . /app
WORKDIR /app/cmd
RUN go build -o kvrouter
ENTRYPOINT ./kvrouter -c config-docker.json

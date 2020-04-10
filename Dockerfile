FROM golang
ENV GO111MODULE=on


RUN mkdir /app
RUN mkdir -p /var/lib/kvstore
COPY . /app
WORKDIR /app/cmd
RUN go build -o kvstore
ENTRYPOINT ./kvstore -c config-docker.json

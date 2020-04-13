# build the kvstore binary
FROM root:latest AS kvstore-builder
WORKDIR /app/store/cmd
RUN go build -o kvstore

# small packaging for kvstore
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
RUN mkdir /var/lib/kvstore
COPY --from=kvstore-builder /app/store/cmd/config-docker.json /
COPY --from=kvstore-builder /app/store/cmd/kvstore /kvstore
ENTRYPOINT /kvstore -c config-docker.json
LABEL Name=kvstore Version=0.0.1
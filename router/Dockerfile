# build the kvrouter binary
FROM root:latest AS kvrouter-builder
WORKDIR /app/router/cmd
RUN go build -o kvrouter

# small packaging for kvrouter
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=kvrouter-builder /app/router/cmd/config-docker.json /config-docker.json
COPY --from=kvrouter-builder /app/router/cmd/kvrouter /kvrouter
ENTRYPOINT /kvrouter -c config-docker.json
LABEL Name=kvrouter Version=0.0.1
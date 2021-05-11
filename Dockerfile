FROM golang:alpine
WORKDIR /opt/app
COPY ./ ./
RUN go build
ENTRYPOINT ["./enry-server"]

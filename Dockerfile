FROM golang:1.25.0 as builder
WORKDIR /go/app
COPY . .
RUN go mod download
RUN go build -o ./src/build ./cmd/executor

FROM golang:1.25.0
WORKDIR /root/
COPY config/todoapp.yml ./config/todoapp.yml
COPY ./assets ./assets
COPY ./cmd/migration/scripts ./cmd/migration/scripts
COPY --from=builder /go/app/src/build .
EXPOSE 1212
ENTRYPOINT ["./build"]

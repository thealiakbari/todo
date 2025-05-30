FROM golang:1.23.3 as builder
WORKDIR /go/app
COPY . .
RUN go mod download
RUN go build -o ./src/build ./cmd/executor

FROM golang:1.23.3
WORKDIR /root/
COPY ./config/hichapp.yml ./config/hichapp.yml
COPY ./assets ./assets
COPY ./cmd/migration/scripts ./cmd/migration/scripts
COPY --from=builder /go/app/src/build .
EXPOSE 1212
ENTRYPOINT ["./build"]

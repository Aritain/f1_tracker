FROM golang:1.24.0-alpine3.21 as app-builder
WORKDIR /go/src/app
COPY . .
RUN go mod init client
RUN go mod tidy
RUN go build -o /f1_tracker

FROM alpine:3.16
COPY --from=app-builder /f1_tracker /f1_tracker
RUN apk add tzdata
RUN mkdir /data
ENTRYPOINT ["/f1_tracker"]

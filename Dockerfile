FROM golang:1.12-alpine as builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOPROXY=https://goproxy.cn
WORKDIR /project
COPY . .
RUN go build -v -a -installsuffix cgo -o /project/atom-server /project/main.go

FROM alpine
WORKDIR /hello-k8s
COPY --from=builder /project/atom-server .
COPY --from=builder /project/docs ./docs
ENTRYPOINT ["/hello-k8s/atom-server"]
CMD ["-c", "conf/config.yaml"]
EXPOSE 8080

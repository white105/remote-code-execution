FROM golang:latest AS builder
# Copy the code from the host and compile it
RUN go get github.com/golang/dep/cmd/dep
WORKDIR $GOPATH/src/remote-code-execution
COPY . ./
RUN dep ensure -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GODEBUG=http2debug=2 go build  -o /remote-code-execution

FROM alpine
COPY --from=builder /remote-code-execution ./
CMD ["./remote-code-execution"]
EXPOSE 1323

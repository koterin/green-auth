FROM golang:1.18
COPY ./ /password-manager/telegram
WORKDIR /password-manager/telegram
RUN go build .
ENTRYPOINT ["./main"]

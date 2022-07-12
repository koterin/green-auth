FROM golang:1.18
COPY ./ /green-auth/telegram
WORKDIR /green-auth/telegram
RUN go build .
ENTRYPOINT ["./main"]

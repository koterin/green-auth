FROM golang:1.18
COPY ./ /telegram
WORKDIR /telegram

RUN go build -o telegram
ENTRYPOINT ["./telegram"]

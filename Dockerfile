FROM golang:1.20.1

WORKDIR /app

COPY ./proxy ./proxy

CMD ["./proxy"]
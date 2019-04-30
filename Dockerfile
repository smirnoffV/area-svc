FROM golang:latest

WORKDIR /app

COPY ./bin/linux-app /app/linux-app

EXPOSE 10000

CMD ["/app/linux-app", "-addr=10000"]
FROM golang:1.22.1-bullseye

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE 8080

CMD sh deploy.sh

FROM golang:latest

LABEL maintainer="Stuko <avplaying@gmail.com>"

WORKDIR /app

COPY go.mod .

COPY go.sum .

COPY scarlet-419401.json /app/serviceAccountKey.json

RUN go mod download

COPY . .

ENV PORT 8000

ENV GOOGLE_APPLICATION_CREDENTIALS=/app/serviceAccountKey.json

RUN go build

CMD ["./scarlet_backend"]
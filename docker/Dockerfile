FROM golang:1.24.1-alpine

WORKDIR /app

COPY ../go.mod .
COPY ../go.sum .

RUN go mod download

COPY . .

RUN go build -o url_shortener ./cmd/url_shortener/

CMD [ "./url_shortener" ]
FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/uploads/submissions

EXPOSE 8080

CMD ["go", "run", "./cmd/server"]

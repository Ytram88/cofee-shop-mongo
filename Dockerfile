FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cofee-shop ./cmd/myapp/.

ENV MONGO_USER="cofeeStaff"
ENV MONGO_PASSWORD="cofeeAdmin"
ENV JWT_SECRET="secretJWT123"
ENV JWT_EXPIRATION_IN_SECONDS=60*120

EXPOSE 8080

CMD ["./cofee-shop"]



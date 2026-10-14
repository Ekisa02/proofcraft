FROM golang:1.20-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN apk add --no-cache git
RUN go mod download
COPY . .
RUN go build -o proofcraft main.go
EXPOSE 8080
CMD ["./proofcraft"]

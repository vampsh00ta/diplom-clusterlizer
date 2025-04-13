FROM golang:latest
WORKDIR /app
COPY ./go.mod ./go.sum ./
COPY cmd/publicapi/preprod.yaml .
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/publicapi
CMD ["./app","-f","preprod.yaml"]
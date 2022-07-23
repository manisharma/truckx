FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
WORKDIR /root
COPY . .
RUN go mod download 
COPY . .
WORKDIR /root/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o truckx .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /root/cmd/truckx .
COPY --from=builder /root/.env .env  
EXPOSE 8080
ENTRYPOINT [ "./truckx" ]
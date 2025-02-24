# Use Golang image
FROM golang:1.19-alpine AS builder


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . . 


RUN if [ ! -f .env ]; then echo "PORT=8080" > .env; fi


RUN go build -o main ./cmd/*


FROM alpine:latest  

WORKDIR /root/


RUN apk add --no-cache tzdata  


ENV TZ=Asia/Kolkata  


COPY --from=builder /app/main .
COPY --from=builder /app/.env .  


EXPOSE 8080

# Run the application
CMD ["/root/main"]
FROM golang:1.19
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o receipt-processor
EXPOSE 8080
CMD ["/app/receipt-processor"]
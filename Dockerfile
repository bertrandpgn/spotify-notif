FROM golang:1.20 AS build

WORKDIR /app
COPY . /app
RUN go mod download \
 && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch

WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE 80
CMD ["./main"]

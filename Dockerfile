FROM golang:latest as builder
WORKDIR /go/src/github.com/alekssaul/urlsigner
COPY . .
RUN mkdir -p /app
RUN CGO_ENABLED=0 GOOS=linux go test . 
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/urlsigner .


FROM alpine:latest
RUN apk update ;  apk add --no-cache ca-certificates ; update-ca-certificates ; mkdir /app
WORKDIR /app
COPY --from=builder /app/urlsigner /app/urlsigner
CMD /app/urlsigner
EXPOSE 8080
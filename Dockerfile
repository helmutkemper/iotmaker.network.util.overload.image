FROM golang:alpine3.12 as builder

RUN mkdir /app
RUN chmod 700 /app

COPY . /app

WORKDIR /app

# import golang packages to be used inside image "scratch"
ARG CGO_ENABLED=0
# RUN go mod init github.com/kempertrasdesclub/queue
RUN go build -ldflags="-w -s" -o /app/main /app/main.go

FROM scratch
#FROM golang:alpine3.12

COPY --from=builder /app/main .

EXPOSE 8080,8000
CMD ["/main"]

FROM golang:alpine as builder
WORKDIR /app

COPY *.go ./
RUN go mod init mod
RUN go mod tidy && go mod download
RUN go build -o /appbin/gnu_chess

FROM alpine
COPY --from=builder /appbin /
COPY ./index.html ./favicon.ico /
COPY ./images/* /images/
COPY ./client/* /client/
EXPOSE 8080
CMD ["/gnu_chess"]

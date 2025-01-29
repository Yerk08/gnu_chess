FROM golang:1.19 as builder
WORKDIR /app
COPY *.go ./
RUN go mod init mod && go mod tidy && go mod download && go build -o /appbin/gnu_chess

FROM alpine
COPY --from=builder /appbin /
EXPOSE 8080
CMD ["/gnu_chess"]

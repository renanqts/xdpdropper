FROM golang:1.18 as builder

WORKDIR /go/src
COPY . /go/src

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o xdpdropper ./cmd

FROM gcr.io/distroless/base-debian11 AS runtime

COPY --from=builder /go/src/xdpdropper ./
CMD ["./xdpdropper"]  

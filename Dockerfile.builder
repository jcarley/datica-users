# vi:syntax=dockerfile
FROM golang:1.9.2-stretch as builder

ARG binary

WORKDIR /go/src/github.com/jcarley/datica-users
COPY . ./
RUN CGO_ENABLED=0 go install -a -tags netgo -ldflags '-extldflags "-static"'
RUN ldd /go/bin/"$binary" | grep -q "not a dynamic executable"

FROM scratch
ARG binary
COPY --from=builder /go/bin/"$binary" /"$binary"
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
EXPOSE "3000"
CMD ["/datica-users", "server"]

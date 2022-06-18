FROM golang:1.15 as builder
WORKDIR /go/src/istio.io/jwt-server/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o jwt-server main.go

FROM gcr.io/distroless/static-debian10@sha256:4433370ec2b3b97b338674b4de5ffaef8ce5a38d1c9c0cb82403304b8718cde9 as distroless

WORKDIR /bin/
# copy the jwt-server binary to a separate container based on BASE_DISTRIBUTION
COPY --from=builder /go/src/istio.io/jwt-server .
ENTRYPOINT [ "/bin/jwt-server" ]
EXPOSE 8000

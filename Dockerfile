FROM golang:alpine as app-builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -tags timetzdata -o /go/bin/iamme

FROM scratch
COPY --from=app-builder /go/bin/iamme /iamme
COPY .env /.env
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
HEALTHCHECK NONE
ENTRYPOINT [ "/iamme" ]
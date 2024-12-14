FROM golang:alpine as app-builder
WORKDIR /go/src/app
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -tags timetzdata -o /go/bin/iamme

FROM scratch
COPY --from=app-builder /go/bin/iamme /iamme
COPY .env* /.env
# ref: https://github.com/charmbracelet/log/issues/90
ENV CI=true
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
HEALTHCHECK NONE
ENTRYPOINT [ "/iamme" ]
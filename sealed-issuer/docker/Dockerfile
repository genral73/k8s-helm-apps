FROM alpine:3.14.0 as download
RUN apk add curl \
    && curl https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.16.0/kubeseal-linux-amd64 -L --output /kubeseal \
    && chmod +x /kubeseal

FROM golang:1.12 as build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:3.14.0
LABEL maintainer="Dafaallah <eng.dafaallah@gmail.com>"
COPY ./kubeseal/ ./kubeseal
COPY --from=download /kubeseal ./kubeseal/bin
COPY --from=build /app/sealedissuer sealedissuer
COPY ./static/ ./static

EXPOSE 8080

ENTRYPOINT [ "./sealedissuer" ]

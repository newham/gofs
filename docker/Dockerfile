FROM golang:1.10-alpine

# RUN apk add --no-cache ca-certificates

COPY copy /gofs

EXPOSE 8087

WORKDIR /gofs

ENTRYPOINT /gofs/gofs

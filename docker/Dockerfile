FROM debian:latest

# RUN apk add --no-cache ca-certificates

COPY copy /gofs

EXPOSE 8087

WORKDIR /gofs

ENTRYPOINT /gofs/gofs

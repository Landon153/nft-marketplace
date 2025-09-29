FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates

WORKDIR /app/
ADD ./app /app/
# ADD ./zoneinfo.zip /usr/local/go/lib/time/
ADD ./assets /app/assets
ADD ./marketplace.abi /app/marketplace.abi
EXPOSE 3000
ENTRYPOINT ["./app"]
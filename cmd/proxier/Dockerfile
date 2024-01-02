FROM ubuntu:latest

LABEL org.opencontainers.image.source='https://github.com/Jacobbrewer1/proxy'
LABEL org.opencontainers.image.description="This is a proxy app written in Go that proxies requests to a given URL."
LABEL org.opencontainers.image.licenses='GNU General Public License v3.0'

WORKDIR /proxy

# Copy the binary from the build
COPY ./bin/app /proxy/app

RUN ["chmod", "+x", "./app"]

ENTRYPOINT ["/proxy/app"]

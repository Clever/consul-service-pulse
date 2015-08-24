FROM gliderlabs/alpine:3.2
ENTRYPOINT ["/bin/consul-service-pulse"]

COPY . /go/src/github.com/Clever/consul-service-pulse
RUN apk-install -t build-deps go git \
    && cd /go/src/github.com/Clever/consul-service-pulse \
    && export GOPATH=/go \
    && go get \
    && go build -ldflags "-X main.Version $(cat VERSION)" -o /bin/consul-service-pulse \
    && rm -rf /go \
    && apk del --purge build-deps

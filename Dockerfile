FROM alpine:3.4
RUN apk add --no-cache ca-certificates mailcap
ADD kubernetes-honeycomb-agent /
ENTRYPOINT ["/kubernetes-honeycomb-agent"]

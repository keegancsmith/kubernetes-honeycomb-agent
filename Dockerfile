FROM scratch
ADD kubernetes-honeycomb-agent /
ENTRYPOINT ["/kubernetes-honeycomb-agent"]

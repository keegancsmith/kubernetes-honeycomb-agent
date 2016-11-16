# kubernetes-honeycomb-agent

`kubernetes-honeycomb-agent` exports events from the
(Kubernetes)[http://kubernetes.io/] API into
[honeycomb](https://honeycomb.io/). It is like piping
`kubectl --all-namespaces=true get events --watch` into honeycomb.

**Status: experimental**

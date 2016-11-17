# kubernetes-honeycomb-agent

`kubernetes-honeycomb-agent` exports events from the
[Kubernetes](http://kubernetes.io/) API into
[honeycomb](https://honeycomb.io/). It is like piping
`kubectl --all-namespaces=true get events --watch` into honeycomb.

To get started, just run

```
$ kubectl run kubernetes-honeycomb-agent \
    --image=keegancsmith/kubernetes-honeycomb-agent:v0.1 \
    --env='HONEYCOMB_TEAM=<REDACTED>'
```

with `<REDACTED>` as your write key / honeycomb team identifier. You should
start getting a dataset called `kubernetes-events`. To put the events into a
different dataset, just set the environment variable `HONEYCOMB_DATASET`.

**Status: experimental**

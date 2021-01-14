# fmz

Conformance tests for the venturemark api.



```
docker run --rm -p 127.0.0.1:6379:6379 redis
```

```
go build && echo ready && ./apiserver daemon
```

```
go test ./... --tags conformance
```

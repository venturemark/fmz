# fmz

Conformance tests for the venturemark api.



```
docker run --rm -p 127.0.0.1:6379:6379 redis
```

```
docker run --net host --rm -ti redis redis-cli
```

```
go build && ./apiserver daemon
go build && ./apiworker daemon
```

```
go test ./... -tags conformance
```

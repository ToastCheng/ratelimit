# Rate Limit

A middleware version is in https://github.com/ToastCheng/ratelimiter

## Quick start
```bash
# build
./scripts/build.sh

# build image and run in docker container
./scripts/run_container.sh

# directly run in host.
./scripts/run_host.sh

# make query
curl http://localhost:8080
```

## 1. Structure
```
.
├── README.md
├── docker
├── scripts
└── src
```

`README.md` read file.

`docker` docker file.

`scripts` scripts to test, build, and run.

`src` source code of rate limit server.

## 2. Source
note: `go.mod` and `go.sum` are ignored.

```
.
├── config.go
├── config_test.go
├── impl.go
├── impl_test.go
├── record.go
├── record_test.go
└── server
    └── main.go
```

`config.go` defines the necessary configuration for ratelimit server.

`config_test.go` test of `config.go`.

`impl.go` the main handler logic of ratelimit server.

`impl_test.go` test of the ratelimit server.

`record.go` the data model of request record, see [3. Description]().

`record_test.go` test of `record.go`

`server/` the entry point of the server.

## 3. Description

The `ratelimit` package implements a service with rate limiting functionality.

For simplicity, 
1. `ratelimit` server uses in-memory data structure to keep track of the number of requests each user made within 60 second.
2. HTTP server directly serves all incoming HTTP traffic to port `8080`, i.e., no routing. 

### `RateLimitHandler`

`RateLimitHandler` is the struct that handles the request. It implements `ServeHTTP` so that it can be passed to `http.Server` as a `http.Handler`.

It also has fields `limit` and `window`, both default to `60` (second). Representing the limit of number of request, the time interval which the number of request should be count, respectively.

```go
// RateLimitHandler implement ServeHTTP,
// handles the incoming request and performs rate limiting.
type RateLimitHandler struct {
	records map[string]*Record
	limit   int
	window  int
}
```

Inside the `RateLimitHandler`, a hashmap `records` is used to record the timestamp of each user's request.

```go
// Record stores the request's timestamp corresponding to a certain IP.
type Record struct {
	mtx       sync.Mutex
	timestamp []int64
}
```

Every time a request comes in, `RateLimitHandler` will check if this IP address exists in `records`. If not, initialize a `Record` for this IP.

Then `RateLimitHandler` will try to add a timestamp to the `Record`'s `timestamp` array by calling `Record`'s `Add` method.

`func Record.Add(limit, window int) (int, error)` will first aquire a lock to make sure the below operation is thread-safe. Noted that the lock is in `Record` level, meaning that request from different IP will not be blocked.

After that, it then calculate the start time `start`, which is the unix timestamp of `window` second before now. And find the first value in `timestamp` that is greater than `start` by binary search.

Now `Record` will know how many request is made within the past `window` second. If it exceeds `limit`, return an error. Else, return the number of request within `window` second + 1 (the current request).

On the other hand, the records that are before `start` are no longer needed, `timestamp` will be sliced to save space.

## 4. Test

All test can be found in the test files.

### `config_test.go`
- `TestConfigPort` test invalid port value.
- `TestConfigLimit` test invalid limit value.
- `TestConfigWindow` test invalid window value.

### `record_test.go`
- `TestNewRecord` test creating a new record.
- `TestAdd` test adding timestamp in record until it exceeds the limit.

### `impl_test.go`
- `TestQueryOne` test making one query to server.
- `TestQueryRepeat` test making request in a fair time interval so that the limit will not be exceeded.
- `TestQueryFromDifferentIP` test make request from different IP, so that the record is saved on different hashmap key, and the limit will not be exceeded.
- `TestQueryConcurrentRepeat` test making request concurrently, only the `limit` requests will be served, the rest will get `429`.
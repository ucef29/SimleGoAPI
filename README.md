# SimleGoAPI

```
### cache benchMark
# cd to handlers dir and run:

go test -bench .
goos: linux
goarch: amd64
pkg: SimpleGoAPI/handlers
BenchmarkGetAllNonCached-4          5000            424201 ns/op
BenchmarkGetAllCached-4           200000              6941 ns/op
PASS
ok      SimpleGoAPI/handlers    7.940s
```
```
# using Echo framework to implement 3 middlewares

func serveCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if cache.Serve(c.Response(), c.Request()) {
			return nil
		}
		return next(c)
	}
}

func cacheResponse(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Writer = cache.NewWriter(c.Response().Writer, c.Request())
		return next(c)
	}
}

func auth(username, password string, c echo.Context) (bool, error) {
	if username == "joe" && password == "secret" {
		return true, nil
	}
	return false, nil
}

```

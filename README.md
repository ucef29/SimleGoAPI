# SimleGoAPI

```
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
# using Echo framework to implementing 3 middlewares

func cacheResponse(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Writer = cache.NewWriter(c.Response().Writer, c.Request())
		return next(c)
	}
}

func usersOptions(c echo.Context) error {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodOptions}
	c.Response().Header().Set("Allow", strings.Join(methods, ","))
	return c.NoContent(http.StatusOK)
}

func auth(username, password string, c echo.Context) (bool, error) {
	if username == "joe" && password == "secret" {
		return true, nil
	}
	return false, nil
}

```

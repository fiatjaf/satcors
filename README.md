# satcors

Super light-weight CORS-anywhere proxy implemented in just a few lines of Go.

## Usage

Just pass your desired request the the root endpoint like this:

`https://pan.cors/?url=https%3A%2F%2Fmichaljanocko.com`

Please include the scheme (`https` or `http`) to qualify as a valid URL. satcors can't proxy other protocols so it has to check.

If you need to request something from a server that checks the `Referer` header, you can set the `&referer=` URL parameter and satcors will send the request with the correct referer.

Also, don't forget to query encode the address because some URI implementations merge neighboring slashes into one if they're part of the path or query (e.g. `https://` wouldn't work and would come out as `https:/` so that why we have to encode it as a param; this is the safest way IMO)

For custom response headers, use the snippet below if you're using the handler functions
```go
func HandleProxyWith(origin string, credentials string)
```
or set these environment for the binary:
- `ALLOW_ORIGIN` to whatever you want your *Access-Control-Allow-Origin* header to be (default is `*`)
- `ALLOW_CREDENTIALS` to either `true` or `false` (default is `true`)

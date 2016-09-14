# Go-Shortener

A small url shortener written in Go.

## Requirements

The only requirements are docker and docker-compose, along with GNU make.

## Running (development)

To build the container, run `make build`.

To run in development you can simply run `make run`.  This will download a Redis image, start it, and properly run the container itself.  This will build the `go-shortener` container if it does not already exist.

## Running (without `docker-compose`)

To run the container without `docker-compose` (as would be appropriate in production), the `REDIS_URL` environment variable must be set to the host and port of the Redis instance.

## Running the Tests

Run the tests with `make run`.  This builds the image and compiles the tests before running them.

## Endpoints

### GET /:shortUrl

Retrieve `shortUrl` from redis using the key `url:{shortUrl}`, which contains the original, unshortened url.  This endpoint returns a 301 Moved Permanently redirect to the original url, and returns a 404 Not Found if `shortUrl` does not exist in Redis.  If the url exists, the total and daily hits count will be incremented (further described below).

Example:
```bash
$ curl -XGET http://`docker-machine ip`:8080/RNFIp -v
*   Trying 192.168.99.100...
* Connected to 192.168.99.100 (192.168.99.100) port 8080 (#0)
> GET /RNFIp HTTP/1.1
> Host: 192.168.99.100:8080
> User-Agent: curl/7.43.0
> Accept: */*
>
< HTTP/1.1 301 Moved Permanently
< Location: http://lmgtfy.com
< Date: Wed, 14 Sep 2016 13:17:12 GMT
< Content-Length: 52
< Content-Type: text/html; charset=utf-8
<
<a href="http://lmgtfy.com">Moved Permanently</a>.
```

### POST /create

Create a short link from a json payload `{"Url": "myVerySpecialSite.com"}`.  The url will be hashed using a CRC32 checksum, base 62-encoded.  The result will be a shortlink, which is guaranteed to be a string with a maximum length of six.  The original url will be stored in redis using the key `url:{shortUrl}` and the short link will be returned to the user as `{"Url": "{shortUrl}"}`.

Example:

```bash
$ curl -XPOST http://`docker-machine ip`:8080/create -d '{"Url": "http://lmgtfy.com"}' -v
*   Trying 192.168.99.100...
* Connected to 192.168.99.100 (192.168.99.100) port 8080 (#0)
> POST /create HTTP/1.1
> Host: 192.168.99.100:8080
> User-Agent: curl/7.43.0
> Accept: */*
> Content-Length: 28
> Content-Type: application/x-www-form-urlencoded
>
* upload completely sent off: 28 out of 28 bytes
< HTTP/1.1 200 OK
< Date: Wed, 14 Sep 2016 13:16:36 GMT
< Content-Length: 15
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host 192.168.99.100 left intact
{"Url":"RNFIp"}
```

### GET /stats/:shortUrl

Fetch total and daily hits for the past year for `shortUrl`. Hits are stored as a hash in redis under the key `hits:{shortUrl}`.  The hash fields are integers representing of a day within the year, between 1 and 366, and the values are the number of hits on that day.  Numbers higher than the current year day represent that day in the previous year.  There is also a `Total` field to represent the total number of hits for a short url.  The structure of the returned payload is `{"Count": {totalHits}, "Days": {"{day1}": {hitsDay1}, "{day2}": {hitsDay2}, ...}}`.

Example:

```bash
$ curl -XGET http://`docker-machine ip`:8080/stats/RNFIp -v
*   Trying 192.168.99.100...
* Connected to 192.168.99.100 (192.168.99.100) port 8080 (#0)
> GET /stats/RNFIp HTTP/1.1
> Host: 192.168.99.100:8080
> User-Agent: curl/7.43.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Wed, 14 Sep 2016 17:17:48 GMT
< Content-Length: 45
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host 192.168.99.100 left intact
{"Count":7,"Days":{"2016-09-14T00:00:00Z":7}}
```

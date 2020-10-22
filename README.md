# go-http-mbtiles

Go HTTP handler for serving MBTiles databases.

## Example

```
package main

import (
	"github.com/aaronland/go-http-mbtiles"
	"net/http"
	"regexp"
)

func main() {

	tiles_source := "/path/to/folder/containing/mbtiles/"
	tiles_extenion := ".db"
	tiles_path := "/tiles"
	tiles_pattern := `/tiles/([a-z-]+)/(\d+)/(\d+)/(\d+)\.([a-z]+)$`
	
	tiles_re, _ := regexp.Compile(tiles_pattern)

	tiles_opts := &mbtiles.MBTilesHandlerOptions{
		Root:         tiles_source,
		TilesPattern: tiles_re,
	}

	tiles_handler, _ := mbtiles.MBTilesHandler(tiles_opts)

	mux := http.NewServeMux()
	mux.Handle(tiles_path, tiles_handler)

	// serve mux here
}
```

_Error handling omitted for brevity._

## Tools

```
$> make cli
go build -mod vendor -o bin/server cmd/server/main.go
```

### server

```
$> ./bin/server -h
  -server-uri string
    	A valid aaronland/go-http-server URI string. (default "http://localhost:8080")
  -tiles-extension string
    	The extension (minus the leading dot) for your MBTiles databases. (default ".mbtiles")
  -tiles-path string
    	The relative path to serve tiles from. (default "/tiles/")
  -tiles-pattern string
    	A valid Go language regular expression for validating requests. The pattern needs to return five values: name of the MBTiles file, Z, X and Y tile values and a file extension used to determine content type. (default "/tiles/([a-z-]+)/(\\d+)/(\\d+)/(\\d+)\\.([a-z]+)$")
  -tiles-source string
    	Path to the directory containing your MBTiles databases.
```	

For example:

```
$> ./bin/server -tiles-source /usr/local/mbtiles/
2020/10/21 21:53:00 Listening on http://localhost:8080
```

#### Lambda

_Please write me._

## See also

* https://github.com/mattn/go-sqlite3
* https://github.com/aaronland/go-http-mbtiles
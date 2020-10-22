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
	tiles_extension := ".db"
	tiles_path := "/tiles"
	tiles_pattern := `/tiles/([a-z-]+)/(\d+)/(\d+)/(\d+)\.([a-z]+)$`
	
	tiles_re, _ := regexp.Compile(tiles_pattern)

	tiles_opts := &mbtiles.MBTilesHandlerOptions{
		Root:         tiles_source,
		Extension:    ties_extension,
		Pattern: tiles_re,
	}

	tiles_handler, _ := mbtiles.MBTilesHandler(tiles_opts)

	mux := http.NewServeMux()
	mux.Handle(tiles_path, tiles_handler)

	// serve mux here
}
```

_Error handling omitted for brevity._

## See also

* https://github.com/mattn/go-sqlite3
* https://github.com/aaronland/go-mbtiles-server
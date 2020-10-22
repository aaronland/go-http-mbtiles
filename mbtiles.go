package mbtiles

import (
	"database/sql"
	"fmt"
	"github.com/aaronland/go-mimetypes"
	_ "github.com/mattn/go-sqlite3"
	"math"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type MBTilesHandlerOptions struct {
	Root         string
	Extension    string
	TilesPattern *regexp.Regexp
}

func MBTilesHandler(opts *MBTilesHandlerOptions) (http.Handler, error) {

	db_conns := make(map[string]*sql.DB)
	mu := new(sync.RWMutex)

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		path := req.URL.Path

		m := opts.TilesPattern.FindStringSubmatch(path)

		if len(m) != 6 {
			http.Error(rsp, "Invalid tile URI", http.StatusBadRequest)
			return
		}

		tileset := m[1]
		str_z := m[2]
		str_x := m[3]
		str_y := m[4]
		format := m[5]

		z, err := strconv.ParseInt(str_z, 10, 64)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		x, err := strconv.ParseInt(str_x, 10, 64)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		y, err := strconv.ParseInt(str_y, 10, 64)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		t := mimetypes.TypesByExtension(format)

		if len(t) == 0 {
			http.Error(rsp, "Invalid or unrecognized format", http.StatusBadRequest)
			return
		}

		content_type := t[0]

		mu.RLock()

		db, ok := db_conns[tileset]

		mu.RUnlock()

		if !ok {

			ext := strings.TrimPrefix(opts.Extension, ".")

			db_fname := fmt.Sprintf("%s.%s", tileset, ext)
			db_path := filepath.Join(opts.Root, db_fname)

			db_dsn := fmt.Sprintf("file:%s?cache=shared&mode=ro", db_path)
			conn, err := sql.Open("sqlite3", db_dsn)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			mu.Lock()

			db_conns[tileset] = conn
			db = conn

			mu.Unlock()
		}

		// https://github.com/TileStache/TileStache/blob/master/TileStache/MBTiles.py#L169

		fl_z := float64(z)
		mb_y := (int64(math.Pow(2.0, fl_z)) - 1) - y

		q := "SELECT tile_data FROM tiles WHERE zoom_level = ? AND tile_column = ? AND tile_row = ?"
		row := db.QueryRowContext(ctx, q, z, x, mb_y)

		var body []byte

		err = row.Scan(&body)

		if err != nil {

			if err == sql.ErrNoRows {
				http.Error(rsp, "Not found", http.StatusNotFound)
				return
			}

			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-type", content_type)
		rsp.Write(body)
		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}

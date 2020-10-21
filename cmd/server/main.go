package main

import (
	"context"
	"github.com/aaronland/go-http-mbtiles"
	"github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-flags/flagset"
	"log"
	"net/http"
	"regexp"
)

func main() {

	fs := flagset.NewFlagSet("prettymaps")

	server_uri := fs.String("server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI string.")
	tiles_source := fs.String("tiles-source", "", "...")
	tiles_path := fs.String("tiles-path", "/tiles/", "...")
	tiles_pattern := fs.String("tiles-pattern", `/tiles/([a-z-]+)/(\d+)/(\d+)/(\d+)\.([a-z]+)$`, "...")

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVarsWithFeedback(fs, "MBTILES", true)

	if err != nil {
		log.Fatalf("Failed to set flags, %v", err)
	}

	ctx := context.Background()

	tiles_re, err := regexp.Compile(*tiles_pattern)

	if err != nil {
		log.Fatalf("Failed to compile tiles pattern, %v", err)
	}

	tiles_opts := &mbtiles.MBTilesHandlerOptions{
		Root:         *tiles_source,
		TilesPattern: tiles_re,
	}

	tiles_handler, err := mbtiles.MBTilesHandler(tiles_opts)

	if err != nil {
		log.Fatalf("Failed to create MBTiles handler, %v", err)
	}

	mux := http.NewServeMux()

	mux.Handle(*tiles_path, tiles_handler)

	s, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Failed to create server, %v", err)
	}

	log.Printf("Listening on %s", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}

}

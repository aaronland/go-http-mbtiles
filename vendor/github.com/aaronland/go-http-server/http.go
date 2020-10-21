package server

// https://medium.com/@simonfrey/go-as-in-golang-standard-net-http-config-will-break-your-production-environment-1360871cb72b

// https://ieftimov.com/post/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/

// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
// https://blog.cloudflare.com/exposing-go-on-the-internet/

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func init() {
	ctx := context.Background()
	RegisterServer(ctx, "http", NewHTTPServer)
}

type HTTPServer struct {
	Server
	url         *url.URL
	http_server *http.Server
	cert        string
	key         string
}

func NewHTTPServer(ctx context.Context, uri string) (Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	u.Scheme = "http"

	read_timeout := 2 * time.Second
	write_timeout := 10 * time.Second
	idle_timeout := 15 * time.Second
	header_timeout := 2 * time.Second

	q := u.Query()

	if q.Get("read_timeout") != "" {

		to, err := strconv.Atoi(q.Get("read_timeout"))

		if err != nil {
			return nil, err
		}

		read_timeout = time.Duration(to) * time.Second
	}

	if q.Get("write_timeout") != "" {

		to, err := strconv.Atoi(q.Get("write_timeout"))

		if err != nil {
			return nil, err
		}

		write_timeout = time.Duration(to) * time.Second
	}

	if q.Get("idle_timeout") != "" {

		to, err := strconv.Atoi(q.Get("idle_timeout"))

		if err != nil {
			return nil, err
		}

		idle_timeout = time.Duration(to) * time.Second
	}

	if q.Get("header_timeout") != "" {

		to, err := strconv.Atoi(q.Get("header_timeout"))

		if err != nil {
			return nil, err
		}

		header_timeout = time.Duration(to) * time.Second
	}

	tls_cert := q.Get("cert")
	tls_key := q.Get("key")

	if (tls_cert != "") && (tls_key != "") {

		_, err = os.Stat(tls_cert)

		if err != nil {
			return nil, err
		}

		_, err = os.Stat(tls_key)

		if err != nil {
			return nil, err
		}

		u.Scheme = "https"

	} else if (tls_cert != "") && (tls_key == "") {
		return nil, errors.New("Missing TLS key parameter")
	} else if (tls_key != "") && (tls_key == "") {
		return nil, errors.New("Missing TLS cert parameter")
	} else {
		// pass
	}

	srv := &http.Server{
		Addr:              u.Host,
		ReadTimeout:       read_timeout,
		WriteTimeout:      write_timeout,
		IdleTimeout:       idle_timeout,
		ReadHeaderTimeout: header_timeout,
	}

	server := HTTPServer{
		url:         u,
		http_server: srv,
		cert:        tls_cert,
		key:         tls_key,
	}

	return &server, nil
}

func (s *HTTPServer) Address() string {
	return s.url.String()
}

func (s *HTTPServer) ListenAndServe(ctx context.Context, mux *http.ServeMux) error {

	idleConnsClosed := make(chan struct{})

	go func() {

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.

		err := s.http_server.Shutdown(context.Background())

		if err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}

		close(idleConnsClosed)
	}()

	s.http_server.Handler = mux

	var err error

	if s.cert != "" && s.key != "" {
		err = s.http_server.ListenAndServeTLS(s.cert, s.key)
	} else {
		err = s.http_server.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	return nil
}

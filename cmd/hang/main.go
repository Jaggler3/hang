package main

import (
	"fmt"
	"os"

	render "hang.sh/internal/render"
	"hang.sh/internal/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	srv, err := server.New(render.TeaHandler)
	if err != nil {
		return err
	}

	err = srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

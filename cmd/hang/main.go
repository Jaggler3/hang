package main

import (
	"fmt"
	"os"

	"hang.sh/internal/server"
	"hang.sh/internal/ui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	srv, err := server.New(ui.TeaHandler)
	if err != nil {
		return err
	}

	err = srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

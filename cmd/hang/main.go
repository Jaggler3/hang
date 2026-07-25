package main

import (
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/ssh"
	"hang.sh/internal/data"
	"hang.sh/internal/entity"
	"hang.sh/internal/render"
	"hang.sh/internal/server"
	"hang.sh/internal/world"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	w := world.NewDefaultWorld()

	if err := data.LoadWorld(w); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not load world: %v\n", err)
	}

	go func() {
		for {
			time.Sleep(30 * time.Second)
			if err := data.SaveWorld(w); err != nil {
				fmt.Fprintf(os.Stderr, "warning: could not save world: %v\n", err)
			}
		}
	}()

	handler := func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
		ctx := sess.Context()
		p := ctx.Value(data.PlayerKey).(*entity.Player)
		ctx.SetValue(data.WorldKey, w)
		m := render.NewGameModel(p, w)
		return m, nil
	}

	srv, err := server.New(handler)
	if err != nil {
		return err
	}

	return srv.ListenAndServe()
}

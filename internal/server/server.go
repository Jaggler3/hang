package server

import (
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"
	"hang.sh/internal/player"
)

type contextKey string

const playerKey contextKey = "player"

// Initialize a wish server with configurations
func New(handler bubbletea.Handler) (*ssh.Server, error) {
	server, err := wish.NewServer(
		wish.WithAddress(":2222"),
		wish.WithHostKeyPath(".ssh/hang_host_key"),
		wish.WithPublicKeyAuth(keyAuthHandler),
		wish.WithMiddleware(
			bubbletea.Middleware(handler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func keyAuthHandler(ctx ssh.Context, key ssh.PublicKey) bool {
	// use 'key' in the future to persist player data in db
	ctx.SetValue(playerKey, player.NewPlayer())
	return true
}

package server

import (
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"
	"hang.sh/internal/data"
	"hang.sh/internal/entity"
	"hang.sh/internal/world"
)

func New(handler bubbletea.Handler) (*ssh.Server, error) {
	server, err := wish.NewServer(
		wish.WithAddress(":2222"),
		wish.WithHostKeyPath(".ssh/hang_host_key"),
		wish.WithPublicKeyAuth(keyAuthHandler),
		wish.WithMiddleware(
			cleanupMiddleware(bubbletea.Middleware(handler)),
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
	ctx.SetValue(data.PlayerKey, entity.NewPlayer(""))
	return true
}

func cleanupMiddleware(next wish.Middleware) wish.Middleware {
	return func(handler ssh.Handler) ssh.Handler {
		return func(sess ssh.Session) {
			ctx := sess.Context()
			p := ctx.Value(data.PlayerKey).(*entity.Player)
			defer func() {
				if w, ok := ctx.Value("world").(*world.World); ok {
					w.RemovePlayer(p.ID)
				}
			}()
			next(handler)(sess)
		}
	}
}

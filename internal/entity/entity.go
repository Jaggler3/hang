package entity

import (
	"crypto/rand"
	"fmt"

	"hang.sh/internal/shared"
	"hang.sh/internal/sprite"
)

type Entity struct {
	ID        string
	X, Y      int
	Direction shared.Direction
	Moving    bool
	Sprite    *sprite.Sprite
}

type Player struct {
	Entity
	Name    string
	HeadChar string
	HeadFG  string
	HeadBG  string
}

func NewPlayer(name string) *Player {
	return &Player{
		Entity: Entity{
			ID:        randID(),
			X:         0,
			Y:         0,
			Direction: shared.Down,
		},
		Name:    name,
		HeadChar: "@",
		HeadFG:  "#ffffff",
		HeadBG:  "",
	}
}

type NPC struct {
	Entity
	Name         string
	DialogID     string
	InteractText string
}

func randID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

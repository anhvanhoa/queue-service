package goidS

import (
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type GoId struct {
	alphabet string
	size     int
}

func NewGoId() *GoId {
	return &GoId{
		alphabet: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		size:     18,
	}
}

func (g *GoId) Gen() string {
	return gonanoid.MustGenerate(g.alphabet, g.size)
}

func (g *GoId) GenWithLength(length int) string {
	return gonanoid.MustGenerate(g.alphabet, length)
}

func (g *GoId) SetAlphabet(alphabet string) {
	g.alphabet = alphabet
}
func (g *GoId) SetSize(size int) {
	g.size = size
}

func (g *GoId) NewUUID() string {
	return uuid.New().String()
}

package console

import (
	"log"

	"github.com/rivo/tview"
)

type Screen interface {
	GetBox() *tview.Box
	GetDefaultFocus() tview.Primitive
	Build()
}

type ScreenBind struct {
	name   string
	screen Screen
}

func (sb *ScreenBind) GetName() string {
	return sb.name
}

func (sb *ScreenBind) GetScreen() Screen {
	return sb.screen
}

func NewScreenBind(name string, screen Screen) ScreenBind {
	sb := ScreenBind{
		name:   name,
		screen: screen,
	}

	return sb
}

type ScreenCollection struct {
	binds []ScreenBind
}

func NewScreenCollection() ScreenCollection {
	sc := ScreenCollection{}
	sc.binds = make([]ScreenBind, 0)

	return sc
}

func (sc *ScreenCollection) AddScreen(name string, screen Screen) {
	sb := NewScreenBind(name, screen)
	sc.binds = append(sc.binds, sb)
}

func (sc *ScreenCollection) GetScreen(name string) *Screen {
	log.Println("all binds", sc.binds)
	for _, sb := range sc.binds {
		log.Println("Bind name", sb.name)
		if sb.GetName() == name {
			return &sb.screen
		}
	}

	return nil
}

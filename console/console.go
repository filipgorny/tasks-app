package console

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Console struct {
	screens ScreenCollection

	tApplication  *tview.Application
	tBox          *tview.Box
	currentScreen *Screen

	isRunning bool

	screenBox tview.Box
}

func NewConsole() Console {
	c := Console{}

	c.screens = NewScreenCollection()
	c.tApplication = tview.NewApplication()
	c.currentScreen = nil

	c.tBox = tview.NewBox()

	return c
}

func (c *Console) RegisterScreen(name string, screen Screen) {
	c.screens.AddScreen(name, screen)
}

func (c *Console) LoadScreen(name string) {
	c.currentScreen = c.screens.GetScreen(name)
	currentScreen := *c.currentScreen

	currentScreen.Build()

	c.screenBox = *currentScreen.GetBox()

	c.tBox.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		c.screenBox.Draw(screen)

		return x, y, width, height
	})

	if c.isRunning {
		c.tApplication.Draw()
		c.tApplication.SetFocus(currentScreen.GetBox())
	}
}

func (c *Console) Run() {
	log.Println("Running app")

	err := c.tApplication.SetRoot(c.tBox, true).Run()
	c.isRunning = true

	currentScreen := c.getCurrentScreen()

	if currentScreen != nil {
		c.tApplication.SetFocus(currentScreen.GetBox())
	}

	if err != nil {
		log.Panicln("ERROR", err)
	}
}

func (c *Console) getCurrentScreen() Screen {
	if c.currentScreen != nil {
		return *c.currentScreen
	}

	return nil
}

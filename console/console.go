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

	c.tBox = tview.NewBox() /*
		c.tBox.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
			cc.getCurrentScreen().GetBox().Draw(screen)

			return x, y, width, height
		})
	*/
	return c
}

func (c *Console) RegisterScreen(name string, screen Screen) {
	c.screens.AddScreen(name, screen)
}

func (c *Console) LoadScreen(name string) {
	c.currentScreen = c.screens.GetScreen(name)
	currentScreen := *c.currentScreen

	currentScreen.Build()

	c.tBox.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		c.getCurrentScreen().GetBox().Draw(screen)

		return x, y, width, height
	})

	if c.isRunning {
		c.tApplication.Draw()
		c.tApplication.SetFocus(currentScreen.GetDefaultFocus())
	}
}

func (c *Console) Run() {
	log.Println("Running app")

	currentScreen := c.getCurrentScreen()

	if currentScreen != nil {
		c.tApplication.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			c.getCurrentScreen().GetBox().GetInputCapture()(event)
			return event
		})
	}

	err := c.tApplication.SetRoot(c.tBox, true).Run()

	c.tApplication.SetFocus(currentScreen.GetDefaultFocus())

	log.Println("screen", currentScreen)
	c.isRunning = true

	log.Println("after focus")
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

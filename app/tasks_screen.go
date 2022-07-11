package app

import (
	"log"

	"github.com/filipgorny/org-tool/organizer"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TasksScreen struct {
	orgService   organizer.OrganizerService
	box          *tview.Box
	tasks        []organizer.Task
	selectedTask *organizer.Task
	limit        int
	page         int
	focus        bool
	drawTasks    func()
	tasksTable   *tview.Grid
	focusTarget  tview.Primitive
}

func NewTaskScreen() TasksScreen {
	ts := TasksScreen{}

	ts.limit = 20
	ts.page = 0
	ts.drawTasks = nil

	ts.orgService = organizer.InitializeOrgService()

	ts.box = tview.NewBox()
	ts.tasks = make([]organizer.Task, 0)

	return ts
}

func (ts *TasksScreen) loadTasks() {
	ts.tasks = ts.orgService.LoadTasks(ts.limit, ts.page*ts.limit)

	ts.selectedTask = nil
	for i := 0; i < ts.limit; i++ {
		if i >= len(ts.tasks) {
			break
		}

		if ts.selectedTask == nil {
			ts.selectedTask = &ts.tasks[i]
		}
	}
}

func (ts *TasksScreen) Build() {

	ts.loadTasks()

	ts.tasksTable = tview.NewGrid()
	ts.tasksTable.SetColumns(3, 10, 0)
	ts.tasksTable.SetRows(0)
	ts.tasksTable.SetBorders(false)

	box := ts.tasksTable.Clear().Box

	ts.focusTarget = box

	doneCheck := ""

	taskColor := tcell.ColorWhite

	newPrimitive := func(color tcell.Color, text string) *tview.TextView {
		return tview.NewTextView().
			SetTextColor(color).
			SetText(text)
	}

	lineY := 0

	ts.drawTasks = func() {
		for _, task := range ts.tasks {
			if task.Done {
				doneCheck = "\u25c9"
			} else {
				doneCheck = "\u25ef"
				taskColor = tcell.ColorOrange
			}

			label := task.Label

			if len(task.Label) == 0 {
				label = "<no label>"
			}

			focus := task == *ts.getSelectedTask()

			labelColor := tcell.ColorWhite
			numberColor := tcell.ColorGray

			if focus {
				labelColor = tcell.ColorYellow
				numberColor = tcell.ColorWhite
			}

			cols := make([]*tview.TextView, 3)
			cols[0] = newPrimitive(taskColor, doneCheck)
			cols[1] = newPrimitive(numberColor, task.Number)
			cols[2] = newPrimitive(labelColor, label)

			ts.tasksTable.AddItem(cols[0], lineY, 0, 1, 1, 1, 1, false)
			ts.tasksTable.AddItem(cols[1], lineY, 1, 1, 1, 1, 1, false)
			ts.tasksTable.AddItem(cols[2], lineY, 2, 1, 1, 1, 1, false)

			lineY++
		}
	}

	ts.drawTasks()
	redrawTable := func(screen tcell.Screen) {}
	var currentTviewScreen tcell.Screen
	currentTviewScreen = nil

	ts.box.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		redrawTable = func(screen tcell.Screen) {
			screenWidth, _ := screen.Size()
			ts.tasksTable.SetRect(2, 1, screenWidth-2, height-2)
			ts.tasksTable.SetSize(len(ts.tasks), 3, 1, 0)
			ts.tasksTable.SetColumns(5, 15, 0)

			ts.tasksTable.Draw(screen)
		}

		redrawTable(screen)

		currentTviewScreen = screen

		return x, y, width, height
	})

	box.SetFocusFunc(func() {
		log.Println("focus")
		ts.focus = true

		ts.loadTasks()
		ts.drawTasks()

		ts.focusTarget = ts.tasksTable.Clear().Box

		if currentTviewScreen != nil {
			ts.loadTasks()
			ts.tasksTable.Clear()
			ts.drawTasks()
		}
	})

	ts.box.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' {
			ts.keyDown()
		}

		if event.Rune() == 'k' {
			ts.keyUp()
		}

		redrawTable(currentTviewScreen)

		return event
	})
}

func (ts *TasksScreen) GetBox() *tview.Box {
	return ts.box
}

func (ts *TasksScreen) keyUp() {
	var prev *organizer.Task
	for _, task := range ts.tasks {
		if prev != nil {
			ts.selectTask(prev)
		}

		prev = &task
	}

	ts.drawTasks()
}

func (ts *TasksScreen) keyDown() {
	selectNext := false

	for _, task := range ts.tasks {
		if task == *ts.selectedTask {
			selectNext = true
		}

		if selectNext {
			ts.selectTask(&task)
			selectNext = false
			break
		}
	}

	ts.drawTasks()
}

func (ts *TasksScreen) selectTask(task *organizer.Task) {
	ts.selectedTask = task
}

func (ts *TasksScreen) getSelectedTask() *organizer.Task {
	return ts.selectedTask
}

func (ts *TasksScreen) GetDefaultFocus() tview.Primitive {
	if ts.focusTarget != nil {
		return ts.focusTarget
	}

	return ts.box
}

func (ts *TasksScreen) isFocused() bool {
	return ts.focus
}

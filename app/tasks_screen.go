package app

import (
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
}

func NewTaskScreen() TasksScreen {
	ts := TasksScreen{}

	ts.limit = 20
	ts.page = 0

	ts.orgService = organizer.InitializeOrgService()

	ts.box = tview.NewBox()
	ts.tasks = make([]organizer.Task, 0)

	return ts
}

func (ts *TasksScreen) loadTasks() {
	ts.tasks = ts.orgService.LoadTasks(ts.limit, ts.page*ts.limit)

	selectedTask := false
	for i := 0; i < ts.limit; i++ {
		if i >= len(ts.tasks) {
			break
		}

		if selectedTask == false {
			ts.selectTask(&ts.tasks[i])
			selectedTask = true
		}
	}

	if len(ts.tasks) > 0 {
		ts.selectedTask = &ts.tasks[0]
	}
}

func (ts *TasksScreen) Build() {
	ts.loadTasks()

	tasksTable := tview.NewGrid()
	tasksTable.SetColumns(0, 0, 0)
	tasksTable.SetRows(1)
	tasksTable.SetBorders(true)

	doneCheck := ""

	taskColor := tcell.ColorWhite

	newPrimitive := func(color tcell.Color, text string) tview.Primitive {
		return tview.NewTextView().
			SetTextColor(color).
			SetText(text)
	}
	for index, task := range ts.tasks {
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

		tasksTable.AddItem(newPrimitive(taskColor, doneCheck), index, 0, 1, 1, 1, 1, false)
		tasksTable.AddItem(newPrimitive(tcell.ColorOrange, task.Number), index, 1, 1, 0, 10, 1, false)
		tasksTable.AddItem(newPrimitive(tcell.ColorWhite, label), index, 2, 1, 0, 10, 10, false)
	}

	ts.box.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		screenWidth, screenHeight := screen.Size()
		tasksTable.SetRect(2, 1, screenWidth-2, screenHeight)
		tasksTable.Draw(screen)

		return x, y, width, height
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
}

func (ts *TasksScreen) keyDown() {
	selectNext := false

	for _, task := range ts.tasks {
		if task == *ts.selectedTask {
			selectNext = true
		}

		if selectNext {
			ts.selectTask(&task)
			break
		}
	}
}

func (ts *TasksScreen) selectTask(task *organizer.Task) {
	ts.selectedTask = task
}

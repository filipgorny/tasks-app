package app

import (
	"fmt"

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
	tasksTable := tview.NewGrid()
	tasksTable.SetRows(len(ts.tasks))
	tasksTable.SetColumns(3)

	doneCheck := ""
	labelColor := ""

	for index, task := range ts.tasks {
		taskBox := tview.NewBox()

		if task == *ts.selectedTask {
			taskBox.SetBackgroundColor(tcell.ColorWhite)
		}
		if task.Done {
			labelColor = "gray"
			doneCheck = "\u25c9"
		} else {
			labelColor = "white"
			doneCheck = "\u25ef"
		}

		label := task.Label

		if len(task.Label) == 0 {
			label = "<no label>"
		}

		doneCol := fmt.Sprintf("[white]%s", doneCheck)
		numberCol := fmt.Sprintf("[gray]%s", task.Number)
		labelCol := fmt.Sprintf("[%s]%s", labelColor, label)

		doneCell := tview.NewTextView().SetText(doneCol)
		numberCell := tview.NewTextView().SetText(numberCol)
		labelCell := tview.NewTextView().SetText(labelCol)

		tasksTable.AddItem(doneCell, index, 0, 0, 0, 0, 0, false)
		tasksTable.AddItem(numberCell, index, 1, 0, 0, 0, 0, false)
		tasksTable.AddItem(labelCell, index, 2, 0, 0, 0, 0, false)
	}

	ts.box.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		tasksTable.DrawForSubclass(screen, ts.box)

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

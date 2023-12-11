package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"time"
)

const (
	Completed  = "🗹"
	Open       = "☐"
)

type Tag struct {
	name string
}

type Task struct {
	description string
	tag         Tag
	dateCreated time.Time
	status      int8 // 0: Open, 1: Completed
}

func TaskFactory(description string, tag string) *Task {
	return &Task{description: description, dateCreated: time.Now(), status: 0, tag: Tag{name: tag}}
}

func (t Task) String() string {
	var statusIcon string

	switch t.status {
	case 1:
		statusIcon = Completed
	default:
		statusIcon = Open
	}
	return fmt.Sprintf("| %v | %v | %v", statusIcon, t.description, t.dateCreated.UTC().Format("Mon Jan 2"))
}

func (t *Task) Update(status int8) {
	t.status = status
}

func PrintTasks(tasks map[int]*Task) {
	tasksOrder := make([]int, 0)
	for k := range tasks {
		tasksOrder = append(tasksOrder, k)
	}
	sort.SliceStable(tasksOrder, func(i, j int) bool {return tasksOrder[i] < tasksOrder[j]})
	for _, v := range tasksOrder {
			fmt.Println("|", v, tasks[v], "|")
	}
}

func main() {
	tasks := map[int]*Task{
		1: TaskFactory("This is a test", "All"),
		2: TaskFactory("This is another test", "All"),
		3: TaskFactory("This is a different test", "All"),
	}
	var newTaskFlag bool
	var taskDescFlag string
	var taskTagFlag string
	var taskNumberFlag bool

	flag.BoolVar(&newTaskFlag, "t", false, "")
	flag.StringVar(&taskDescFlag, "d", "", "Usage:")
	flag.BoolVar(&taskNumberFlag, "update", false, "Usage:")
	flag.StringVar(&taskTagFlag, "tag", "All", "Usage:")
	flag.BoolFunc("l", "Usage", func(string) error {
		PrintTasks(tasks)
		return nil
	})
	flag.Parse()

	if newTaskFlag {
		if taskDescFlag == "" {
			fmt.Println("Cannot have empty description")
			os.Exit(1)
		}
		task := TaskFactory(taskDescFlag, taskTagFlag)

		tasks[len(tasks) - 1] = task
		PrintTasks(tasks)
	}

	if taskNumberFlag {
		var tasksIndex []int 
		for _, v := range flag.Args() {
			idx, err := strconv.Atoi(v)

			if err != nil {
				os.Exit(1)
			}
			if idx < 1 {
				os.Exit(1)
			}
			tasksIndex = append(tasksIndex, idx)
		}

		if len(tasksIndex) <= 0 {
			os.Exit(1)
		}

		if len(tasksIndex) > len(tasks) {
			os.Exit(1)
		}

		tasksIndex = slices.Compact(tasksIndex)
		for _, v := range tasksIndex {
			task, ok := tasks[v]
			if !ok {
				os.Exit(1)
			}
			task.Update(1)
		}
		PrintTasks(tasks)
	}
}

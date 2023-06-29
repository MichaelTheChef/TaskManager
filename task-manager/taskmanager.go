package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type TaskManager struct {
	Tasks []*Task
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		Tasks: make([]*Task, 0),
	}
}

func (tm *TaskManager) AddTask(name string, memory int) {
	task := &Task{
		Name:   name,
		Memory: memory,
	}
	tm.Tasks = append(tm.Tasks, task)
}

func (tm *TaskManager) RemoveTask(name string) {
	for i, task := range tm.Tasks {
		if task.Name == name {
			tm.Tasks = append(tm.Tasks[:i], tm.Tasks[i+1:]...)
			break
		}
	}
}

func (tm *TaskManager) GetTasks() []*Task {
	return tm.Tasks
}

func (tm *TaskManager) Run() {
	for {
		switch runtime.GOOS {
		case "windows":
			tm.getWindowsTasks()
		case "darwin":
			tm.getDarwinTasks()
		case "linux":
			tm.getLinuxTasks()
		}

		time.Sleep(1 * time.Second) 
	}
}

func (tm *TaskManager) getWindowsTasks() {
	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error retrieving task list:", err)
		return
	}

	tasks := strings.Split(string(output), "\n")
	for _, task := range tasks {
		fields := strings.Split(task, `","`)
		if len(fields) >= 5 {
			name := strings.Trim(fields[0], `"`)
			memory := strings.Trim(fields[4], `"`)
			tm.AddTask(name, parseMemory(memory))
		}
	}
}

func (tm *TaskManager) getDarwinTasks() {
	cmd := exec.Command("ps", "-A", "-o", "pid,rss,comm")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error retrieving task list:", err)
		return
	}

	tasks := strings.Split(string(output), "\n")
	for _, task := range tasks {
		fields := strings.Fields(task)
		if len(fields) >= 3 {
			name := fields[2]
			memory := fields[1]
			tm.AddTask(name, parseMemory(memory))
		}
	}
}

func (tm *TaskManager) getLinuxTasks() {
	cmd := exec.Command("ps", "-e", "-o", "pid,rss,comm")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error retrieving task list:", err)
		return
	}

	tasks := strings.Split(string(output), "\n")
	for _, task := range tasks {
		fields := strings.Fields(task)
		if len(fields) >= 3 {
			name := fields[2]
			memory := fields[1]
			tm.AddTask(name, parseMemory(memory))
		}
	}
}

func parseMemory(memory string) int {
	var size int
	unit := strings.ToUpper(memory[len(memory)-2:])
	memory = memory[:len(memory)-2]
	fmt.Sscanf(memory, "%d", &size)

	switch unit {
	case "KB":
		size /= 1024
	case "GB":
		size *= 1024
	}

	return size
}

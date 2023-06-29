package main

import (
	"fmt"
	"time"
)

func main() {
	taskManager := NewTaskManager()

	go taskManager.Run()

	for {
		tasks := taskManager.GetTasks()

		fmt.Println("Current Tasks:")
		for _, task := range tasks {
			fmt.Printf("Name: %s, Memory: %d MB\n", task.Name, task.Memory)
		}

		time.Sleep(5 * time.Second)
	}
}

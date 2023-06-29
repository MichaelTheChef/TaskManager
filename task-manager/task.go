package main

import (
	"fmt"
)

type Task struct {
	Name   string
	Memory int
}

func (t *Task) String() string {
	return fmt.Sprintf("Name: %s, Memory: %d MB", t.Name, t.Memory)
}

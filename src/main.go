package main

import (
	"f"
	"fmt"
)

func main() {
	flow := f.New("sample", true)
	s1 := f.NewState("state-1")
	s1.AddTasks(
		f.Task{
			Name: "Task-1: Get user profilefrom Github",
			Run: func() (*f.Result, error) {
				fmt.Println("This is task 1")
				r := &f.Result{}
				return r, nil

			},
		},
		f.Task{
			Name: "Task-2: Get user profilefrom Github",
			Run: func() (*f.Result, error) {
				fmt.Println("This is task 2")
				r := &f.Result{}
				return r, nil

			},
		},
	)
	s2 := f.NewState("state-2")
	s3 := f.NewState("state-3")
	s4 := f.NewState("state-4")
	flow.AddStates(s1, s2, s3, s4)
	flow.Execute()
}

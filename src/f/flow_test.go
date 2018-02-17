package f

import (
	"fmt"
	"testing"
)

func TestAddFlow(t *testing.T) {
	// flow := New("Test", true)
	s1 := NewState("state-1")
	t.Log("should return tasks")
	{
		s1.AddTasks(
			Task{
				Name: "Task-1",
				Execute: func() (*Output, error) {
					fmt.Println("This is task 1")
					r := &Result{}
					return r, nil

				},
			},
			Task{
				Name: "Task-2",
				Run: func() (*Result, error) {
					fmt.Println("This is task 2")
					r := &Result{}
					return r, nil

				},
			},
		)

		tasks := s1.GetTasks()
		if len(tasks) != 2 {
			t.Fatal("Expected 2 tasks, got ", len(tasks))
		}
		if tasks[0].Name != "Task-1" && tasks[1].Name != "Task-2" {
			t.Fatal(
				"Expected asks name to be Task-1 & Task-2, but got ",
				tasks,
			)
		}

		task3 := Task{
			Name: "Task-3",
			Run: func() (*Result, error) {
				fmt.Println("This is task 3")
				r := &Result{}
				return r, nil

			},
		}

		s1.AddTasks(task3)
		tasks = s1.GetTasks()
		if len(tasks) != 3 {
			t.Fatal("Expected 3 tasks, got ", len(tasks))
		}
		if tasks[2].Name != "Task-3" {
			t.Fatal(
				"Expected task 3's name to be Task-3 but got ",
				tasks,
			)
		}
	}

}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
}

type TodoList struct {
	todos    []Todo
	filename string
}

func NewTodoList(filename string) (*TodoList, error) {
	tl := &TodoList{filename: filename}
	err := tl.load()
	return tl, err
}

func (tl *TodoList) Add(title string) {
	todo := Todo{
		ID:        len(tl.todos) + 1,
		Title:     title,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	tl.todos = append(tl.todos, todo)
	tl.save()
	fmt.Printf("✓ Added: %s\n", title)
}

func (tl *TodoList) List() {
	if len(tl.todos) == 0 {
		fmt.Println("No todos yet.")
		return
	}
	for _, todo := range tl.todos {
		status := "[ ]"
		if todo.Completed {
			status = "[✓]"
		}
		fmt.Printf("%d. %s %s (%s)\n", todo.ID, status, todo.Title, todo.CreatedAt)
	}
}

func (tl *TodoList) Reset() {
	tl.todos = []Todo{}
	tl.save()
	fmt.Println("✓ All todos reset.")
}

func (tl *TodoList) Done(id int) error {
	for i := range tl.todos {
		if tl.todos[i].ID == id {
			tl.todos[i].Completed = true
			tl.save()
			fmt.Printf("✓ Marked done: %s\n", tl.todos[i].Title)
			return nil
		}
	}
	return fmt.Errorf("no todo with ID %d", id)
}

func (tl *TodoList) Delete(id int) error {
	for i := range tl.todos {
		if tl.todos[i].ID == id {
			fmt.Printf("✗ Deleted: %s\n", tl.todos[i].Title)
			tl.todos = append(tl.todos[:i], tl.todos[i+1:]...)
			tl.save()
			return nil
		}
	}
	return fmt.Errorf("no todo with ID %d", id)
}

func (tl *TodoList) save() error {
	data, err := json.MarshalIndent(tl.todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tl.filename, data, 0644)
}

func (tl *TodoList) load() error {
	data, err := os.ReadFile(tl.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, &tl.todos)
}

func main() {
	tl, err := NewTodoList("todos.json")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  todo add <title>")
		fmt.Println("  todo list")
		fmt.Println("  todo done <id>")
		fmt.Println("  todo delete <id>")
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo add <title>")
			os.Exit(1)
		}
		tl.Add(os.Args[2])

	case "list":
		tl.List()

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo done <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID — must be a number")
			os.Exit(1)
		}
		if err := tl.Done(id); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: todo delete <id>")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID — must be a number")
			os.Exit(1)
		}
		if err := tl.Delete(id); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	case "reset":
		tl.Reset()

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
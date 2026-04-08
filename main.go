package main

import (
	"fmt"
	"os"
)

type Todo struct {
	ID int
	Text string
	Done bool
}

func add_todo(todos []Todo, todo_text string) []Todo {
	thisid := len(todos) + 1

	todo := Todo{
		ID: thisid,
		Text: todo_text,
		Done: false,
	}

	todos = append(todos, todo)

	return todos
}

func delete_todo(todos []Todo, id int) []Todo {

	result := []Todo{}

	for _, todo := range todos {
		if todo.ID != id {
			result = append(result, todo)
		}
	}

	return result
}

func complete_task(todos []Todo, id int) []Todo {

	for i := range todos {

		if todos[i].ID == id {
			todos[i].Done = true
			break
		}
	}

	return todos
}

func print_todos(todos []Todo) {
	
	if len(todos) == 0 {
		fmt.Println("Not found todo")
		return
	}

	for _, todo := range todos {
		status := " "

		if todo.Done == true {
			status = "✓"
		}

		fmt.Printf("[%s] %d: %s\n", status, todo.ID, todo.Text)
	}
	// %s 文字列
	// %d 整数
	// %v 適当な時
	// %T 型
	// %f 小数
	// %t Bool
}

func arg_commands(todos []Todo, Args []string) {
	if len(Args) <= 2 {
		fmt.Println("Type: todo (list | add <task> | delete <id> | complete <id>)")
		return
	}
	if Args[1] != "todo" {
		fmt.Println("Type: todo (list | add <task> | delete <id> | complete <id>)")
		return
	}
	commands := Args[2]

	switch commands {
	case "add":
		if len(Args) <= 3 {
			fmt.Println("Hint: add <task>")
			return
		}

		task := Args[3]
		fmt.Printf("add task \"%s\"", task)

		return

	case "delete":
		if len(Args) <= 3 {
			fmt.Println("Hint: delete <id>")
			return
		}

		id := Args[3]
		fmt.Printf("delete task id: %s", id)

		return

	case "complete":
		if len(Args) <= 3 {
			fmt.Println("Hint: complete <id>")
			return
		}

		id := Args[3]
		fmt.Printf("complete task id: %s", id)

		return

	case "list":
		fmt.Println("list")
		
		for _, todo := range todos {
			completed := " "
			if todo.Done == true {
				completed = "✓"
			}
			fmt.Printf("[%s]%d: %s\n", completed, todo.ID, todo.Text)
		}

		return

	default:
		fmt.Println("this command not found\nType: todo (list | add <task> | delete <id> | complete <id>)")

		return
	}
}

func main() {

	todos := []Todo{}

	newtodo := "コード書く"

	// append は新しい領域で変更を加える場合があるため
	// 更新後のtodosを返し適応する
	todos = add_todo(todos, newtodo)
	newtodo = "ピスタチオ食べる"
	todos = add_todo(todos, newtodo)
	newtodo = "コミットする"
	todos = add_todo(todos, newtodo)

	complete_task_id := 1
	todos = complete_task(todos, complete_task_id)
	complete_task_id = 3
	todos = complete_task(todos, complete_task_id)

	arg_commands(todos, os.Args)

	//print_todos(todos)
}
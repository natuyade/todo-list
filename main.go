package main

import "fmt"

type Todo struct {
	ID int
	Text string
	Done bool
}

func addtodo(todos []Todo, todo_text string) []Todo {
	thisid := len(todos) + 1

	todo := Todo{
		ID: thisid,
		Text: todo_text,
		Done: false,
	}

	todos = append(todos, todo)

	return todos
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

func main() {
	todos := []Todo{}

	newtodo := "コード書く"

	// append は新しい領域で変更を加える場合があるため
	// 更新後のtodosを返し適応する
	todos = addtodo(todos, newtodo)
	newtodo = "コミットする"
	todos = addtodo(todos, newtodo)

	complete_task_id := 1
	todos = complete_task(todos, complete_task_id)
	complete_task_id = 2
	todos = complete_task(todos, complete_task_id)

	print_todos(todos)
}
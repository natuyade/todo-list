package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Todo struct {
	ID int `json:"id"`
	Text string `json:"text"`
	Done bool `json:"done"`
}

func add_todo(todos []Todo, task string) {
	todo := Todo{
		ID: len(todos) + 1,
		Text: task,
		Done: false,
	}

	todos = append(todos, todo)

	marshaled_todos, ms_err := json.MarshalIndent(todos, "", "    ")
	if ms_err != nil {
		fmt.Println(ms_err)
		return
	}

	os.WriteFile("todo_list.json", marshaled_todos, 0666)
}

func delete_todo(todos []Todo, id int) {
	
	result := []Todo{}
	set_id := 1

	for i, todo := range todos {

		fmt.Println(set_id)
		if todos[i].ID != id {

			todo.ID = set_id
			set_id += 1

			result = append(result, todo)
		}
	}

	marshaled_todos, ms_err := json.MarshalIndent(result, "", "    ")
	if ms_err != nil {
		fmt.Println(ms_err)
		return
	}

	os.WriteFile("todo_list.json", marshaled_todos, 0666)
}

func complete_task(todos []Todo, id int) {

	for i := range todos {

		if todos[i].ID == id {
			todos[i].Done = true
			break
		}
	}

	marshaled_todos, ms_err := json.MarshalIndent(todos, "", "    ")
	if ms_err != nil {
		fmt.Println(ms_err)
		return
	}

	os.WriteFile("todo_list.json", marshaled_todos, 0666)
}

func print_todos(todos []Todo) {

	for _, todo := range todos {

		completed := " "
		if todo.Done == true {
			completed = "✓"
		}

		fmt.Printf("[%s]%d: %s\n", completed, todo.ID, todo.Text)
	}
	// %s 文字列
	// %d 整数
	// %v 適当な時
	// %T 型
	// %f 小数
	// %t Bool
}

func arg_commands(Args []string) {
	if len(Args) <= 2 {
		fmt.Println("Type: todo (list | add <task> | delete <id> | complete <id>)")
		return
	}
	if Args[1] != "todo" {
		fmt.Println("Type: todo (list | add <task> | delete <id> | complete <id>)")
		return
	}
	commands := Args[2]

	todos_bytes, read_err := os.ReadFile("todo_list.json")
	if read_err != nil {
		fmt.Println(read_err)
		return
	}

	var todos []Todo
	// bytes -> json 変換
	unms_err := json.Unmarshal(todos_bytes, &todos)
	if unms_err != nil {
		fmt.Println(unms_err)
	}

	switch commands {
	case "add":
		if len(Args) <= 3 {
			fmt.Println("Hint: add <task>")
			return
		}

		task := Args[3]

		fmt.Printf("add task \"%s\"\n", task)

		add_todo(todos, task)

		return

	case "delete":
		if len(Args) <= 3 {
			fmt.Println("Hint: delete <id>")
			return
		}

		id_string := Args[3]

		id, conv_err := strconv.Atoi(id_string)
		if conv_err != nil {
			fmt.Println(conv_err)
			return
		}

		delete_todo(todos, id)

		fmt.Printf("delete task id: %d\n", id)

	case "complete":
		if len(Args) <= 3 {
			fmt.Println("Hint: complete <id>")
			return
		}

		id_string := Args[3]

		id, conv_err := strconv.Atoi(id_string)
		if conv_err != nil {
			fmt.Println(conv_err)
			return
		}

		complete_task(todos, id)

		fmt.Printf("complete task id: %d\n", id)

	case "list":
		fmt.Println("list")

		print_todos(todos)

		return

	default:
		fmt.Println("this command not found\nType: todo (list | add <task> | delete <id> | complete <id>)")
	}
}

func main() {

	arg_commands(os.Args)

	fmt.Println("Finished")
}
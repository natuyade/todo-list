package main

// %s 文字列
// %d 整数
// %v 適当な時
// %T 型
// %f 小数
// %t Bool
import (
	"encoding/json"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	textinput "charm.land/bubbles/v2/textinput"
	lipgloss "charm.land/lipgloss/v2"
)

type Todo struct {
	Text string `json:"text"`
	Done bool `json:"done"`
}

type Layer struct {
	SizeX int
	SizeY int
	RenderS string
	PosX int
	PosY int
}

type model struct{
	inputText textinput.Model
	todos []Todo
	cursor int
	// mapはキーの型を自由に変えられる(stringならselected["yamada"]などで指定できる)
	// [int]がキーの型, struct{}は値の型
	selected map[int]struct{}
	layers []Layer
	enableLayer int
}

type todosLoadedMessage struct {
	todos []Todo
	Err error
}

func loadTodos() tea.Msg {
	todos_bytes, read_err := os.ReadFile("todo_list.json")
	if read_err != nil {
		var first_todo = []Todo{}
		// aにappend(aを代入している理由...
		// 普通にappend(a, ..)だけすればいいのではと思っていましたが
		// append側で代入されたaは元のaと別ものとして作られるため
		// 元のaには代入されないのでa = append(a, ..)とする必要がある
		first_todo = append(first_todo, Todo{
			Text: "コミットする",
			Done: false,
		})
		return todosLoadedMessage{
			todos: first_todo,
			Err: nil,
		}
	}

	var todos []Todo
	// bytes -> json 変換
	unms_err := json.Unmarshal(todos_bytes, &todos)
	if unms_err != nil {
		return todosLoadedMessage{
			todos: nil,
			Err: unms_err,
		}
	}

	return todosLoadedMessage{
		todos: todos,
		Err: nil,
	}
}

func saveTodosToFile(m model) {
	for i := range m.todos {
		m.todos[i].Done = false
		_, ok := m.selected[i]
		if ok {
			m.todos[i].Done = true
		}
	}

	todos_json, Err := json.MarshalIndent(m.todos, "", "    ")
	if Err != nil {
		fmt.Println(Err)
		return
	}
	os.WriteFile("todo_list.json", todos_json, 0666)
}

// io init
func (m model) Init() tea.Cmd {

	// Sequenceで順に処理
	return tea.Sequence(
		loadTodos,
	)
}

// model init
func initialModel() model {
	// textinputの設定
	// text入力の場
	ti := textinput.New()
	// place holderは入力値がないときの仮テキスト
	// 十分なwidthを設定しないとplaceholderが見えなくなる
	ti.Placeholder = "enter here"
	// 大文字を考慮したwidthにすべき
	ti.SetWidth(32)
	ti.CharLimit = 16


	return model{
		enableLayer: 0,
		// modelに持たせる
		inputText: ti,

		todos: []Todo{},
		// 空のmapはmake()で作成
		selected: make(map[int]struct{}),
	}
}

// 操作があれば更新される
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case todosLoadedMessage:
		if msg.Err != nil {
			return m, tea.Quit
		}
		m.todos = msg.todos

		for i := range m.todos {
			if m.todos[i].Done == true {
				m.selected[i] = struct{}{}
			}
		}
		
		return m, nil

	// 更新時押されていたkey
	case tea.KeyPressMsg:

		switch m.enableLayer {
		case 0:
			switch msg.String() {
			// それがこのいずれかなら...
			case "q", "ctrl+c":
				// tea.Quitコマンドで終了
				return m, tea.Quit
			
			case "up", "w":
				if m.cursor > 0 {
					m.cursor--
				}

			case "down", "s":
				if m.cursor < len(m.todos) -1 {
					m.cursor++
				}
			case "enter":
				// selected[key]に値が存在しているかを二つ目のvalueでみれる
				// 存在していればtrue, いなければfalseが返る.

				// boolで切り替えればいいのではとか思ったけど,
				// 中身を取り出してboolを得るより外から中身の存在だけでboolを得れるため
				// この方法が使われてるのだと思う. 多分だけど...

				// 中身のvalueはどうでもいいので, 存在だけを確かめる
				_, ok := m.selected[m.cursor]
				if ok {
					// 中身があればdelete(mapType, key)の値を削除
					delete(m.selected, m.cursor)
				} else {
					// なければ, 存在boolをtrueにするために空の構造体を入れる
					m.selected[m.cursor] = struct{}{}
				}
			case "t":
				saveTodosToFile(m)
			case "f":
				m.enableLayer = 1
				m.inputText.Focus()
				return m, nil
			}
		case 1:
			switch msg.String() {
			case "enter":
				task := m.inputText.Value()
				if task != "" {
					m.todos = append(m.todos, Todo{
						Text: task,
						Done: false, 
					})
				}
				m.inputText.SetValue("")
				// blurをかけてfocusを解除
				m.inputText.Blur()
				m.enableLayer = 0
				return m, nil
			case "esc":
				m.enableLayer = 0
				m.inputText.Blur()
				return m, nil
			}
		}

			var cmd tea.Cmd
			m.inputText, cmd = m.inputText.Update(msg)
			return m, cmd
	}
	return m, nil
}

func (m model) View() tea.View {

	//windowSize := tea.RequestWindowSize()
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		Margin(1).
		Padding(1)
	
	lineS := "todo list\n\n"
	for i, todo := range m.todos {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		selected := " "
		_, ok := m.selected[i]
		if ok {
			selected = "✓"
		}
		// Sprintはto_stringへのformat
		lineS += fmt.Sprintf("%s[%s] %s\n", cursor, selected, todo.Text)
	}

	addtS := "enter new todo\n\n"
	addtS += m.inputText.View() + "\n"
	
	helpS := "help\n----------------\n"
	switch m.enableLayer{
	case 0:
		helpS += "|t|Save |f|AddTodo |q|Quit"
	case 1:
		helpS += "|Enter|AddTodo |esc|BackToList"
	}

	// sizex,sizey,posx,posy
	lineLayer := []int{64, 32, 0, 0}
	addTodoLayer := []int{64, 32, 65, 0}
	helpLayer := []int{64, 8, 0, 33}

	layerString := []string{}
	layerString = append(layerString, lineS, addtS, helpS)

	layerSize := [][]int{}
	layerSize = append(layerSize, lineLayer, addTodoLayer, helpLayer)

	layerLen := len(layerString)
	if layerLen < len(layerSize) {
		layerLen = len(layerSize)
	} 

	layer := []Layer{}
	for i := range layerLen {
		layer = append(layer, Layer{
			SizeX: layerSize[i][0],
			SizeY: layerSize[i][1],
			RenderS: layerString[i],
			PosX: layerSize[i][2],
			PosY: layerSize[i][3],
		})
	}
	m.layers = layer
	
	var newLayers []*lipgloss.Layer
	for _, l := range m.layers {
		listLayer := boxStyle.Width(l.SizeX).Height(l.SizeY).
			BorderForeground(lipgloss.Magenta).Render(l.RenderS)
		layer := lipgloss.NewLayer(listLayer).
			X(l.PosX).Y(l.PosY)
		newLayers = append(newLayers, layer)
	}
	
//lipgloss.Color("#e464f0")

	// 可変長引数が必要なものにスライスを渡す際は
	// '...'を付けないと型不一致でエラーを吐く可能性がある
	output := lipgloss.NewCompositor(newLayers...).Render()

	view := tea.NewView(output)
	view.AltScreen = true

    return view
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
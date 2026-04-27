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
	"time"

	textinput "charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

const Description = "私natuyadeがGo言語に慣れるため\n色々適当に触ってあそんでます.\n課題まみれでよく頭を抱えてます\n；；"

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
	logs []string
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

func deleteTodo(m model) model {
	result := []Todo{}
	for i, todo := range m.todos {
		_, ok := m.selected[i]
		if ok {
			delete(m.selected, i)
			continue
		}
		result = append(result, todo)
	}
	m.todos = result
	m.cursor = 0

	return m
}

func pushToLog(m model, log string) model {
	now := time.Now()
	nowString := now.Format("15:04:05")

	logS := fmt.Sprintf("[%s]%s", nowString, log)
	m.logs = append(m.logs, logS)

	return m
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

	// layerが何枚になるかを初期の時点で決めておきたいのでここである程度作成
	lineLayer := Layer{SizeX: 41, SizeY: 32, RenderS: "", PosX: 0, PosY: 0}
	addTodoLayer := Layer{SizeX: 42, SizeY: 8, RenderS: "", PosX: 0, PosY: 0}
	tes1Layer := Layer{SizeX: 42, SizeY: 8, RenderS: "", PosX: 0, PosY: 0}
	tes2Layer := Layer{SizeX: 42, SizeY: 8, RenderS: "", PosX: 0, PosY: 0}
	tes3Layer := Layer{SizeX: 42, SizeY: 8, RenderS: "", PosX: 0, PosY: 0}
	tes4Layer := Layer{SizeX: 42, SizeY: 8, RenderS: "", PosX: 0, PosY: 0}
	tes5Layer := Layer{SizeX: 42, SizeY: 8, RenderS: "", PosX: 0, PosY: 0}
	tes6Layer := Layer{SizeX: 42, SizeY: 8, RenderS: "", PosX: 0, PosY: 0}
	helpLayer := Layer{SizeX: 42, SizeY: 8, RenderS: "", PosX: 0, PosY: 0}

	var layers []Layer
	layers = append(layers, lineLayer, addTodoLayer, tes1Layer, tes2Layer, tes3Layer, tes4Layer, tes5Layer, tes6Layer, helpLayer)

	return model{
		enableLayer: 0,
		// modelに持たせる
		inputText: ti,
		layers: layers,
		todos: []Todo{},
		// 空のmapはmake()で作成
		selected: make(map[int]struct{}),
	}
}

// 操作があれば更新される
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.layers[0].SizeY = msg.Height - 1
		for i := range m.layers {
			if i != 0 {
				if m.layers[i - 1].PosX + m.layers[i - 1].SizeX + m.layers[i].SizeX + 1 > msg.Width {
					m.layers[i].PosX = m.layers[0].SizeX + 1
					m.layers[i].PosY = m.layers[i - 1].PosY + m.layers[i - 1].SizeY + 1
				} else {
					m.layers[i].PosX = m.layers[i - 1].PosX + m.layers[i - 1].SizeX + 1
					m.layers[i].PosY = m.layers[i - 1].PosY
				}
			}
			
		}

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
		// 0=list 1=addtodo +1=help
		switch msg.String() {
			case "tab":
				switch m.enableLayer {
				case 0:
					m.inputText.Focus()
				case 1:
					m.inputText.Blur()
				}
				m.enableLayer = min(m.enableLayer + 1, len(m.layers) - 2)
				return m, nil
			case "esc":
				switch m.enableLayer {
				case 1:
					m.inputText.Blur()
				case 2:
					m.inputText.Focus()
				}
				m.enableLayer = max(m.enableLayer - 1, 0)
				return m, nil
			case "q", "ctrl+c":
				// tea.Quitコマンドで終了
				return m, tea.Quit
			case "enter":
				switch m.enableLayer {
				case 0:
					_, ok := m.selected[m.cursor]
					if ok {
						delete(m.selected, m.cursor)
					} else {
						m.selected[m.cursor] = struct{}{}
					}
				case 1:
					task := m.inputText.Value()
					if task != "" {
						m.todos = append(m.todos, Todo{
							Text: task,
							Done: false, 
						})
					}
					m.inputText.SetValue("")
					m.inputText.Blur()
					m.enableLayer = 0
					return m, nil
				}
			case "up", "w":
				switch m.enableLayer{
				case 0:
					if m.cursor > 0 {
						m.cursor--
					}
				}
			case "down", "s":
				switch m.enableLayer{
				case 0:
					if m.cursor < len(m.todos) -1 {
						m.cursor++
					}
				}
			case "d":
				switch m.enableLayer{
				case 0:
					m = deleteTodo(m)
					m = pushToLog(m, "deleted complete task")
				}
			case "t":
				switch m.enableLayer {
				case 0:
					saveTodosToFile(m)
					m = pushToLog(m, "saved to file")
				}
		}

		var cmd tea.Cmd
		m.inputText, cmd = m.inputText.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() tea.View {

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		Margin(-1).
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
		// Sprintはstringへのformat
		lineS += fmt.Sprintf("%s[%s] %s\n", cursor, selected, todo.Text)
	}

	addtS := "enter new todo\n\n"
	addtS += m.inputText.View() + "\n"

	logS := ""
	limit := 0
	for i := len(m.logs); i != 0; i-- {
		limit += 1
		if limit > 4 {
			break
		}

		if i != len(m.logs) {
			logS += "\n"
		}
		logS += fmt.Sprintf("%s", m.logs[i - 1])
	}
	
	helpS := "help\n----------------\n|tab|Next |esc|Prev |q|Quit "
	switch m.enableLayer{
	case 0:
		helpS += "|t|Save |d|DeleteCompTodo"
	case 1:
		helpS += "|Enter|AddTodo"
	}
	var layerStrings []string
	layerStrings = append(layerStrings, lineS, addtS, logS, "", "", "", "", Description, helpS)

	for i := range m.layers {
		m.layers[i].RenderS = layerStrings[i]
	}
	
	var newLayers []*lipgloss.Layer
	for i, l := range m.layers {
		borderColor := lipgloss.Magenta
		if i == m.enableLayer {
			borderColor = lipgloss.BrightMagenta
		}
		listLayer := boxStyle.Width(l.SizeX).Height(l.SizeY).
			BorderForeground(borderColor).Render(l.RenderS)
		layer := lipgloss.NewLayer(listLayer).
			X(l.PosX).Y(l.PosY)
		newLayers = append(newLayers, layer)
	}

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
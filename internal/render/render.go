package render

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"hang.sh/internal/entity"
	"hang.sh/internal/shared"
	"hang.sh/internal/sprite"
	"hang.sh/internal/ui"
	"hang.sh/internal/world"
)

type TickMsg struct{}

type GameModel struct {
	width, height int
	sized         bool

	screenType ui.ScreenType
	world      *world.World
	player     *entity.Player

	camX, camY int

	nameInput textinput.Model
	headIndex int

	lastInputTime time.Time

	animFrame int
	animTick  int

	buffer [][]*shared.Cell

	chatInput  textinput.Model
	chatActive bool
	chatOffset int
	msgOffset  int
	shiftHeld  bool
}

func NewGameModel(p *entity.Player, w *world.World) GameModel {
	ti := textinput.New()
	ti.Placeholder = "your name"
	ti.Focus()
	ti.CharLimit = 20
	ti.SetWidth(20)

	ci := textinput.New()
	ci.Placeholder = "chat..."
	ci.CharLimit = 200
	ci.SetWidth(40)

	return GameModel{
		screenType: ui.ScreenStart,
		world:      w,
		player:     p,
		nameInput:  ti,
		chatInput:  ci,
	}
}

func (m GameModel) Init() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func setupView(v *tea.View) {
	v.AltScreen = true
	v.MouseMode = tea.MouseModeAllMotion
	v.KeyboardEnhancements.ReportEventTypes = true
}

func (m GameModel) View() tea.View {
	if !m.sized {
		v := tea.NewView("")
		setupView(&v)
		return v
	}

	switch m.screenType {
	case ui.ScreenStart:
		return m.renderTitle()
	case ui.ScreenGame:
		return m.renderGame()
	}

	v := tea.NewView("")
	setupView(&v)
	return v
}

// --- Title screen ---

func (m GameModel) renderTitle() tea.View {
	lines := make([]string, 0, m.height)

	push := func(s string) { lines = append(lines, s) }

	for i := 0; i < m.height/4; i++ {
		push("")
	}

	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00ff88")).Render("HANG")
	push(centerText(title, m.width))
	push("")

	namePrompt := lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa")).Render("what is your name?")
	push(centerText(namePrompt+"  "+m.nameInput.View(), m.width))
	push("")

	headLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa")).Render("choose your head")
	push(centerText(headLabel, m.width))

	prev := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render("◀")
	next := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render("▶")
	current := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color("#333333")).
		Render(" " + ui.HeadOptions[m.headIndex] + " ")
	push(centerText(prev+"  "+current+"  "+next, m.width))
	push("")

	preview := playerPreview(ui.HeadOptions[m.headIndex])
	for _, pl := range strings.Split(preview, "\n") {
		push(centerText(pl, m.width))
	}
	push("")

	play := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff88")).Render("Press Enter to play")
	push(centerText(play, m.width))

	for len(lines) < m.height {
		push("")
	}

	v := tea.NewView(strings.Join(lines, "\n"))
	setupView(&v)
	return v
}

func playerPreview(headChar string) string {
	return " " + headChar + " \n" + "/|\\\n" + "/ \\"
}

// --- Game screen ---

func (m GameModel) renderGame() tea.View {
	m.buildBuffer()
	m.renderTiles()
	m.renderEntities()

	lines := make([]string, m.height)
	for y, row := range m.buffer {
		lines[y] = renderCellRow(row)
	}

	// HUD: player name
	if m.player.Name != "" {
		nameRow := m.height - 1
		nameStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff88")).
			Background(lipgloss.Color("#1a1a2e"))
		hud := nameStyle.Render(" " + m.player.Name + " ")
		hudVW := visualWidth(hud)
		if hudVW < m.width {
			hud += strings.Repeat(" ", m.width-hudVW)
		}
		lines[nameRow] = hud
	}

	// online count
	count := 0
	m.world.ForEachPlayer(func(p *entity.Player) {
		count++
	})
	onlineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	onlineText := onlineStyle.Render(fmt.Sprintf(" %d online", count))
	lines[0] = onlineText + strings.Repeat(" ", m.width-visualWidth(onlineText))

	// chat overlay
	if m.chatActive {
		lines = m.renderChatOverlay(lines)
	} else {
		lines = m.renderChatHistory(lines)
	}

	v := tea.NewView(strings.Join(lines, "\n"))
	setupView(&v)
	return v
}

func (m *GameModel) buildBuffer() {
	m.buffer = make([][]*shared.Cell, m.height)
	for y := range m.buffer {
		m.buffer[y] = make([]*shared.Cell, m.width)
		for x := range m.buffer[y] {
			m.buffer[y][x] = &shared.Cell{Char: " ", FG: "#ffffff", BG: "#000000"}
		}
	}
}

func (m *GameModel) renderTiles() {
	m.camX = m.player.X - m.width/2
	m.camY = m.player.Y - m.height/2

	for sy := 0; sy < m.height; sy++ {
		wy := m.camY + sy
		for sx := 0; sx < m.width; sx++ {
			wx := m.camX + sx
			tile := m.world.TileAt(wx, wy)
			if tile != nil {
				m.buffer[sy][sx] = &shared.Cell{
					Char: tile.Char,
					FG:   tile.FG,
					BG:   tile.BG,
				}
			}
		}
	}
}

func (m *GameModel) renderEntities() {
	// NPCs
	m.world.ForEachNPC(func(npc *entity.NPC) {
		m.renderSprite(npc.X, npc.Y, npc.Sprite, 0)
	})

	// Other players
	m.world.ForEachPlayer(func(p *entity.Player) {
		if p.ID == m.player.ID {
			return
		}
		if p.Sprite == nil {
			return
		}
		frameIdx := 0
		if p.Moving {
			frameIdx = p.Sprite.WalkFrame(p.Direction, m.animFrame)
		}
		m.renderSprite(p.X, p.Y, p.Sprite, frameIdx)
	})

	// Self
	frameIdx := 0
	if m.player.Sprite != nil {
		if m.player.Moving {
			frameIdx = m.player.Sprite.WalkFrame(m.player.Direction, m.animFrame)
		}
		m.renderSprite(m.player.X, m.player.Y, m.player.Sprite, frameIdx)
	}
}

func (m *GameModel) renderSprite(worldX, worldY int, s *sprite.Sprite, frameIdx int) {
	if s == nil || frameIdx >= len(s.Frames) {
		return
	}
	frame := s.Frames[frameIdx]
	for sy, row := range frame.Cells {
		for sx, cell := range row {
			if cell == nil || cell.Char == " " {
				continue
			}
			screenX := worldX + sx - 1 - m.camX
			screenY := worldY + sy - 2 - m.camY
			if screenX >= 0 && screenX < m.width && screenY >= 0 && screenY < m.height {
				c := *cell
				if c.BG == "" && m.buffer[screenY][screenX] != nil {
					c.BG = m.buffer[screenY][screenX].BG
				}
				m.buffer[screenY][screenX] = &c
			}
		}
	}
}

func (m *GameModel) renderChatHistory(lines []string) []string {
	msgs, newOffset := m.world.Messages.Since(m.msgOffset)
	m.msgOffset = newOffset

	for _, msg := range msgs {
		line := fmt.Sprintf("%s: %s", msg.PlayerName, msg.Text)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))
		rendered := style.Render(line)
		// push existing lines up
		if m.height > 2 {
			copy(lines[2:], lines[1:m.height-1])
			lines[1] = rendered + strings.Repeat(" ", m.width-visualWidth(rendered))
			if len(lines[1]) > m.width {
				lines[1] = lines[1][:m.width]
			}
		}
	}

	return lines
}

func (m *GameModel) renderChatOverlay(lines []string) []string {
	// chat input at bottom
	inputLine := m.height - 2
	prompt := "say: "
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff88"))
	rendered := style.Render(prompt + m.chatInput.View())
	lines[inputLine] = rendered + strings.Repeat(" ", m.width-visualWidth(rendered))

	// show recent messages above input
	msgs, _ := m.world.Messages.Since(0)
	start := len(msgs) - 5
	if start < 0 {
		start = 0
	}
	for i, msg := range msgs[start:] {
		line := fmt.Sprintf("%s: %s", msg.PlayerName, msg.Text)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa"))
		rendered := style.Render(line)
		y := inputLine - 5 + i
		if y >= 0 {
			lines[y] = rendered + strings.Repeat(" ", m.width-visualWidth(rendered))
		}
	}

	return lines
}

// --- Rendering utilities ---

func renderCellRow(cells []*shared.Cell) string {
	type seg struct{ chars, fg, bg string }
	var segs []seg
	for _, c := range cells {
		if c == nil {
			continue
		}
		if len(segs) > 0 && segs[len(segs)-1].fg == c.FG && segs[len(segs)-1].bg == c.BG {
			segs[len(segs)-1].chars += c.Char
		} else {
			segs = append(segs, seg{c.Char, c.FG, c.BG})
		}
	}
	var sb strings.Builder
	for _, s := range segs {
		sb.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color(s.fg)).
			Background(lipgloss.Color(s.bg)).
			Render(s.chars))
	}
	return sb.String()
}

func centerText(s string, width int) string {
	vw := visualWidth(s)
	if vw >= width {
		return s
	}
	return strings.Repeat(" ", (width-vw)/2) + s
}

func visualWidth(s string) int {
	w := 0
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
			continue
		}
		w++
	}
	return w
}

// --- Update ---

func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.sized = true

	case tea.KeyPressMsg:
		cmds = append(cmds, m.handleKeyPress(msg)...)

	case tea.KeyReleaseMsg:
		m.handleKeyRelease(msg)

	case TickMsg:
		if m.screenType == ui.ScreenGame {
			m.advanceAnim()
			if m.player.Moving && time.Since(m.lastInputTime) > 750*time.Millisecond {
				m.player.Moving = false
			}
			if m.player.Moving {
				m.movePlayer()
			}
		}
		return m, m.tickCmd()
	}

	return m, tea.Batch(cmds...)
}

func (m GameModel) tickCmd() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg { return TickMsg{} })
}

// --- Key handling ---

func (m *GameModel) handleKeyPress(msg tea.KeyPressMsg) []tea.Cmd {
	key := msg.String()

	if key == "ctrl+c" {
		m.world.RemovePlayer(m.player.ID)
		return []tea.Cmd{tea.Quit}
	}

	if m.chatActive {
		return m.handleChatKey(msg)
	}

	if key == "enter" {
		m.chatActive = true
		m.chatInput.Focus()
		m.chatInput.SetValue("")
		return nil
	}

	switch m.screenType {
	case ui.ScreenStart:
		return m.handleTitleKey(msg)
	case ui.ScreenGame:
		m.handleGameKeyPress(msg)
	}

	return nil
}

func (m *GameModel) handleChatKey(msg tea.KeyPressMsg) []tea.Cmd {
	key := msg.String()
	switch key {
	case "enter":
		text := m.chatInput.Value()
		if text != "" {
			name := m.player.Name
			if name == "" {
				name = "anon"
			}
			m.world.BroadcastChat(name, text)
		}
		m.chatActive = false
		m.chatInput.Blur()
		return nil
	case "esc":
		m.chatActive = false
		m.chatInput.Blur()
		return nil
	default:
		var cmd tea.Cmd
		m.chatInput, cmd = m.chatInput.Update(msg)
		if cmd != nil {
			return []tea.Cmd{cmd}
		}
	}
	return nil
}

func (m *GameModel) handleTitleKey(msg tea.KeyPressMsg) []tea.Cmd {
	key := msg.String()

	switch key {
	case "enter":
		name := m.nameInput.Value()
		if name != "" {
			m.player.Name = name
		}
		m.player.HeadChar = ui.HeadOptions[m.headIndex]
		m.player.Sprite = sprite.NewPlayerSprite(m.player.HeadChar, m.player.HeadFG, m.player.HeadBG)
		m.world.AddPlayer(m.player)
		m.screenType = ui.ScreenGame
		return nil

	case "left":
		m.headIndex = (m.headIndex - 1 + len(ui.HeadOptions)) % len(ui.HeadOptions)
		return nil

	case "right":
		m.headIndex = (m.headIndex + 1) % len(ui.HeadOptions)
		return nil

	default:
		var cmd tea.Cmd
		m.nameInput, cmd = m.nameInput.Update(msg)
		if cmd != nil {
			return []tea.Cmd{cmd}
		}
	}

	return nil
}

func (m *GameModel) handleGameKeyPress(msg tea.KeyPressMsg) {
	key := msg.String()

	switch key {
	case "w", "up":
		m.player.Direction = shared.Up
	case "s", "down":
		m.player.Direction = shared.Down
	case "a", "left":
		m.player.Direction = shared.Left
	case "d", "right":
		m.player.Direction = shared.Right
	default:
		return
	}

	m.shiftHeld = msg.Mod.Contains(tea.ModShift)

	if time.Since(m.lastInputTime) > 750*time.Millisecond {
		m.player.Moving = true
	}
	m.lastInputTime = time.Now()
}

func (m *GameModel) movePlayer() {
	var dx, dy int
	switch m.player.Direction {
	case shared.Up:
		dy = -1
	case shared.Down:
		dy = 1
	case shared.Left:
		dx = -1
	case shared.Right:
		dx = 1
	}

	if m.shiftHeld {
		dx *= 2
		dy *= 2
	}

	newX := m.player.X + dx
	newY := m.player.Y + dy
	if m.world.IsWalkable(newX, newY) {
		m.player.X = newX
		m.player.Y = newY
	} else if m.shiftHeld {
		singleX := m.player.X + dx/2
		singleY := m.player.Y + dy/2
		if m.world.IsWalkable(singleX, singleY) {
			m.player.X = singleX
			m.player.Y = singleY
		}
	}
}

func (m *GameModel) handleKeyRelease(msg tea.KeyReleaseMsg) {
	key := msg.String()
	switch key {
	case "w", "up", "s", "down", "a", "left", "d", "right":
		m.lastInputTime = time.Now()
	case "shift":
		m.shiftHeld = false
	}
}

// --- Animation ---

func (m *GameModel) advanceAnim() {
	if m.player.Moving && m.player.Sprite != nil {
		m.animTick++
		if m.animTick >= m.player.Sprite.FrameRate {
			m.animTick = 0
			m.animFrame = (m.animFrame + 1) % m.player.Sprite.WalkCycleLen(m.player.Direction)
		}
	} else {
		m.animFrame = 0
		m.animTick = 0
	}
}

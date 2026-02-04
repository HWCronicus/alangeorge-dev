package models

import (
	"fmt"
	"os"
	"strings"
	"time"

	gradient "github.com/HWCronicus/ssh-resume/src/utils"
	"github.com/Nomadcxx/sysc-Go/animations"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type view int

const (
	splashView view = iota
	aboutMeView
	skillsView
	workExperienceView
	projectsView
	contactInfoView
)

var tabs = []string{"About Me", "Skills", "Work Experience", "Projects", "Contact Info"}

type splashTimeoutMsg struct{}

type Model struct {
	Animation      animations.Animation
	AsciiArt       string
	CurrentView    view
	Cursor         int
	MinHeight      int
	MinWidth       int
	SetHeight      int
	SetWidth       int
	TerminalHeight int
	TerminalWidth  int
	Viewport       viewport.Model
	Ready          bool
	renderer       *lipgloss.Renderer
}

func (m Model) highlightColor() lipgloss.Color {
	return lipgloss.Color("208")
}

func (m Model) titleStyle() lipgloss.Style {
	return m.renderer.NewStyle().
		Border(lipgloss.Border{Bottom: "─"}, true).
		BorderForeground(lipgloss.Color("208")).
		Foreground(lipgloss.Color("208"))
}

func (m Model) aboutStyle() lipgloss.Style {
	return m.renderer.NewStyle().
		Foreground(lipgloss.Color("241"))
}

func (m Model) helpStyle() lipgloss.Style {
	return m.renderer.NewStyle().
		Foreground(lipgloss.Color("241"))
}

func (m Model) tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func (m Model) inactiveTabStyle() lipgloss.Style {
	return m.renderer.NewStyle().
		Border(m.tabBorderWithBottom("┴", "─", "┴"), true).
		BorderForeground(m.highlightColor()).
		Padding(0, 1)
}

func (m Model) activeTabStyle() lipgloss.Style {
	return m.renderer.NewStyle().
		Border(m.tabBorderWithBottom("┘", " ", "└"), true).
		BorderForeground(m.highlightColor()).
		Padding(0, 1).
		Foreground(lipgloss.Color("208")).
		Underline(true).
		Bold(true)
}

func (m Model) windowStyle() lipgloss.Style {
	return m.renderer.NewStyle().
		BorderForeground(m.highlightColor()).
		Padding(2, 2).
		Align(lipgloss.Left).
		Border(lipgloss.NormalBorder()).
		UnsetBorderTop().
		UnsetBorderBottom()
}

func (m Model) infoStyle() lipgloss.Style {
	return m.renderer.NewStyle().
		BorderForeground(m.highlightColor()).
		Padding(0, 1)
}

func InitialModel(minHeight, minWidth, terminalHeight, terminalWidth int, renderer *lipgloss.Renderer) Model {
	asciiArt, err := os.ReadFile("assets/alan.txt")
	if err != nil {
		asciiArt = []byte("AlanGeorge.Dev")
	}

	return Model{
		CurrentView:    splashView,
		Cursor:         0,
		MinHeight:      minHeight,
		MinWidth:       minWidth,
		SetHeight:      max(terminalHeight-40, minHeight),
		SetWidth:       max(terminalWidth-40, minWidth),
		TerminalHeight: terminalHeight,
		TerminalWidth:  terminalWidth,
		AsciiArt:       string(asciiArt),
		Animation:      nil,
		renderer:       renderer,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return splashTimeoutMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case splashTimeoutMsg:
		if m.CurrentView == splashView {
			if m.TerminalWidth >= m.MinWidth && m.TerminalHeight >= m.MinHeight {
				if m.Animation != nil {
					m.Animation.Update()
				}
				return m, tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
					return splashTimeoutMsg{}
				})
			}
			return m, tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
				return splashTimeoutMsg{}
			})
		}
		m.CurrentView = aboutMeView
		m.loadViewportContent()

	case tea.WindowSizeMsg:
		// oldWidth := m.TerminalWidth
		// oldHeight := m.Height

		// m.TerminalWidth = msg.Width
		// m.Width = min(msg.Width, 150)
		// m.Height = min(msg.Height, 50)

		m.TerminalHeight = msg.Height
		m.TerminalWidth = msg.Width
		m.SetHeight = max(m.TerminalHeight, m.MinHeight)
		m.SetWidth = max(m.TerminalWidth-40, m.MinWidth)

		if m.CurrentView == splashView {
			if m.TerminalWidth >= m.MinWidth && m.TerminalHeight >= m.MinHeight && m.Animation == nil {
				config := animations.BeamTextConfig{
					Width:                m.TerminalWidth,
					Height:               m.TerminalHeight,
					Text:                 m.AsciiArt,
					Auto:                 false,
					Display:              true,
					BeamRowSymbols:       []rune{'▂', '▁', '_'},
					BeamColumnSymbols:    []rune{'▌', '▍', '▎', '▏'},
					BeamDelay:            2,
					BeamRowSpeedRange:    [2]int{20, 80},
					BeamColumnSpeedRange: [2]int{15, 30},
					BeamGradientStops:    []string{"#ffffff", "#ff7300", "#ff9933"},
					BeamGradientSteps:    5,
					BeamGradientFrames:   1,
					FinalGradientStops:   []string{"#666666", "#ff7300", "#ff9933"},
					FinalGradientSteps:   8,
					FinalGradientFrames:  1,
					FinalWipeSpeed:       3,
				}
				m.Animation = animations.NewBeamTextEffect(config)

			}
		}

		if !m.Ready {
			contentWidth := m.SetWidth - 20
			viewportWidth := contentWidth - 4
			viewportHeight := m.SetHeight - 30
			m.Viewport = viewport.New(viewportWidth, viewportHeight)
			m.Viewport.YPosition = 0
			m.Ready = true
		} else {
			contentWidth := m.SetWidth - 20
			viewportWidth := contentWidth - 4
			viewportHeight := m.SetHeight - 30
			m.Viewport.Width = viewportWidth
			m.Viewport.Height = viewportHeight
		}
		return m, nil

	case tea.KeyMsg:
		if m.CurrentView == splashView {
			m.CurrentView = aboutMeView
			m.Cursor = 0
			m.loadViewportContent()
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "left", "h":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "right", "l":
			if m.Cursor < 4 {
				m.Cursor++
			}

		case "enter":
			m.CurrentView = view(m.Cursor + 1)
			m.loadViewportContent()

		case "b", "backspace", "esc":
			if m.CurrentView != aboutMeView {
				m.CurrentView = aboutMeView
				m.loadViewportContent()
			}

		case "1", "2", "3", "4", "5":
			selection := int(msg.String()[0] - '1')
			if selection >= 0 && selection < 5 {
				m.Cursor = selection
				m.CurrentView = view(selection + 1)
				m.loadViewportContent()
			}
		}
	}

	if m.CurrentView != splashView && m.Ready {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.TerminalWidth < m.MinWidth || m.TerminalHeight < m.MinHeight {
		return m.renderResizeMessage()
	}

	if m.CurrentView == splashView {
		return m.renderSplashScreen()
	}

	return m.RenderView()
}

func (m Model) renderSplashScreen() string {
	if m.Animation == nil {
		splashStyle := m.renderer.NewStyle().
			Foreground(lipgloss.Color("208")).
			Bold(true).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Width(m.TerminalWidth).
			Height(m.TerminalHeight)
		return splashStyle.Render(string(m.AsciiArt))
	}

	animationOutput := m.Animation.Render()

	promptText := m.renderer.NewStyle().
		Foreground(lipgloss.Color("241")).
		Italic(true).
		AlignHorizontal(lipgloss.Center).
		Width(m.TerminalWidth).
		Render("\n\nPress Enter to continue...")

	combined := animationOutput + promptText

	splashStyle := m.renderer.NewStyle().
		Width(m.TerminalWidth).
		Height(m.TerminalHeight)

	return splashStyle.Render(combined)
}

func (m Model) renderResizeMessage() string {
	messageStyle := m.renderer.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("208")).
		Align(lipgloss.Center).
		Width(m.TerminalWidth).
		Height(m.TerminalHeight)

	message := fmt.Sprintf(
		"Your Terminal window is too small.\n\n"+
			"For an optimal experience, please resize your terminal window.\n\n"+
			"Current size: %dx%d\n"+
			"Minimum size: %dx%d\n\n",
		m.TerminalWidth, m.TerminalHeight,
		m.MinWidth, m.MinHeight,
	)

	return messageStyle.Render(message)
}

func (m Model) RenderMainTitle() string {
	return m.titleStyle().Render(string(m.AsciiArt))
}

func (m Model) RenderTabs(content string) string {
	var renderedTabs []string
	for i, t := range tabs {
		var style lipgloss.Style
		if i == m.Cursor {
			style = m.activeTabStyle()
		} else {
			style = m.inactiveTabStyle()
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	contentWidth := m.SetWidth - 20
	totalWidth := contentWidth + m.windowStyle().GetHorizontalFrameSize() - 4

	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	tabsWidth := lipgloss.Width(tabRow)

	leftGap := (totalWidth - tabsWidth) / 2
	rightGap := totalWidth - tabsWidth - leftGap

	var leftBorder string
	if leftGap > 0 {
		leftBorder = m.renderer.NewStyle().
			Foreground(m.highlightColor()).
			Render("┌" + strings.Repeat("─", leftGap-1))
	}

	var rightBorder string
	if rightGap > 0 {
		rightBorder = m.renderer.NewStyle().
			Foreground(m.highlightColor()).
			Render(strings.Repeat("─", rightGap-1) + "┐")
	}

	row := lipgloss.JoinHorizontal(lipgloss.Bottom, leftBorder, tabRow, rightBorder)

	doc := strings.Builder{}
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(m.windowStyle().Width(contentWidth).Render(content))

	footer := m.FooterView()
	if footer != "" {
		doc.WriteString("\n")
		doc.WriteString(footer)
	}

	return doc.String()
}

func (m Model) RenderHelp() string {
	return m.helpStyle().Render("←/→: navigate • ↑/↓: scroll • 1-5: quick select • enter: select • q: quit")
}

func (m Model) FooterView() string {
	contentWidth := m.SetWidth - 20

	info := m.infoStyle().Render(fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100))
	infoWidth := lipgloss.Width(info)

	leftLineWidth := (contentWidth - infoWidth) / 2
	if leftLineWidth < 1 {
		leftLineWidth = 1
	}

	rightLineWidth := contentWidth - infoWidth - leftLineWidth
	if rightLineWidth < 1 {
		rightLineWidth = 1
	}

	leftLine := m.renderer.NewStyle().
		Foreground(m.highlightColor()).
		Render("└" + strings.Repeat("─", leftLineWidth-1) + "┤")

	rightLine := m.renderer.NewStyle().
		Foreground(m.highlightColor()).
		Render("├" + strings.Repeat("─", rightLineWidth-1) + "┘")

	return lipgloss.JoinHorizontal(lipgloss.Center, leftLine, info, rightLine)
}

func (m Model) RenderView() string {
	if m.TerminalWidth < m.MinWidth || m.TerminalHeight < m.MinHeight {
		return m.renderResizeMessage()
	}

	title := m.RenderMainTitle()
	about := m.aboutStyle().Render("Welcome to Terminal based version of AlanGeorge.Dev, navigate through the sections to learn more about me.")

	var viewContent string
	if m.Ready {
		viewContent = m.Viewport.View()
		if viewContent == "" {
			viewContent = "Press a number key (1-5) or use arrow keys and Enter to select a section."
		}
	} else {
		viewContent = "Initializing..."
	}

	contentStyle := m.renderer.NewStyle().
		Padding(0, 2).
		Width(m.SetWidth - 15)

	middleContent := contentStyle.Render(viewContent)

	var tabs string

	tabs = m.RenderTabs(middleContent)

	help := m.RenderHelp()

	topContent := fmt.Sprintf("%s\n\n%s\n\n", title, about)

	top := m.renderer.NewStyle().
		Width(m.SetWidth - 8).
		AlignHorizontal(lipgloss.Center).
		Render(topContent)

	availableHeight := m.SetHeight - 8
	topContentHeight := lipgloss.Height(topContent)
	helpHeight := lipgloss.Height(help)
	spacerHeight := availableHeight - topContentHeight - helpHeight + 4
	if spacerHeight < 0 {
		spacerHeight = 0
	}

	middle := m.renderer.NewStyle().
		Height(spacerHeight - 10).
		Width(m.SetWidth - 8).
		PaddingLeft(3).
		AlignHorizontal(lipgloss.Center).
		Render(tabs)

	bottom := m.renderer.NewStyle().
		Width(m.SetWidth - 8).
		AlignHorizontal(lipgloss.Center).
		Render(help)

	joinedVert := lipgloss.JoinVertical(lipgloss.Top, top, middle, bottom)

	bordered := gradient.RenderGradientBorder("#ff7300", "#666666", joinedVert, m.SetWidth, m.SetHeight, m.renderer)

	return m.renderer.NewStyle().
		Width(m.TerminalWidth).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(bordered)
}

func (m *Model) loadViewportContent() tea.Cmd {
	if !m.Ready {
		return nil
	}

	var filename string
	switch m.CurrentView {
	case aboutMeView:
		filename = "content/about-me.md"
	case skillsView:
		filename = "content/skills.md"
	case workExperienceView:
		filename = "content/work-experience.md"
	case projectsView:
		filename = "content/projects.md"
	case contactInfoView:
		filename = "content/contact-info.md"
	default:
		return nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		errorMsg := fmt.Sprintf("File not found: %s\n\nError: %s", filename, err.Error())
		m.Viewport.SetContent(errorMsg)
		return nil
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONFile("assets/md_styles.json"),
		glamour.WithWordWrap(m.Viewport.Width-2),
		glamour.WithColorProfile(termenv.TrueColor),
	)
	if err != nil {
		m.Viewport.SetContent(string(content))
		return nil
	}

	renderedContent, err := renderer.Render(string(content))
	if err != nil {
		m.Viewport.SetContent(string(content))
		return nil
	}

	m.Viewport.SetContent(renderedContent)
	m.Viewport.GotoTop()
	return nil
}

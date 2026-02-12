package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

func makeGradientStyles(colorA, colorB string, steps int) (s []lipgloss.Color) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0; i < steps; i++ {
		c := cA.BlendLuv(cB, float64(i)/float64(steps))
		s = append(s, lipgloss.Color(colorToHex(c)))
	}
	return
}

func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}

func RenderGradientBorder(colorA, colorB, content string, width, height int, renderer *lipgloss.Renderer) string {
	if width < 6 || height < 6 {
		return content
	}

	innerWidth := width - 2
	innerHeight := height - 2

	lines := renderer.NewStyle().
		Width(innerWidth).
		Height(innerHeight).
		Render(content)

	contentLines := strings.Split(lines, "\n")
	totalLines := len(contentLines)

	var result strings.Builder

	topColor := makeGradientStyles(colorA, colorB, totalLines)[0]
	result.WriteString(renderer.NewStyle().Foreground(topColor).Render("╔" + strings.Repeat("═", innerWidth) + "╗"))

	for i, line := range contentLines {
		color := makeGradientStyles(colorA, colorB, totalLines)[i]
		borderStyle := renderer.NewStyle().Foreground(color)
		paddedLine := renderer.NewStyle().Width(innerWidth).Render(line)
		result.WriteString("\n" + borderStyle.Render("║") + paddedLine + borderStyle.Render("║"))
	}

	result.WriteString("\n")
	bottomColor := makeGradientStyles(colorA, colorB, totalLines)[totalLines-1]
	result.WriteString(renderer.NewStyle().Foreground(bottomColor).Render("╚" + strings.Repeat("═", innerWidth) + "╝"))

	return result.String()
}

// Package theme shares the slidev-theme-ikura-hamu dark palette across the UI.
package theme

import (
	"charm.land/bubbles/v2/help"
	"charm.land/glamour/v2"
	"charm.land/glamour/v2/ansi"
	"charm.land/glamour/v2/styles"
	"charm.land/lipgloss/v2"
)

// Source: https://github.com/ikura-hamu/slidev-theme-ikura-hamu/blob/main/styles/layout.css
const (
	Primary    = "#0f71f2"
	Secondary  = "#c7d2fe"
	Background = "#071018"
	Accent     = "#eee8a9"
	Text       = "#f4f7fb"
	Inverse    = "#ffffff"
)

func HelpStyles() help.Styles {
	key := lipgloss.NewStyle().Foreground(lipgloss.Color(Accent))
	desc := lipgloss.NewStyle().Foreground(lipgloss.Color(Secondary))
	separator := lipgloss.NewStyle().Foreground(lipgloss.Color(Primary))
	return help.Styles{
		ShortKey: key, ShortDesc: desc, ShortSeparator: separator,
		FullKey: key, FullDesc: desc, FullSeparator: separator,
		Ellipsis: separator,
	}
}

func MarkdownStyle() ansi.StyleConfig {
	s := styles.DarkStyleConfig
	s.Document.Color = colorString(Text)
	s.Document.BackgroundColor = colorString(Background)
	s.Heading.Color = colorString(Primary)
	s.H1.Color = colorString(Inverse)
	s.H1.BackgroundColor = colorString(Primary)
	s.H6.Color = colorString(Secondary)
	s.Strong.Color = colorString(Primary)
	s.BlockQuote.Color = colorString(Secondary)
	s.HorizontalRule.Color = colorString(Secondary)
	s.Link.Color = colorString(Primary)
	s.LinkText.Color = colorString(Secondary)
	s.Image.Color = colorString(Primary)
	s.ImageText.Color = colorString(Secondary)
	s.Code.Color = colorString(Accent)
	s.Code.BackgroundColor = colorString(Background)
	s.CodeBlock.Color = colorString(Text)
	s.CodeBlock.BackgroundColor = colorString(Background)
	s.CodeBlock.Theme = ""
	s.CodeBlock.Chroma = &ansi.Chroma{
		Text:          ansi.StylePrimitive{Color: colorString(Text)},
		Background:    ansi.StylePrimitive{BackgroundColor: colorString(Background)},
		Comment:       ansi.StylePrimitive{Color: colorString(Secondary)},
		Keyword:       ansi.StylePrimitive{Color: colorString(Primary)},
		NameFunction:  ansi.StylePrimitive{Color: colorString(Secondary)},
		LiteralString: ansi.StylePrimitive{Color: colorString(Accent)},
		LiteralNumber: ansi.StylePrimitive{Color: colorString(Accent)},
	}
	return s
}

func RenderMarkdown(content string, width int) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(MarkdownStyle()),
		glamour.WithWordWrap(max(1, width)),
	)
	if err != nil {
		return "", err
	}
	return renderer.Render(content)
}

func colorString(value string) *string { return &value }

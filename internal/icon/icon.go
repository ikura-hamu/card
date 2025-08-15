package icon

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	tea "github.com/charmbracelet/bubbletea"
	"go.ikura-hamu.work/card/internal/common/size"
)

const iconURL = "https://avatars.githubusercontent.com/u/104292023"

type Model struct {
	size     size.Size
	filePath string
}

type imageFetchedMsg struct {
	filePath string
}

func fetchImage() (msg tea.Msg) {
	resp, err := http.Get(iconURL)
	if err != nil {
		return fmt.Errorf("fetch icon: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			msg = fmt.Errorf("close resp.Body: %w", err)
		}
	}()

	f, err := os.CreateTemp("", TempFilePrefix())
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			msg = fmt.Errorf("close temp file: %w", err)
		}
	}()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("copy resp.Body to file: %w", err)
	}

	return imageFetchedMsg{filePath: f.Name()}
}

func NewModel() Model { return Model{} }

func (m Model) Name() string {
	return "Icon"
}

func (m Model) Init() tea.Cmd {
	return fetchImage
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size.Width = msg.Width
		m.size.Height = msg.Height
	case imageFetchedMsg:
		m.filePath = msg.filePath
	}
	return m, nil
}

func (m Model) View() string {
	flags := aic_package.DefaultFlags()

	flags.Height = m.size.Height
	flags.Colored = true

	if m.filePath == "" {
		return "Loading icon..."
	}

	icon, err := aic_package.Convert(m.filePath, flags)
	if err != nil {
		return fmt.Errorf("failed to convert icon: %w", err).Error()
	}

	return icon
}

// TempFilePrefix returns the prefix used for icon temporary files for this process.
// Including PID prevents accidental removal of other running instances' temp files.
func TempFilePrefix() string { return fmt.Sprintf("ikura-hamu-card-%d-", os.Getpid()) }

func CleanupTempIcons() {
	// Clean up icon temp files created by this process after program exit.
	// We match files in the system temp dir with our PID-specific prefix.
	pattern := filepath.Join(os.TempDir(), TempFilePrefix()+"*")
	matches, _ := filepath.Glob(pattern)
	for _, f := range matches {
		_ = os.Remove(f)
	}
}

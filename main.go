package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/devguilhermeriberiiro/tekton/ui"
)

type Layer struct {
	Name        string
	Description string
}

type ArchTemplate struct {
	Name        string
	Description string
	Layers      []Layer
}

var templates = []ArchTemplate{
	{
		Name:        "Clean Architecture (Uncle Bob)",
		Description: "Separation by business rule. Framework-independent.",
		Layers: []Layer{
			{
				Name:        "domain",
				Description: "Pure entities and business rules. Zero external dependencies.",
			},
			{
				Name:        "usecase",
				Description: "Application use cases. Orchestrates the entities.",
			},
			{
				Name:        "repository",
				Description: "Data access interfaces (contracts, not implementations).",
			},
			{
				Name:        "infra/database",
				Description: "Concrete database implementations.",
			},
			{
				Name:        "infra/http",
				Description: "HTTP handlers, middleware, routers.",
			},
			{
				Name:        "config",
				Description: "Application configuration (env, flags).",
			},
			{
				Name:        "cmd",
				Description: "Application entry point.",
			},
		},
	},

	{
		Name:        "Hexagonal Architecture (Ports & Adapters)",
		Description: "Domain at the center. Ports define contracts, adapters implement them.",
		Layers: []Layer{
			{
				Name:        "core/domain",
				Description: "Pure domain model.",
			},
			{
				Name:        "core/ports/inbound",
				Description: "Interfaces the outside world uses to enter the application.",
			},
			{
				Name:        "core/ports/outbound",
				Description: "Interfaces the application uses to reach out (DB, external APIs).",
			},
			{
				Name:        "adapters/primary/http",
				Description: "HTTP adapter (inbound).",
			},
			{
				Name:        "adapters/primary/grpc",
				Description: "gRPC adapter (inbound).",
			},
			{
				Name:        "adapters/secondary/postgres",
				Description: "PostgreSQL adapter (outbound).",
			},
			{
				Name:        "adapters/secondary/redis",
				Description: "Redis adapter (outbound).",
			},
			{
				Name:        "cmd",
				Description: "Bootstrap and dependency wiring.",
			},
		},
	},

	{
		Name:        "Standard Go Layout (github.com/golang-standards)",
		Description: "Conventional Go community layout for larger projects.",
		Layers: []Layer{
			{
				Name:        "cmd",
				Description: "Main application binaries.",
			},
			{
				Name:        "internal",
				Description: "Private application logic.",
			},
			{
				Name:        "pkg",
				Description: "Public libraries (can be imported by other projects).",
			},
			{
				Name:        "api",
				Description: "API specifications (OpenAPI, Protobuf, GraphQL schemas).",
			},
			{
				Name:        "configs",
				Description: "Configuration templates and default values.",
			},
			{
				Name:        "scripts",
				Description: "Build, install, and analysis scripts.",
			},
			{
				Name:        "test",
				Description: "Fixtures and integration tests.",
			},
		},
	},
}

type screen int

const (
	screenWelcome screen = iota
	screenGithub
	screenName
	screenModule
	screenTemplate
	screenPreview
	screenDone
)

type Model struct {
	screen screen

	workDir      string
	githubUser   string
	projectName  string
	modulePath   string
	selectedTmpl int

	inputBuffer string
	inputError  string

	generated bool
	genError  string
	genPaths  []string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		return m, nil
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		if m.screen > screenWelcome {
			m.screen--
			switch m.screen {
			case screenGithub:
				m.inputBuffer = m.githubUser
			case screenName:
				m.inputBuffer = m.projectName
			case screenModule:
				m.inputBuffer = m.modulePath
			default:
				m.inputBuffer = ""
			}
			m.inputError = ""
		}
		return m, nil
	}

	switch m.screen {

	case screenWelcome:
		if msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace {
			m.screen = screenGithub
		}

	case screenGithub:
		switch msg.Type {
		case tea.KeyEnter:
			github := strings.TrimSpace(m.inputBuffer)
			if github == "" {
				m.inputError = "GitHub username cannot be empty."
			} else if strings.ContainsAny(github, " /\\:*?\"<>|") {
				m.inputError = "Use only letters, numbers, hyphens, and underscores."
			} else {
				m.githubUser = github
				m.inputBuffer = ""
				m.inputError = ""
				m.screen = screenName
			}
		case tea.KeyBackspace:
			if len(m.inputBuffer) > 0 {
				m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
			}
		default:
			if msg.Type == tea.KeyRunes {
				m.inputBuffer += string(msg.Runes)
			}
		}

	case screenName:
		switch msg.Type {
		case tea.KeyEnter:
			name := strings.TrimSpace(m.inputBuffer)
			if name == "" {
				m.inputError = "Project name cannot be empty."
			} else if strings.ContainsAny(name, " /\\:*?\"<>|") {
				m.inputError = "Use only letters, numbers, hyphens, and underscores."
			} else {
				m.projectName = name
				m.inputBuffer = "github.com/" + m.githubUser + "/" + name
				m.inputError = ""
				m.screen = screenModule
			}
		case tea.KeyBackspace:
			if len(m.inputBuffer) > 0 {
				m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
			}
		default:
			if msg.Type == tea.KeyRunes {
				m.inputBuffer += string(msg.Runes)
			}
		}

	case screenModule:
		switch msg.Type {
		case tea.KeyEnter:
			mod := strings.TrimSpace(m.inputBuffer)
			if mod == "" {
				m.inputError = "Module path cannot be empty."
			} else if !strings.Contains(mod, "/") {
				m.inputError = "Use the format: github.com/user/project"
			} else {
				m.modulePath = mod
				m.inputBuffer = ""
				m.inputError = ""
				m.screen = screenTemplate
			}
		case tea.KeyBackspace:
			if len(m.inputBuffer) > 0 {
				m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
			}
		default:
			if msg.Type == tea.KeyRunes {
				m.inputBuffer += string(msg.Runes)
			}
		}

	case screenTemplate:
		switch msg.String() {
		case "j", "down":
			if m.selectedTmpl < len(templates)-1 {
				m.selectedTmpl++
			}
		case "k", "up":
			if m.selectedTmpl > 0 {
				m.selectedTmpl--
			}
		case "enter":
			m.screen = screenPreview
		}

	case screenPreview:
		switch msg.String() {
		case "enter", "y":
			paths, err := generateStructure(m.projectName, m.workDir, templates[m.selectedTmpl])
			if err != nil {
				m.genError = err.Error()
			} else {
				m.genPaths = paths
				m.generated = true
			}
			m.screen = screenDone
		case "n":
			m.screen = screenTemplate
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch m.screen {
	case screenWelcome:
		return welcome()
	case screenGithub:
		return input(m, "GitHub User", "github user", "Used to create the local git repository (.git)")
	case screenName:
		return input(m, "Project Name", "my-project", "This will be the name of the root folder created.")
	case screenModule:
		return input(m, "Go module path", "github.com/user/project", "Used in go.mod and internal imports.")
	case screenTemplate:
		return template(m)
	case screenPreview:
		return preview(m)
	case screenDone:
		return done(m)
	}
	return ""
}

func welcome() string {
	logo := `████████╗███████╗██╗  ██╗████████╗ ██████╗ ███╗   ██╗
╚══██╔══╝██╔════╝██║ ██╔╝╚══██╔══╝██╔═══██╗████╗  ██║
   ██║   █████╗  █████╔╝    ██║   ██║   ██║██╔██╗ ██║
   ██║   ██╔══╝  ██╔═██╗    ██║   ██║   ██║██║╚██╗██║
   ██║   ███████╗██║  ██╗   ██║   ╚██████╔╝██║ ╚████║
   ╚═╝   ╚══════╝╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═══╝`

	var b strings.Builder
	b.WriteString(ui.StyleBrand.Render(logo))
	b.WriteString("\n\n")
	b.WriteString(ui.StyleTitle.Render("  Clean Architecture Folder Generator for Go"))
	b.WriteString("\n")
	b.WriteString(ui.StyleMuted.Render("  Generates the folder structure for your next project."))
	b.WriteString("\n\n")
	b.WriteString(ui.StyleMuted.Render("  ─────────────────────────────────────────────────"))
	b.WriteString("\n\n")
	b.WriteString("  " + ui.StyleWarning.Render("→") + "  Clean Architecture  " + ui.StyleMuted.Render("(Uncle Bob)"))
	b.WriteString("\n")
	b.WriteString("  " + ui.StyleWarning.Render("→") + "  Hexagonal Architecture  " + ui.StyleMuted.Render("(Ports & Adapters)"))
	b.WriteString("\n")
	b.WriteString("  " + ui.StyleWarning.Render("→") + "  Standard Go Layout  " + ui.StyleMuted.Render("(golang-standards)"))
	b.WriteString("\n\n")
	b.WriteString(renderFooter("enter/space  start", "ctrl+c  quit"))
	return b.String()
}

func input(m Model, label, placeholder, hint string) string {
	var b strings.Builder

	b.WriteString(renderHeader(label))
	b.WriteString("\n")
	b.WriteString(ui.StyleMuted.Render("  " + hint))
	b.WriteString("\n\n")

	cursor := ui.StyleActive.Render("▌")
	inputText := m.inputBuffer
	if inputText == "" {
		inputText = ui.StyleMuted.Render(placeholder) + cursor
	} else {
		inputText = ui.StyleTitle.Render(inputText) + cursor
	}

	b.WriteString(ui.StylePrompt.Render("  › ") + inputText)
	b.WriteString("\n")

	if m.inputError != "" {
		b.WriteString("\n")
		b.WriteString("  " + ui.StyleError.Render("✗ "+m.inputError))
	}

	b.WriteString("\n\n")
	b.WriteString(renderFooter("enter  confirm", "esc  back", "ctrl+c  quit"))
	return b.String()
}

func template(m Model) string {
	var b strings.Builder

	b.WriteString(renderHeader("Choose architecture"))
	b.WriteString("\n")
	b.WriteString(ui.StyleMuted.Render("  Project: ") + ui.StyleActive.Render(m.projectName))
	b.WriteString("  " + ui.StyleMuted.Render("│") + "  ")
	b.WriteString(ui.StyleMuted.Render("Module: ") + ui.StyleMuted.Render(m.modulePath))
	b.WriteString("\n\n")

	for i, tmpl := range templates {
		prefix := "  "
		if i == m.selectedTmpl {
			prefix = ui.StyleActive.Render("  ▸ ")
			b.WriteString(prefix + ui.StyleActive.Render(tmpl.Name))
		} else {
			b.WriteString(prefix + "  " + ui.StyleTitle.Render(tmpl.Name))
		}
		b.WriteString("\n")
		b.WriteString("       " + ui.StyleMuted.Render(tmpl.Description))
		b.WriteString("\n\n")
	}

	b.WriteString(renderFooter("↑/k  up", "↓/j  down", "enter  select", "esc  back"))
	return b.String()
}

func preview(m Model) string {
	tmpl := templates[m.selectedTmpl]
	var b strings.Builder

	b.WriteString(renderHeader("Structure preview"))
	b.WriteString("\n")
	b.WriteString(ui.StyleMuted.Render("  Architecture: ") + ui.StyleActive.Render(tmpl.Name))
	b.WriteString("\n\n")

	tree := buildTree(m.projectName, tmpl)
	b.WriteString(ui.StyleBorderLeft.Render(tree))
	b.WriteString("\n\n")

	b.WriteString("  " + ui.StyleActive.Render("Generate this structure at ./"+m.projectName+"?"))
	b.WriteString("\n\n")
	b.WriteString(renderFooter("enter/y  generate", "n  back", "esc  back"))
	return b.String()
}

func done(m Model) string {
	var b strings.Builder
	b.WriteString(renderHeader("Done"))
	b.WriteString("\n")

	if m.genError != "" {
		b.WriteString(ui.StyleError.Render("  ✗ Error generating: " + m.genError))
		b.WriteString("\n\n")
		b.WriteString(renderFooter("ctrl+c  quit"))
		return b.String()
	}

	summary := fmt.Sprintf(
		"  Project: %s\n  Module:  %s\n  Folders: %d created",
		m.projectName, m.modulePath, len(m.genPaths),
	)

	b.WriteString(ui.StyleBox.Render(
		ui.StyleSuccess.Render("✓ Structure generated successfully!") + "\n\n" + summary,
	))
	b.WriteString("\n\n")

	b.WriteString(ui.StyleMuted.Render("  Folders created:"))
	b.WriteString("\n\n")
	for _, p := range m.genPaths {
		b.WriteString(ui.StyleMuted.Render("  ") + ui.StyleSuccess.Render("  ✓") + "  " + ui.StyleTitle.Render(p) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(ui.StyleMuted.Render("  Next steps:"))
	b.WriteString("\n\n")
	b.WriteString("  " + ui.StyleWarning.Render("1.") + "  cd " + ui.StyleActive.Render(m.workDir+m.projectName) + "\n")
	b.WriteString("  " + ui.StyleWarning.Render("2.") + "  go mod init " + ui.StyleActive.Render(m.modulePath) + "\n")
	b.WriteString("  " + ui.StyleWarning.Render("3.") + "  Start coding! 🚀" + "\n")

	b.WriteString("\n")
	b.WriteString(renderFooter("ctrl+c  quit"))
	return b.String()
}

func renderHeader(title string) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(ui.StyleBrand.Render("  TEKTON") + ui.StyleMuted.Render("  /  ") + ui.StyleTitle.Render(title))
	b.WriteString("\n")
	b.WriteString(ui.StyleMuted.Render("  " + strings.Repeat("─", 50)))
	b.WriteString("\n")
	return b.String()
}

func renderFooter(hints ...string) string {
	parts := make([]string, len(hints))
	for i, h := range hints {
		parts[i] = ui.StyleMuted.Render("  " + h)
	}
	return "\n" + strings.Join(parts, ui.StyleMuted.Render("  │")) + "\n"
}

func buildTree(projectName string, tmpl ArchTemplate) string {
	var b strings.Builder

	b.WriteString(ui.StyleActive.Render(projectName+"/") + "\n")

	for i, layer := range tmpl.Layers {
		isLast := i == len(tmpl.Layers)-1
		parts := strings.Split(layer.Name, "/")

		for depth, part := range parts {
			prefix := "│   "
			connector := "├── "
			if isLast && depth == len(parts)-1 {
				prefix = "    "
				connector = "└── "
			}

			indent := strings.Repeat(prefix, depth)

			if depth == len(parts)-1 {
				b.WriteString(
					ui.StyleMuted.Render(indent+connector) +
						ui.StyleTitle.Render(part+"/") +
						ui.StyleMuted.Render("   # "+layer.Description) +
						"\n",
				)
			} else {
				b.WriteString(
					ui.StyleMuted.Render(indent+connector) +
						ui.StyleTitle.Render(part+"/") +
						"\n",
				)
			}
		}
	}

	return b.String()
}

func generateStructure(projectName string, workDir string, tmpl ArchTemplate) ([]string, error) {
	var created []string
	root := filepath.Join(workDir, projectName)

	for _, layer := range tmpl.Layers {
		dirPath := filepath.Join(root, filepath.FromSlash(layer.Name))

		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			return created, fmt.Errorf("error creating %s: %w", dirPath, err)
		}
		created = append(created, dirPath)
	}
	return created, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tekton <path>")
		return
	}

	workDir := os.Args[1]
	if !strings.HasSuffix(workDir, "/") {
		workDir += "/"
	}

	initialModel := Model{
		screen:       screenWelcome,
		selectedTmpl: 0,
		workDir:      workDir,
	}

	p := tea.NewProgram(
		initialModel,
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

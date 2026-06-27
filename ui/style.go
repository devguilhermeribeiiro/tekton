package ui

import (
	// lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Paleta de cores — escolhida para parecer um terminal de desenvolvedor:
	// fundo escuro implícito, acentos em verde/cyan como shells modernos
	ColorPrimary   = lipgloss.Color("#00D7AF") // verde-turquesa (accent principal)
	ColorSecondary = lipgloss.Color("#5F87FF") // azul periwinkle (itens ativos)
	ColorMuted     = lipgloss.Color("#626262") // cinza escuro (texto secundário)
	ColorWarning   = lipgloss.Color("#FFD700") // dourado (avisos, destaques)
	ColorError     = lipgloss.Color("#FF5F5F") // vermelho suave (erros)
	ColorSuccess   = lipgloss.Color("#87FF5F") // verde limão (sucesso)
	ColorText      = lipgloss.Color("#D7D7D7") // branco suave (texto principal)

	// Estilos compostos
	StyleBrand = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	StyleTitle = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true)

	StyleSubtitle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true)

	StyleActive = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)

	StyleMuted = lipgloss.NewStyle().
			Foreground(ColorMuted)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	StyleError = lipgloss.NewStyle().
			Foreground(ColorError)

	StyleWarning = lipgloss.NewStyle().
			Foreground(ColorWarning)

	StylePrompt = lipgloss.NewStyle().
			Foreground(ColorPrimary)

	// Borda lateral esquerda — usada no painel de preview
	StyleBorderLeft = lipgloss.NewStyle().
			BorderLeft(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(ColorMuted).
			PaddingLeft(1)

	// Box de confirmação final
	StyleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)
)

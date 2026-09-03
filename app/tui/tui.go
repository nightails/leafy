// Package tui implements the main loop of the Leafy app in TUI mode.
// This includes user interactions and prompts.
package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/internal/style"
)

type step int

const (
	media = iota
	destination
	copying
	finished
)

type Model struct {
	version    string
	state      state
	currStep   step
	mediaList  list.Model
	destInput  textinput.Model
	copyEvents <-chan tea.Msg
	err        error
}

func New(ver string) Model {
	l := list.New([]list.Item{}, mediaItemDelegate{}, 0, 1)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	ti := textinput.New()
	ti.Placeholder = "Destination Path"

	return Model{
		version:    ver,
		state:      state{},
		currStep:   media,
		mediaList:  l,
		destInput:  ti,
		copyEvents: nil,
		err:        nil,
	}
}

func (m Model) Init() tea.Cmd {
	return initDevicesCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		default:
			if m.err == nil {
				var cmd tea.Cmd
				switch m.currStep {
				case media:
					m.mediaList, cmd = m.mediaList.Update(msg)
					return m, cmd
				case destination:
					m.destInput, cmd = m.destInput.Update(msg)
					return m, cmd
				case copying:
					return m, nil
				case finished:
					return m, nil
				}
			}
			return m, nil

		case "ctrl+c":
			if m.currStep == copying {
				return m, nil
			}
			// TODO: implement canceling
			return m, tea.Batch(
				removeDevicesCmd(m.state.devices),
				tea.Quit,
			)

		case " ":
			index := m.mediaList.Index()
			item := m.mediaList.SelectedItem()
			selectItem := item.(mediumItem)
			selectItem.selected = !selectItem.selected
			m.mediaList.SetItem(index, selectItem)
			return m, nil

		case "enter", "return":
			switch m.currStep {
			default:
				// do nothing while copying and finished
				return m, nil

			case media:
				for _, i := range m.mediaList.Items() {
					mi := i.(mediumItem)
					if mi.selected {
						m.state.media = append(m.state.media, mi.medium)
					}
				}
				if m.state.media == nil {
					return m, nil
				}
				m.currStep++
				m.destInput.Focus()
				return m, nil

			case destination:
				dest := m.destInput.Value()
				if dest == "" {
					return m, nil
				}

				dest, err := expandHome(dest)
				if err != nil {
					m.err = err
					return m, nil
				}

				// TODO: apply naming scheme with increment
				for i := range m.state.media {
					m.state.media[i].dest = filepath.Join(dest, m.state.media[i].name)
				}
				m.currStep++
				m.destInput.Blur()
				return m, copyMediaCmd(m.state.media)
			}
		}

	case errMsg:
		m.err = msg
		return m, nil

	case tea.WindowSizeMsg:
		m.mediaList.SetWidth(msg.Width)
		return m, nil

	case devicesMsg:
		if len(msg) == 0 {
			return m, nil
		}
		m.state.devices = msg
		return m, findMediaCmd(msg)

	case mediaMsg:
		if len(msg) == 0 {
			return m, nil
		}
		items := make([]list.Item, 0, len(msg))
		for _, i := range msg {
			items = append(items, mediumItem{i, false})
		}
		m.mediaList.SetItems(items)
		m.mediaList.SetHeight(len(items))
		return m, nil

	case copyStartedMsg:
		m.copyEvents = msg.Ch
		return m, listenCopyCmd(m.copyEvents)

	case copyProgressMsg:
		if msg.Index >= 0 && msg.Index < len(m.state.media) {
			m.state.media[msg.Index].copied = msg.Copied
			m.state.media[msg.Index].total = msg.Total
		}
		return m, listenCopyCmd(m.copyEvents)

	case copyDoneMsg:
		if msg.Index >= 0 && msg.Index < len(m.state.media) {
			m.state.media[msg.Index].copied = m.state.media[msg.Index].total
		}
		return m, listenCopyCmd(m.copyEvents)

	case copyErrorMsg:
		m.err = msg.Err
		return m, nil

	case copyFinishedMsg:
		m.currStep = finished
		return m, nil

	case deleteDoneMsg:
		m.mediaList.SetItems(nil)
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString(style.HeaderStyle.Render(fmt.Sprintf("Leafy %s", m.version)))
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(m.err.Error())
		return b.String()
	}

	if m.currStep >= media {
		b.WriteString("Found media:\n")
		b.WriteString(m.mediaList.View())
	}
	if m.currStep >= destination {
		b.WriteString("\n\n")
		b.WriteString("Destination Path:\n")
		b.WriteString(m.destInput.View())
	}
	if m.currStep >= copying {
		b.WriteString("\n\n")
		b.WriteString("Copying Progress:\n")
		for i := range m.state.media {
			percent := 0.0
			if m.state.media[i].total > 0 {
				percent = float64(m.state.media[i].copied) / float64(m.state.media[i].total) * 100
			}
			b.WriteString(fmt.Sprintf(
				"%s %s/%s (%.0f%%)\n",
				m.state.media[i].name,
				formatFileSize(m.state.media[i].copied),
				formatFileSize(m.state.media[i].total),
				percent,
			))
		}
	}
	if m.currStep == finished && len(m.state.media) == 0 {
		b.WriteString("Selected media has been deleted.")
	}
	body := style.TextStyle.Render(b.String())

	var h strings.Builder
	h.WriteString("\n\n")
	if m.currStep == copying {
		h.WriteString("[copying in progress: quitting disabled]")
	} else {
		h.WriteString("[j/k] up/down | [space] select | [enter] confirm | [ctrl+c] quit")
	}
	help := style.HelpTextStyle.Render(h.String())

	return body + help
}

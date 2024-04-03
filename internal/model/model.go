/*
	bqt - BigQuery TUI
	Copyright (C) 2024  Corbin Staaben<cstaaben@gmail.com>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// Package model defines the state for bqt.
package model

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cstaaben/bqt/internal/client"
	"github.com/cstaaben/bqt/internal/formatter"
)

// Model represents the state of the application.
type Model struct {
	// resources
	client    *client.Client
	formatter formatter.Formatter

	// UI elements
	QueryArea   *textarea.Model
	ResultsArea any

	// values
	Query        string
	QueryResults []QueryResult
}

// QueryResult is a map of column names to values.
type QueryResult map[string]any

// New creates a new Model.
func New(client *client.Client, formatter formatter.Formatter) *Model {
	area := textarea.New()
	area.Placeholder = "Enter a query..."
	area.ShowLineNumbers = true
	area.Cursor.Blink = true

	return &Model{
		client:       client,
		formatter:    formatter,
		QueryArea:    &area,
		ResultsArea:  nil, // TODO: add results area
		QueryResults: make([]QueryResult, 0),
	}
}

// Init sets the window title and initializes the model.
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("bqt"),
	)
}

// Update updates the model based on the message.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.QueryArea.SetWidth(msg.Width)
		m.QueryArea.SetHeight(msg.Height)
		return m, nil
	default:
		slog.Warn("unrecognized message type", slog.Any("msg_type", fmt.Sprintf("%T", msg)))
	}

	return m, nil
}

func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		if msg.String() == "q" && m.QueryArea.Focused() {
			return m.updateQueryArea(msg)
		}

		// quit when query area is not focused and either key combination is pressed
		return m, tea.Sequence(tea.ClearScreen, tea.Quit)
	case "i":
		// enable or disable insert mode
		if !m.QueryArea.Focused() {
			return m, m.QueryArea.Focus()
		}

		return m.updateQueryArea(msg)
	case "esc":
		if m.QueryArea.Focused() {
			m.QueryArea.Blur()
		}

		return m, nil
	default:
		return m.updateQueryArea(msg)
	}
}

func (m *Model) updateQueryArea(msg tea.Msg) (tea.Model, tea.Cmd) {
	area, cmd := m.QueryArea.Update(msg)
	m.QueryArea = &area
	return m, cmd
}

// View returns the view for the model.
func (m *Model) View() string {
	return m.QueryArea.View()
}

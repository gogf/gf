// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file implements a Bubble Tea single-select list used by
// PromptSingleSelection. It intentionally does NOT use charmbracelet/huh
// Select for long lists: huh pins the selected row to the top of its
// viewport (YOffset = selected) on every update once the option list is
// taller than the terminal, which makes every row above the cursor
// disappear while navigating. This picker keeps a sliding window that
// only scrolls when the cursor would leave the visible range, so rows
// above the cursor stay visible whenever they still fit.
//
// Section rows (SingleOption.Section == true) are rendered as headers and
// are skipped by cursor movement so users cannot "select" a group title.

package common

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	// defaultListWindow is the visible option rows when terminal height
	// is unknown. Title, filter, and help lines sit outside this budget.
	defaultListWindow = 12
	// minListWindow keeps at least a few options on screen on tiny terminals.
	minListWindow = 5
	// listChromeLines is the reserved vertical space for title, filter
	// status, blank separators and the help footer.
	listChromeLines = 6
)

var (
	listTitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	listHelpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	listFilterStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	listCursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
	listNormalStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	listSectionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
	listDimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// singleListItem is one row in the custom single-select list.
type singleListItem struct {
	value   string
	label   string
	section bool
}

// singleListModel is the Bubble Tea model for PromptSingleSelection.
type singleListModel struct {
	title       string
	items       []singleListItem
	filtered    []int // indexes into items
	cursor      int   // index into filtered
	windowStart int   // first visible filtered index
	windowSize  int   // visible filtered rows
	filter      string
	chosen      string
	cancelled   bool
	quitting    bool
	width       int
	termHeight  int
}

// newSingleListModel builds a model with the cursor on the first
// selectable (non-section) row.
func newSingleListModel(title string, options []SingleOption) singleListModel {
	items := make([]singleListItem, 0, len(options))
	for _, option := range options {
		label := strings.TrimSpace(option.Label)
		if label == "" {
			label = option.Value
		}
		items = append(items, singleListItem{
			value:   option.Value,
			label:   label,
			section: option.Section,
		})
	}
	model := singleListModel{
		title:      title,
		items:      items,
		windowSize: defaultListWindow,
		width:      80,
	}
	model.refilter()
	return model
}

func (m singleListModel) Init() tea.Cmd { return nil }

func (m singleListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.termHeight = msg.Height
		m.windowSize = listWindowSize(msg.Height)
		m.ensureVisible()
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.cancelled = true
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			if value, ok := m.selectedValue(); ok {
				m.chosen = value
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		case tea.KeyUp, tea.KeyCtrlP:
			m.moveCursor(-1)
			return m, nil
		case tea.KeyDown, tea.KeyCtrlN:
			m.moveCursor(1)
			return m, nil
		case tea.KeyPgUp:
			m.moveCursor(-m.windowSize)
			return m, nil
		case tea.KeyPgDown:
			m.moveCursor(m.windowSize)
			return m, nil
		case tea.KeyHome:
			m.cursor = firstSelectableFiltered(m.items, m.filtered)
			m.ensureVisible()
			return m, nil
		case tea.KeyEnd:
			m.cursor = lastSelectableFiltered(m.items, m.filtered)
			m.ensureVisible()
			return m, nil
		case tea.KeyBackspace:
			if m.filter != "" {
				r := []rune(m.filter)
				m.filter = string(r[:len(r)-1])
				m.refilter()
			}
			return m, nil
		case tea.KeyRunes:
			for _, r := range msg.Runes {
				if unicode.IsControl(r) {
					continue
				}
				m.filter += string(r)
			}
			m.refilter()
			return m, nil
		}
	}
	return m, nil
}

func (m singleListModel) View() string {
	if m.quitting {
		return ""
	}
	var builder strings.Builder
	builder.WriteString(listTitleStyle.Render(m.title))
	builder.WriteByte('\n')
	if m.filter != "" {
		builder.WriteString(listFilterStyle.Render("Filter: " + m.filter))
	} else {
		builder.WriteString(listHelpStyle.Render("Type to filter · ↑/↓ move · Enter confirm · Esc cancel"))
	}
	builder.WriteByte('\n')
	builder.WriteByte('\n')

	if len(m.filtered) == 0 {
		builder.WriteString(listDimStyle.Render("  (no matches)"))
		builder.WriteByte('\n')
		return builder.String()
	}

	end := m.windowStart + m.windowSize
	if end > len(m.filtered) {
		end = len(m.filtered)
	}
	if m.windowStart > 0 {
		builder.WriteString(listDimStyle.Render(fmt.Sprintf("  … %d more above", m.windowStart)))
		builder.WriteByte('\n')
	}
	for index := m.windowStart; index < end; index++ {
		item := m.items[m.filtered[index]]
		line := m.renderRow(item, index == m.cursor)
		builder.WriteString(line)
		builder.WriteByte('\n')
	}
	below := len(m.filtered) - end
	if below > 0 {
		builder.WriteString(listDimStyle.Render(fmt.Sprintf("  … %d more below", below)))
		builder.WriteByte('\n')
	}
	builder.WriteByte('\n')
	builder.WriteString(listHelpStyle.Render(fmt.Sprintf("%d/%d", m.positionLabel(), len(m.selectableCount()))))
	return builder.String()
}

func (m singleListModel) renderRow(item singleListItem, selected bool) string {
	if item.section {
		return "  " + listSectionStyle.Render(item.label)
	}
	if selected {
		return listCursorStyle.Render("▸ " + item.label)
	}
	return listNormalStyle.Render("  " + item.label)
}

func (m *singleListModel) selectedValue() (string, bool) {
	if m.cursor < 0 || m.cursor >= len(m.filtered) {
		return "", false
	}
	item := m.items[m.filtered[m.cursor]]
	if item.section {
		return "", false
	}
	return item.value, true
}

func (m *singleListModel) refilter() {
	query := strings.ToLower(strings.TrimSpace(m.filter))
	m.filtered = m.filtered[:0]
	if query == "" {
		for index := range m.items {
			m.filtered = append(m.filtered, index)
		}
	} else {
		matched := make([]bool, len(m.items))
		for index, item := range m.items {
			if item.section {
				if strings.Contains(strings.ToLower(item.label), query) {
					matched[index] = true
				}
				continue
			}
			if strings.Contains(strings.ToLower(item.label), query) || strings.Contains(strings.ToLower(item.value), query) {
				matched[index] = true
			}
		}
		// Keep a section header when any selectable row under it matches,
		// so the two-group layout remains readable while filtering.
		for index, item := range m.items {
			if !item.section {
				continue
			}
			for look := index + 1; look < len(m.items); look++ {
				if m.items[look].section {
					break
				}
				if matched[look] {
					matched[index] = true
					break
				}
			}
		}
		for index, ok := range matched {
			if ok {
				m.filtered = append(m.filtered, index)
			}
		}
	}
	m.cursor = firstSelectableFiltered(m.items, m.filtered)
	m.windowStart = 0
	m.ensureVisible()
}

func (m *singleListModel) moveCursor(delta int) {
	if len(m.filtered) == 0 {
		return
	}
	next := moveSelectableFiltered(m.items, m.filtered, m.cursor, delta)
	if next < 0 {
		return
	}
	m.cursor = next
	m.ensureVisible()
}

// ensureVisible scrolls the window just enough to keep the cursor inside
// the visible range. Unlike huh.Select it does not pin the cursor to the
// top row, so previously visible rows above the cursor remain on screen
// until the cursor would walk off the bottom (or top) edge.
//
// Section headers are not selectable, so a plain cursor clamp can leave a
// group title parked in "… N more above" when the user returns to the
// first item of that group. When the nearest section above the cursor
// still fits in the same window as the cursor, pull windowStart up so
// the title is shown again.
func (m *singleListModel) ensureVisible() {
	total := len(m.filtered)
	sectionStart := sectionStartForCursor(m.items, m.filtered, m.cursor)
	m.windowStart = clampListWindowWithSection(m.cursor, m.windowStart, m.windowSize, total, sectionStart)
}

func (m singleListModel) selectableCount() []int {
	out := make([]int, 0, len(m.filtered))
	for _, index := range m.filtered {
		if !m.items[index].section {
			out = append(out, index)
		}
	}
	return out
}

func (m singleListModel) positionLabel() int {
	pos := 0
	for _, index := range m.filtered {
		if m.items[index].section {
			continue
		}
		pos++
		if m.cursor >= 0 && m.cursor < len(m.filtered) && m.filtered[m.cursor] == index {
			return pos
		}
	}
	return pos
}

// listWindowSize derives visible option rows from the terminal height.
func listWindowSize(termHeight int) int {
	if termHeight <= 0 {
		return defaultListWindow
	}
	size := termHeight - listChromeLines
	if size < minListWindow {
		return minListWindow
	}
	return size
}

// clampListWindow returns a windowStart that keeps cursor visible inside
// a window of the given size over a list of total rows.
func clampListWindow(cursor, windowStart, windowSize, total int) int {
	if total <= 0 || windowSize <= 0 {
		return 0
	}
	if windowSize > total {
		windowSize = total
	}
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= total {
		cursor = total - 1
	}
	if cursor < windowStart {
		return cursor
	}
	if cursor >= windowStart+windowSize {
		return cursor - windowSize + 1
	}
	if windowStart+windowSize > total {
		start := total - windowSize
		if start < 0 {
			return 0
		}
		return start
	}
	return windowStart
}

// clampListWindowWithSection is clampListWindow plus section-header
// recovery: when sectionStart is the filtered index of the nearest
// non-selectable group title above the cursor and that title has scrolled
// out of the top of the window, pull windowStart up to the title — but
// only when the cursor still fits in the same window (cursor-sectionStart
// < windowSize). Deep in a long group the title stays off-screen so the
// viewport can keep following the cursor.
func clampListWindowWithSection(cursor, windowStart, windowSize, total, sectionStart int) int {
	start := clampListWindow(cursor, windowStart, windowSize, total)
	if sectionStart < 0 || sectionStart >= start {
		return start
	}
	if windowSize <= 0 || cursor-sectionStart >= windowSize {
		return start
	}
	return sectionStart
}

// sectionStartForCursor returns the filtered-list index of the nearest
// section header at or above cursor, or -1 when the cursor has no group
// title above it.
func sectionStartForCursor(items []singleListItem, filtered []int, cursor int) int {
	if cursor < 0 || cursor >= len(filtered) {
		return -1
	}
	for index := cursor; index >= 0; index-- {
		itemIndex := filtered[index]
		if itemIndex < 0 || itemIndex >= len(items) {
			continue
		}
		if items[itemIndex].section {
			return index
		}
	}
	return -1
}

// firstSelectableFiltered returns the filtered-index of the first
// non-section row, or 0 when the filtered list is empty.
func firstSelectableFiltered(items []singleListItem, filtered []int) int {
	for cursor, index := range filtered {
		if !items[index].section {
			return cursor
		}
	}
	if len(filtered) == 0 {
		return 0
	}
	return 0
}

// lastSelectableFiltered returns the filtered-index of the last
// non-section row, or the last filtered index when none are selectable.
func lastSelectableFiltered(items []singleListItem, filtered []int) int {
	for cursor := len(filtered) - 1; cursor >= 0; cursor-- {
		if !items[filtered[cursor]].section {
			return cursor
		}
	}
	if len(filtered) == 0 {
		return 0
	}
	return len(filtered) - 1
}

// moveSelectableFiltered walks delta steps across selectable rows only.
// delta may be larger than 1 (page up/down). Returns -1 when there is no
// selectable row at all.
func moveSelectableFiltered(items []singleListItem, filtered []int, cursor, delta int) int {
	if len(filtered) == 0 {
		return -1
	}
	if delta == 0 {
		return cursor
	}
	step := 1
	if delta < 0 {
		step = -1
		delta = -delta
	}
	current := cursor
	moved := 0
	for moved < delta {
		next := current + step
		if next < 0 || next >= len(filtered) {
			// Clamp at the first/last selectable edge instead of wrapping,
			// which keeps long lists predictable while filtering.
			if step < 0 {
				return firstSelectableFiltered(items, filtered)
			}
			return lastSelectableFiltered(items, filtered)
		}
		current = next
		if items[filtered[current]].section {
			continue
		}
		moved++
	}
	// If we landed on a section (should not happen with continue), fix up.
	if items[filtered[current]].section {
		if step < 0 {
			return firstSelectableFiltered(items, filtered)
		}
		return lastSelectableFiltered(items, filtered)
	}
	return current
}

// runSingleListSelection launches the custom list picker on a real TTY.
func runSingleListSelection(in io.Reader, out io.Writer, title string, options []SingleOption) (string, error) {
	model := newSingleListModel(title, options)
	program := tea.NewProgram(model, tea.WithInput(in), tea.WithOutput(out))
	finalModel, err := program.Run()
	if err != nil {
		return "", fmt.Errorf("run single-select list: %w", err)
	}
	result, ok := finalModel.(singleListModel)
	if !ok {
		return "", fmt.Errorf("run single-select list: unexpected model type %T", finalModel)
	}
	if result.cancelled {
		return "", nil
	}
	return result.chosen, nil
}

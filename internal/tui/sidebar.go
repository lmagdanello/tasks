package tui

import "github.com/rivo/tview"

const shortcutsText = `[::b]SHORTCUTS[::-]

[cyan]Tab[/]     navigate
[magenta]Enter[/]  task menu
[magenta]a[/]      add task
[magenta]e[/]      edit task
[magenta]d[/]      delete task
[magenta]n[/]      new bucket
[magenta]/[/]      filter
[magenta]v[/]      archived
[magenta]p[/]      pomodoro
[magenta]q[/]      quit`

func newShortcutsPanel() *tview.TextView {
	tv := tview.NewTextView().
		SetDynamicColors(true).
		SetText(shortcutsText)
	tv.SetBorder(true).
		SetBorderColor(clrCyan).
		SetTitle(" [::b]▸ KEYS [::-] ").
		SetTitleColor(clrMagenta).
		SetBorderPadding(1, 1, 2, 2)
	return tv
}

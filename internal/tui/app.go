package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/lmagdanello/tasks/internal/config"
	"github.com/lmagdanello/tasks/internal/store"
	"github.com/rivo/tview"
)

const (
	pageMain      = "main"
	pageTaskMenu  = "taskmenu"
	pageEdit      = "edit"
	pageNewBucket = "newbucket"
)

// App represents the TUI application.
type App struct {
	app      *tview.Application
	pages    *tview.Pages
	store    *store.Store
	buckets  *tview.List
	tasks    *tview.List
	input    *tview.InputField
	filter   *tview.InputField
	status   *tview.TextView
	mainFlex *tview.Flex
	pomodoro *Pomodoro
	current  string
	showDone bool
}

// New creates and configures the TUI application with cyberpunk theme.
func New() *App {
	a := &App{app: tview.NewApplication()}
	a.setupTheme()
	a.setupLayout()
	a.setupKeybindings()
	return a
}

func (a *App) setupTheme() {
	tview.Styles.PrimitiveBackgroundColor = clrBgDark
	tview.Styles.ContrastBackgroundColor = clrBgContrast
	tview.Styles.MoreContrastBackgroundColor = clrBgMore
	tview.Styles.BorderColor = clrCyan
	tview.Styles.TitleColor = clrMagenta
	tview.Styles.GraphicsColor = clrCyan
	tview.Styles.PrimaryTextColor = clrText
	tview.Styles.SecondaryTextColor = clrCyan
	tview.Styles.TertiaryTextColor = clrMagenta
	tview.Styles.InverseTextColor = clrCyan
	tview.Styles.ContrastSecondaryTextColor = tcell.NewRGBColor(100, 100, 255)
}

func (a *App) setupLayout() {
	a.buckets = tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(clrBgContrast).
		SetSelectedTextColor(clrCyan).
		SetMainTextColor(clrTextDim).
		SetShortcutColor(clrMagenta)
	a.buckets.SetBorder(true).
		SetBorderColor(clrCyan).
		SetTitle(" [::b]▸ BUCKETS [::-] ").
		SetTitleColor(clrMagenta).
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 1)

	a.tasks = tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true).
		SetSelectedBackgroundColor(clrBgMore).
		SetSelectedTextColor(clrMagenta).
		SetMainTextColor(clrText).
		SetShortcutColor(clrCyan)
	a.tasks.SetBorder(true).
		SetBorderColor(clrCyan).
		SetTitle(" [::b]▸ TASKS [::-] ").
		SetTitleColor(clrMagenta).
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 1)

	a.input = tview.NewInputField().
		SetLabel(" [+] ").
		SetLabelColor(clrCyan).
		SetFieldTextColor(clrText).
		SetFieldBackgroundColor(clrInputBg).
		SetPlaceholderTextColor(clrPlaceholder).
		SetPlaceholder("nova task...")
	a.input.SetBorder(true).SetBorderColor(clrCyan).SetBorderPadding(0, 0, 1, 1)

	a.filter = tview.NewInputField().
		SetLabel(" [/] ").
		SetLabelColor(clrCyan).
		SetFieldTextColor(clrText).
		SetFieldBackgroundColor(clrInputBg).
		SetPlaceholderTextColor(clrPlaceholder).
		SetPlaceholder("filtrar...")
	a.filter.SetBorder(true).SetBorderColor(clrCyan).SetBorderPadding(0, 0, 1, 1)
	a.filter.SetChangedFunc(func(string) { a.loadTasks() })
	a.filter.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			a.app.SetFocus(a.tasks)
		}
	})

	a.status = tview.NewTextView().
		SetDynamicColors(true).
		SetText("[cyan]Tab[/] nav  [magenta]Enter[/] menu  [magenta]a[/] add  [magenta]e[/] edit  [magenta]d[/] del  [magenta]n[/] bucket  [magenta]v[/] archived  [magenta]p[/] pomodoro  [magenta]q[/] quit")
	a.status.SetBackgroundColor(clrStatusBg)

	contentFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.tasks, 0, 1, false).
		AddItem(a.filter, 3, 0, false).
		AddItem(a.input, 3, 0, false)

	a.pomodoro = newPomodoro(a.app)
	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(newShortcutsPanel(), 0, 1, false).
		AddItem(a.pomodoro.View(), 10, 0, false)

	a.mainFlex = tview.NewFlex().
		AddItem(a.buckets, 22, 0, true).
		AddItem(contentFlex, 0, 1, false).
		AddItem(rightPanel, 24, 0, false)

	a.pages = tview.NewPages()
	a.pages.AddPage(pageMain, a.mainFlex, true, true)

	a.buckets.SetChangedFunc(func(int, string, string, rune) {})
	a.buckets.SetSelectedFunc(a.onBucketChosen)
	a.input.SetDoneFunc(a.onInputDone)
	a.tasks.SetSelectedFunc(a.onTaskSelected)
}

func (a *App) focusNext() {
	switch a.app.GetFocus() {
	case a.buckets:
		a.app.SetFocus(a.tasks)
	case a.tasks:
		a.app.SetFocus(a.filter)
	case a.filter:
		a.app.SetFocus(a.input)
	case a.input:
		a.app.SetFocus(a.buckets)
	default:
		a.app.SetFocus(a.buckets)
	}
}

func (a *App) focusPrev() {
	switch a.app.GetFocus() {
	case a.buckets:
		a.app.SetFocus(a.input)
	case a.tasks:
		a.app.SetFocus(a.buckets)
	case a.filter:
		a.app.SetFocus(a.tasks)
	case a.input:
		a.app.SetFocus(a.filter)
	default:
		a.app.SetFocus(a.buckets)
	}
}

func (a *App) setupKeybindings() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		focus := a.app.GetFocus()
		front, _ := a.pages.GetFrontPage()
		inInputField := focus == a.input || focus == a.filter
		inModal := front != pageMain

		if event.Key() == tcell.KeyTab && front == pageMain {
			if event.Modifiers()&tcell.ModShift != 0 {
				a.focusPrev()
			} else {
				a.focusNext()
			}
			return nil
		}
		if event.Key() == tcell.KeyEnter && focus == a.tasks {
			a.showTaskMenu()
			return nil
		}
		if inInputField || inModal {
			return event
		}
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q', 'Q':
				a.app.Stop()
				return nil
			case 'a', 'A':
				a.app.SetFocus(a.input)
				return nil
			case 'n', 'N':
				a.showNewBucketModal()
				return nil
			case 'e', 'E':
				a.editSelectedTask()
				return nil
			case 'd', 'D':
				a.deleteSelectedTask()
				return nil
			case '/':
				a.app.SetFocus(a.filter)
				return nil
			case 'v', 'V':
				a.toggleShowDone()
				return nil
			case 'p', 'P':
				a.pomodoro.Start()
				return nil
			}
		}
		return event
	})
}

func (a *App) toggleShowDone() {
	a.showDone = !a.showDone
	if a.current != "" {
		title := " [::b]▸ TASKS [::-] "
		if a.showDone {
			title = " [::b]▸ ARQUIVADAS [::-] "
		}
		a.tasks.SetTitle(fmt.Sprintf("%s [::i]%s[::-]", title, a.current))
	}
	a.loadTasks()
}

func (a *App) onBucketChosen(index int, main, secondary string, shortcut rune) {
	if index == a.buckets.GetItemCount()-1 {
		a.showNewBucketModal()
		return
	}
	bucket := strings.TrimSpace(main)
	if bucket == "" {
		return
	}
	a.loadBucket(bucket)
	a.app.SetFocus(a.tasks)
}

func (a *App) onInputDone(key tcell.Key) {
	if key != tcell.KeyEnter {
		if key == tcell.KeyEscape {
			a.app.SetFocus(a.tasks)
		}
		return
	}

	text := strings.TrimSpace(a.input.GetText())
	if text == "" {
		a.app.SetFocus(a.tasks)
		return
	}

	if a.store == nil {
		a.setStatus("[red]Selecione um bucket primeiro")
		a.app.SetFocus(a.buckets)
		return
	}

	id, err := a.store.CreateTask(text)
	if err != nil {
		a.setStatus("[red]Erro: " + err.Error())
		return
	}

	a.input.SetText("")
	a.loadTasks()
	a.setStatus(fmt.Sprintf("[green]Task #%d adicionada", id))
	a.app.SetFocus(a.tasks)
}

func (a *App) onTaskSelected(index int, main, secondary string, shortcut rune) {
	a.showTaskMenu()
}

func (a *App) loadBuckets() {
	a.buckets.Clear()
	buckets, err := config.ListBuckets()
	if err != nil {
		a.buckets.AddItem(" [red]Erro ao listar", "", 0, nil)
		return
	}

	for _, b := range buckets {
		a.buckets.AddItem(" "+b, "", 0, nil)
	}
	a.buckets.AddItem(" [dim]+ Novo bucket", "", 0, nil)

	if len(buckets) > 0 && a.current == "" {
		a.current = buckets[0]
		a.loadBucket(a.current)
	}
}

func (a *App) loadBucket(bucket string) {
	if a.store != nil {
		a.store.Close()
	}
	a.current = bucket
	s, err := store.Open(bucket)
	if err != nil {
		a.setStatus("[red]" + err.Error())
		a.store = nil
		a.tasks.Clear()
		return
	}
	a.store = s
	a.loadTasks()
	a.tasks.SetTitle(fmt.Sprintf(" [::b]▸ TASKS [::-] [::i]%s[::-]", bucket))
}

func (a *App) loadTasks() {
	a.tasks.Clear()
	if a.store == nil {
		return
	}
	opts := a.buildListOpts()
	var tasks []store.Task
	var err error
	if a.showDone {
		tasks, err = a.store.DoneTasks(opts...)
	} else {
		tasks, err = a.store.AllTasks(opts...)
	}
	if err != nil {
		a.tasks.AddItem(" [red]Erro: "+err.Error(), "", 0, nil)
		return
	}
	for i := range tasks {
		a.tasks.AddItem(a.formatTaskLine(i+1, &tasks[i]), "", 0, nil)
	}
	if len(tasks) == 0 {
		msg := " [dim]Nenhuma task..."
		if a.showDone {
			msg = " [dim]Nenhuma arquivada..."
		}
		a.tasks.AddItem(msg, "", 0, nil)
	}
}

func (a *App) buildListOpts() []store.ListOption {
	opts := make([]store.ListOption, 0, 1)
	if a.filter != nil {
		q := strings.TrimSpace(a.filter.GetText())
		if q != "" {
			opts = append(opts, store.FilterQuery(q))
		}
	}
	return opts
}

func (a *App) formatTaskLine(idx int, t *store.Task) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(" %d. %s", idx, t.Text))
	if t.Priority != store.PriorityMedium {
		b.WriteString(fmt.Sprintf(" [%s]", t.Priority))
	}
	if t.Due != nil {
		b.WriteString(fmt.Sprintf(" (%s)", t.Due.Format("01/02")))
	}
	return b.String()
}

func (a *App) getSelectedTask() (*store.Task, bool) {
	if a.store == nil {
		return nil, false
	}
	idx := a.tasks.GetCurrentItem()
	var tasks []store.Task
	var err error
	if a.showDone {
		tasks, err = a.store.DoneTasks()
	} else {
		tasks, err = a.store.AllTasks()
	}
	if err != nil || idx >= len(tasks) {
		return nil, false
	}
	return &tasks[idx], true
}

func (a *App) showTaskMenu() {
	task, ok := a.getSelectedTask()
	if !ok {
		return
	}

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Task: %s", task.Text)).
		AddButtons([]string{"Editar", "Arquivar", "Excluir", "Cancelar"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			a.pages.RemovePage(pageTaskMenu)
			a.app.SetFocus(a.tasks)
			if buttonIndex < 0 {
				return
			}
			switch buttonLabel {
			case "Editar":
				a.showEditTaskModal(task)
			case "Arquivar":
				if err := a.store.ArchiveTask(task.Key); err != nil {
					a.setStatus("[red]" + err.Error())
				} else {
					a.loadTasks()
					a.setStatus(fmt.Sprintf("[green]✓ \"%s\" arquivada", task.Text))
				}
			case "Excluir":
				if err := a.store.DeleteTask(task.Key); err != nil {
					a.setStatus("[red]" + err.Error())
				} else {
					a.loadTasks()
					a.setStatus("[green]Task excluída")
				}
			}
		})
	modal.SetBorder(true).
		SetTitle(" [::b]▸ AÇÕES [::-] ").
		SetTitleColor(clrMagenta).
		SetBorderColor(clrCyan)

	a.pages.AddPage(pageTaskMenu, modal, true, true)
	a.app.SetFocus(modal)
}

func (a *App) deleteSelectedTask() {
	task, ok := a.getSelectedTask()
	if !ok {
		return
	}
	if err := a.store.DeleteTask(task.Key); err != nil {
		a.setStatus("[red]" + err.Error())
		return
	}
	a.loadTasks()
	a.setStatus("[green]Task excluída")
}

func (a *App) editSelectedTask() {
	task, ok := a.getSelectedTask()
	if !ok {
		return
	}
	a.showEditTaskModal(task)
}

func (a *App) showEditTaskModal(task *store.Task) {
	form := tview.NewForm().AddInputField("Task:", task.Text, 40, nil, nil)
	editInput := form.GetFormItem(0).(*tview.InputField)
	form.AddButton("Salvar", func() {
		text := strings.TrimSpace(editInput.GetText())
		if text != "" {
			if err := a.store.UpdateTask(task.Key, text); err != nil {
				a.setStatus("[red]" + err.Error())
			} else {
				a.loadTasks()
				a.setStatus("[green]Task atualizada")
			}
		}
		a.pages.RemovePage(pageEdit)
		a.app.SetFocus(a.mainFlex)
	})
	form.AddButton("Cancelar", func() {
		a.pages.RemovePage(pageEdit)
		a.app.SetFocus(a.mainFlex)
	})
	form.SetBorder(true).
		SetTitle(" [::b]EDITAR TASK [::-] ").
		SetTitleColor(clrMagenta).
		SetBorderColor(clrCyan)
	a.pages.AddPage(pageEdit, form, true, true)
	a.app.SetFocus(form)
}

func (a *App) setStatus(msg string) {
	if a.status != nil {
		a.status.SetText(msg)
	}
}

func (a *App) showNewBucketModal() {
	form := tview.NewForm().AddInputField("Bucket:", "", 20, nil, nil)
	inputField := form.GetFormItem(0).(*tview.InputField)
	form.AddButton("Criar", func() {
		bucket := strings.TrimSpace(inputField.GetText())
		if bucket == "" {
			return
		}
		if err := store.CreateBucket(bucket); err != nil {
			if err == store.ErrBucketExists {
				a.setStatus("[red]Bucket já existe")
			} else {
				a.setStatus("[red]" + err.Error())
			}
		} else {
			a.loadBuckets()
			a.setStatus(fmt.Sprintf("[green]Bucket \"%s\" criado", bucket))
		}
		a.pages.RemovePage(pageNewBucket)
		a.app.SetFocus(a.mainFlex)
	})
	form.AddButton("Cancelar", func() {
		a.pages.RemovePage(pageNewBucket)
		a.app.SetFocus(a.mainFlex)
	})
	form.SetBorder(true).
		SetTitle(" [::b]NOVO BUCKET [::-] ").
		SetTitleColor(clrMagenta).
		SetBorderColor(clrCyan)
	frame := tview.NewFrame(form).SetBorders(0, 0, 0, 0, 0, 0)
	a.pages.AddPage(pageNewBucket, frame, true, true)
	a.app.SetFocus(form)
}

// Run starts the application.
func (a *App) Run() error {
	a.loadBuckets()
	return a.app.SetRoot(a.pages, true).EnableMouse(false).Run()
}

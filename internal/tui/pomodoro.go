package tui

import (
	"fmt"
	"sync"
	"time"

	"github.com/rivo/tview"
)

const (
	pomodoroWorkMin  = 25
	pomodoroBreakMin = 5
)

// Pomodoro manages the Pomodoro timer state and display.
type Pomodoro struct {
	mu       sync.Mutex
	app      *tview.Application
	view     *tview.TextView
	running  bool
	workTime time.Duration
	remain   time.Duration
	stopCh   chan struct{}
}

func newPomodoro(app *tview.Application) *Pomodoro {
	p := &Pomodoro{
		app:      app,
		workTime: pomodoroWorkMin * time.Minute,
		remain:   pomodoroWorkMin * time.Minute,
		stopCh:   make(chan struct{}),
	}
	p.view = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	p.view.SetBorder(true).
		SetBorderColor(clrCyan).
		SetTitle(" [::b]▸ POMODORO [::-] ").
		SetTitleColor(clrMagenta).
		SetBorderPadding(1, 1, 2, 2)
	p.render()
	return p
}

func (p *Pomodoro) render() {
	p.mu.Lock()
	defer p.mu.Unlock()

	m := int(p.remain.Minutes())
	s := int(p.remain.Seconds()) % 60
	status := " [dim]idle[-]"
	if p.running {
		status = " [green]● running[-]"
	}
	text := fmt.Sprintf("\n[cyan]%02d:%02d[-]%s\n\n[dim]p[/] start/stop", m, s, status)
	if p.view != nil {
		p.view.SetText(text)
	}
}

// Start toggles the timer. If running, stops. If idle, starts countdown.
func (p *Pomodoro) Start() {
	p.mu.Lock()
	if p.running {
		p.running = false
		select {
		case p.stopCh <- struct{}{}:
		default:
		}
		p.mu.Unlock()
		go p.scheduleUpdate()
		return
	}
	p.running = true
	p.mu.Unlock()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		p.scheduleUpdate()

		for {
			select {
			case <-p.stopCh:
				return
			case <-ticker.C:
				p.mu.Lock()
				p.remain -= time.Second
				rem := p.remain
				running := p.running
				p.mu.Unlock()

				p.scheduleUpdate()

				if rem <= 0 && running {
					p.mu.Lock()
					p.running = false
					p.remain = pomodoroWorkMin * time.Minute
					p.mu.Unlock()
					return
				}
			}
		}
	}()
}

// scheduleUpdate queues a UI update from a goroutine to avoid deadlock
// (QueueUpdateDraw blocks when called from the main event handler).
func (p *Pomodoro) scheduleUpdate() {
	p.app.QueueUpdateDraw(p.render)
}

// View returns the TextView for layout.
func (p *Pomodoro) View() *tview.TextView {
	return p.view
}

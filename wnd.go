package main

import (
	"context"
	"strings"
	"sync/atomic"
	"time"

	"github.com/lxn/walk"
)

type Translator struct {
	goBtn     *walk.PushButton
	inText    *walk.TextEdit
	outText   *walk.TextEdit
	wnd       *walk.MainWindow
	cfg       *Config
	inWorking atomic.Bool
	process   *translatorProcess
}

type translatorProcess struct {
	ch     chan *translateMessage
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTranslator(cfg *Config) *Translator {
	println(cfg.String())
	return &Translator{
		outText: nil,
		cfg:     cfg,
	}
}

func (t *Translator) CopyToClipboard() {
	if t.inWorking.Load() {
		return
	}
	text := t.outText.Text()
	if text == "" {
		return
	}

	clipboard := walk.Clipboard()
	err := clipboard.SetText(text)
	if err != nil {
		return
	}

	go func() {
		origin := t.goBtn.Text()
		t.goBtn.SetText("Copied to clipboard")
		time.Sleep(2 * time.Second)
		t.goBtn.SetText(origin)
	}()
}

func (t *Translator) Translate(llm LLModel, systemPrompt string, text string) {
	text = strings.TrimSpace(text)
	if t.inWorking.Load() {
		t.Stop()
		return
	}

	t.inWorking.Store(true)
	t.goBtn.SetText("Stop")
	t.outText.SetText("")

	ctx, cancel := context.WithCancel(context.Background())
	t.process = &translatorProcess{
		ch:     make(chan *translateMessage),
		ctx:    ctx,
		cancel: cancel,
	}
	defer close(t.process.ch)
	go func() {
		var e = llm.GenerateContent(ctx, systemPrompt, text, t.process.ch)
		if e != nil {
			walk.MsgBox(t.wnd, "Error", "Failed to generate content: "+e.Error(), walk.MsgBoxIconError)
		}
	}()

	var finished = false
outerLoop:
	for !finished {
		select {
		case <-t.process.ctx.Done():
			return
		case w, ok := <-t.process.ch:
			if !ok {
				return
			}
			if w.cmd == finishedCmd {
				finished = true
				break outerLoop
			}
			if w.cmd == errorCmd {
				finished = true
				break outerLoop
			}
			t.outText.AppendText(w.text)
		}
	}

	println("come here")
	t.inWorking.Store(false)
	t.goBtn.SetText("Go Translate")
}

func (t *Translator) Stop() {
	if !t.inWorking.Load() {
		panic("Stop called when not working")
	}
	if t.process != nil {
		t.process.cancel()
	}

	t.inWorking.Store(false)
	t.goBtn.SetText("Go Translate")
}

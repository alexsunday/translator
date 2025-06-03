package main

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	// rsrc embedded the logo icon to resource id 2.
	ico, err := walk.NewIconFromResourceIdWithSize(2, walk.Size{Width: 64, Height: 64})
	if err != nil {
		walk.MsgBox(nil, "Error", "Failed to load icon: "+err.Error(), walk.MsgBoxIconError)
		return
	}

	cfg, err := loadConfig()
	if err != nil {
		walk.MsgBox(nil, "Error", "Failed to load configuration: "+err.Error(), walk.MsgBoxIconError)
		return
	}

	// llm := NewOpenAiLLM(cfg.Base, cfg.Model, cfg.Key)
	llm, err := NewLangChainLLM(cfg)
	if err != nil {
		walk.MsgBox(nil, "Error", "Failed to create LLM: "+err.Error(), walk.MsgBoxIconError)
		return
	}

	var translator = NewTranslator(cfg)

	MainWindow{
		Title:    "Translator",
		AssignTo: &translator.wnd,
		Icon:     ico,
		Size: Size{
			Width:  600,
			Height: 300,
		},
		Layout: VBox{
			MarginsZero: true,
			SpacingZero: true,
		},
		Children: []Widget{
			VSplitter{
				Children: []Widget{
					TextEdit{AssignTo: &translator.inText},
					TextEdit{AssignTo: &translator.outText, ReadOnly: true},
				},
			},
			PushButton{
				Text:     "Go Translate",
				AssignTo: &translator.goBtn,
				OnClicked: func() {
					fmt.Printf("translator: %p\n", &translator)
					go translator.Translate(llm, cfg.System, translator.inText.Text())
				},
			},
		},
	}.Run()
}

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

func (t *Translator) Translate(llm LLModel, systemPrompt string, text string) {
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

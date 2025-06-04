package main

import (
	"fmt"

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
	go setSysTray(translator)

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
					TextEdit{
						AssignTo: &translator.outText,
						ReadOnly: true,
						OnMouseDown: func(x, y int, button walk.MouseButton) {
							if button == walk.LeftButton {
								translator.CopyToClipboard()
							}
						},
					},
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

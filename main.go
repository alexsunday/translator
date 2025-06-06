package main

import (
	"flag"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var (
	confName      = flag.String("c", "", "Path to the configuration file")
	builtinConfig = ``
)

/*
考虑将 api key 配置内置的方案
1. 先看看命令行CLI是否指定了配置文件 ，若指定了，则使用命令行的；
2. 看看是否有内置的 若有则使用内置
3. 都没有，则查找 conf 文件
*/

func indexConfig() (*Config, error) {
	flag.Parse()

	if *confName != "" {
		return loadConfig(*confName)
	}

	content := strings.TrimSpace(builtinConfig)
	if content != "" {
		return loadConfigViaFileContent(content)
	}

	return loadConfig("conf")
}

func main() {
	cfg, err := indexConfig()
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

	title := "Translator"
	if cfg.Dict["title"] != "" {
		title = cfg.Dict["title"]
	}

	// rsrc embedded the logo icon to resource id 2.
	ico, err := walk.NewIconFromResourceIdWithSize(2, walk.Size{Width: 64, Height: 64})
	if err != nil {
		walk.MsgBox(nil, "Error", "Failed to load icon: "+err.Error(), walk.MsgBoxIconError)
		return
	}

	err = MainWindow{
		Title:    title,
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
					go translator.Translate(llm, cfg.System, translator.inText.Text())
				},
			},
		},
	}.Create()
	if err != nil {
		walk.MsgBox(nil, "Error", "Failed to create main window: "+err.Error(), walk.MsgBoxIconError)
		return
	}

	translator.inText.KeyDown().Attach(func(key walk.Key) {
		if key == walk.KeyReturn {
			keyMod := walk.ModifiersDown()
			if keyMod&walk.ModControl != 0 {
				if translator.goBtn != nil {
					if translator.inWorking.Load() {
						return
					}
					go translator.Translate(llm, cfg.System, translator.inText.Text())
				}
			}
		}
	})

	translator.wnd.Show()
	translator.wnd.Run()
}

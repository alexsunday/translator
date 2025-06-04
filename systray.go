package main

import (
	_ "embed"

	"github.com/getlantern/systray"
	"github.com/lxn/win"
)

//go:embed tray.ico
var trayPngPic []byte

var topMost = false

func setSysTray(app *Translator) {
	systray.Run(func() {
		systray.SetIcon(trayPngPic)
		systray.SetTitle("Translator")
		systray.SetTooltip("Translator")

		mTop := systray.AddMenuItemCheckbox("总在最上层", "", topMost)
		mQuit := systray.AddMenuItem("退出", "")

		for {
			select {
			case <-mTop.ClickedCh:
				if app.wnd == nil {
					return
				}
				var hwndInsertAfter win.HWND
				topMost = !topMost
				if topMost {
					hwndInsertAfter = win.HWND_TOPMOST
					mTop.Check()
				} else {
					hwndInsertAfter = win.HWND_NOTOPMOST
					mTop.Uncheck()
				}
				win.SetWindowPos(app.wnd.Handle(), hwndInsertAfter, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)
			case <-mQuit.ClickedCh:
				systray.Quit()
				if app.wnd != nil {
					app.wnd.Close()
				}
				return
			}
		}
	}, nil)
}

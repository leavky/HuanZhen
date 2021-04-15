package main

import (
	"HuanZhen/config"
	"HuanZhen/core/dnsProxy"
	"HuanZhen/core/pcapCheck"
	"HuanZhen/core/portConnCheck"
	"HuanZhen/core/portForwarding"
	"HuanZhen/core/processCheck"
	"HuanZhen/pages"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var defaultConfig config.Config

// 读取配置文件
func init() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("读取配置文件失败")
	}
	err = json.Unmarshal(data, &defaultConfig)
	if err != nil {
		panic("解析配置文件失败")
	}
	//fmt.Println(defaultConfig)

	// 检车配置是否正确，主要检测端口是否冲突，要求每个服务的端口不冲突

	// 启动端口转发服务
	go portForwarding.StartPortForward(defaultConfig)

	// 启动 DNS 代理服务
	go dnsProxy.StartDnsProxy()

	// 启动端口连接检测服务
	go portConnCheck.StartPortConnCheck(defaultConfig.PortConnCheck)

	// 启动进程检测
	go processCheck.StartProcessCheck()

	// 启动数据包检测
	go pcapCheck.StartCheckPcap()

}

var topWindow fyne.Window
const preferenceCurrentTutorial = "currentTutorial"

func main() {
	log.Println("@===============⚡幻阵⚡====================@\n ===========虚实之间，防护溯源。=================")

	//a := app.NewWithID("huanzhen")
	//a.Settings().SetTheme(&myThemes.MyTheme{})  // 不自定义中文会乱码
	//w := a.NewWindow("幻阵")
	//topWindow = w
	//
	//content := container.NewMax()
	//title := widget.NewLabel("Component name")
	////intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	////intro.Wrapping = fyne.TextWrapWord
	//
	//setTutorial := func(page pages.Page) {
	//	if fyne.CurrentDevice().IsMobile() {
	//		child := a.NewWindow(page.Title)
	//		topWindow = child
	//		child.SetContent(page.View(topWindow))
	//		child.Show()
	//		child.SetOnClosed(func() {
	//			topWindow = w
	//		})
	//		return
	//	}
	//
	//	title.SetText(page.Title)
	//	//intro.SetText(page.Intro)
	//
	//	content.Objects = []fyne.CanvasObject{page.View(w)}
	//	content.Refresh()
	//}
	//
	//tutorial := container.NewBorder(container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, content)
	//if fyne.CurrentDevice().IsMobile() {
	//	w.SetContent(makeNav(setTutorial, false))
	//} else {
	//	split := container.NewHSplit(makeNav(setTutorial, true), tutorial)
	//	split.Offset = 0.2
	//	w.SetContent(split)
	//}



	//hello := widget.NewLabel("Hello Fyne!")
	//w.SetContent(container.NewVBox(
	//	hello,
	//	widget.NewButton("Hi!", func() {
	//		hello.SetText("Welcome :)")
	//	}),
	//))
	//w.Resize(fyne.NewSize(640, 460))
	//w.ShowAndRun()



	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)
	<-sig
	log.Println("Bye")


}

func makeNav(setTutorial func(page pages.Page), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return pages.PageIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := pages.PageIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := pages.Pages[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := pages.Pages[uid]; ok {
				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
		tree.Select(currentPref)
	}

	//mytheme := myThemes.MyTheme{}
	//themes := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
	//	widget.NewButton("Dark", func() {
	//		a.Settings().SetTheme(theme.DarkTheme())
	//	}),
	//	widget.NewButton("Light", func() {
	//		a.Settings().SetTheme(theme.LightTheme())
	//	}),
	//)

	return container.NewBorder(nil, nil, nil, nil, tree)
}

package main

import (
	"os"
	"os/signal"

	"github.com/nepail/go-sync/server"
	"github.com/zserge/lorca"
	// "github.com/nepail/lorca"
)

func main() {
	// chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	// cmd := exec.Command(chromePath, "--app=https://www.google.com")
	// cmd.Start()

	// 啟動 gin 服務
	go server.Run()

	ui := StartBrowser()

	// 監聽中斷訊號
	chSignal := listenToInterrupt()

	// select 監聽可讀可寫的chan 沒有值的話阻塞目前線程 隨機輪詢
	// 等待中斷訊號
	select {
	case <-ui.Done():
	case <-chSignal:
		// cmd.Process.Kill()
	}

	// select {}

	// 關閉主線程或ui，會主動退出
	ui.Close()

	// 等待命令
	// <-chSignal
	// cmd.Process.Kill()
}

func StartBrowser() lorca.UI {
	var ui lorca.UI
	// ui, _ = lorca.New("http://google.com", "", 800, 600, "--disable-sync", "--disable-translate")
	// ui, _ = lorca.New("https://term.ptt.cc", "", 800, 600, "--disable-sync", "--disable-translate")
	// localhost 不走代理，所以必須使用127.0.0.1
	ui, _ = lorca.New("http://localhost:27149/static/index.html", "", 800, 600, "--disable-sync", "--disable-translate")
	return ui
}

func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}

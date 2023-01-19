package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/nepail/go-sync/config"
	"github.com/nepail/go-sync/server"
	"github.com/zserge/lorca"
	// "github.com/nepail/lorca"
)

func main() {
	// chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	// cmd := exec.Command(chromePath, "--app=https://www.google.com")
	// cmd.Start()

	chChromeDie := make(chan struct{})
	chBackendDie := make(chan struct{})

	// 啟動 gin 服務
	go server.Run()
	go StartBrowser(chChromeDie, chBackendDie)

	// ui := StartBrowser()

	// 監聽中斷訊號
	chSignal := listenToInterrupt()

	// 等待中斷訊號
	for {
		select {

		case <-chSignal:
			fmt.Println("--------cmd close -----------")
			chBackendDie <- struct{}{}

		case <-chChromeDie:
			fmt.Println("--------chrome close --------")
			os.Exit(0)

		}
	}

	// 關閉主線程或ui，會主動退出
	// ui.Close()

	// 等待命令
	// <-chSignal
	// cmd.Process.Kill()
}

// func StartBrowser() lorca.UI {
// 	var ui lorca.UI
// 	// ui, _ = lorca.New("http://google.com", "", 800, 600, "--disable-sync", "--disable-translate")
// 	// ui, _ = lorca.New("https://term.ptt.cc", "", 800, 600, "--disable-sync", "--disable-translate")
// 	// localhost 不走代理，所以必須使用127.0.0.1
// 	ui, _ = lorca.New("http://localhost:"+config.GetPort()+"/static/index.html", "", 800, 600, "--disable-sync", "--disable-translate")
// 	return ui
// }

func StartBrowser(chChromeDie chan struct{}, chBackendDie chan struct{}) {

	// ui, _ = lorca.New("http://google.com", "", 800, 600, "--disable-sync", "--disable-translate")
	// ui, _ = lorca.New("https://term.ptt.cc", "", 800, 600, "--disable-sync", "--disable-translate")
	// localhost 不走代理，所以必須使用127.0.0.1
	// ui, _ = lorca.New("http://localhost:"+config.GetPort()+"/static/index.html", "", 800, 600, "--disable-sync", "--disable-translate")

	var ui lorca.UI
	ui, _ = lorca.New("http://localhost:"+config.GetPort()+"/static/index.html", "", 800, 600, "--disable-sync", "--disable-translate")

	for {
		select {
		case <-ui.Done():
			chChromeDie <- struct{}{}
		case <-chBackendDie:
			ui.Close()
		}
	}

}

// 監聽中斷訊息
func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}

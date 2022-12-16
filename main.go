package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/zserge/lorca"
)

func main() {
	var ui lorca.UI
	ui, _ = lorca.New("http://google.com", "", 800, 600, "--disable-sync", "--disable-translate")
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)

	// select 監聽可讀可寫的chan 沒有值的話阻塞目前線程 隨機輪詢
	select {
	case <-ui.Done():
	case <-chSignal:
	}

	// 關閉主線程或ui，會主動退出
	ui.Close()

}

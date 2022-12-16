package main

import (
	"net/http"
	"os"
	"os/signal"

	// "os/exec"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"
)

func main() {
	// chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	// cmd := exec.Command(chromePath, "--app=https://www.google.com")
	// cmd.Start()

	go func() {
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		router.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "<h1> Hello World </h1>")
		})
		router.Run(":8080")
	}()

	var ui lorca.UI
	// ui, _ = lorca.New("http://google.com", "", 800, 600, "--disable-sync", "--disable-translate")
	// localhost 不走代理，所以必須使用127.0.0.1
	ui, _ = lorca.New("http://127.0.0.1:8080", "", 800, 600, "--disable-sync", "--disable-translate")
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)

	// select 監聽可讀可寫的chan 沒有值的話阻塞目前線程 隨機輪詢
	select {
	case <-ui.Done():
	case <-chSignal:
	}

	// select {}

	// 關閉主線程或ui，會主動退出
	ui.Close()

}

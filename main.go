package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	// "os/exec"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"
	// "github.com/nepail/lorca"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	// chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	// cmd := exec.Command(chromePath, "--app=https://www.google.com")
	// cmd.Start()

	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		// router.GET("/", func(c *gin.Context) {
		// 	c.String(http.StatusOK, "<h1> Hello World </h1>")
		// })
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.StaticFS("/static", http.FS(staticFiles))
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/static/") {
				reader, err := staticFiles.Open("index.html")
				if err != nil {
					log.Fatal(err)
				}
				defer reader.Close()
				stat, err := reader.Stat()
				if err != nil {
					log.Fatal(err)
				}
				c.DataFromReader(http.StatusOK, stat.Size(), "text/html;charset=utf-8", reader, nil)
			} else {
				// c.Status(http.StatusNotFound)
				c.Status(404)
			}
		})
		router.Run(":8080")
	}()

	var ui lorca.UI
	// ui, _ = lorca.New("http://google.com", "", 800, 600, "--disable-sync", "--disable-translate")
	// localhost 不走代理，所以必須使用127.0.0.1
	ui, _ = lorca.New("http://localhost:8080/static/index.html", "", 800, 600, "--disable-sync", "--disable-translate")
	// ui, _ = lorca.New("https://term.ptt.cc", "", 800, 600, "--disable-sync", "--disable-translate")
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// select 監聽可讀可寫的chan 沒有值的話阻塞目前線程 隨機輪詢
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

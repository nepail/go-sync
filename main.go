package main

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"

	// "os/exec"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"

	// "github.com/nepail/lorca"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

// 把指定目錄的檔案嵌入到exe
//
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
		router.GET("/api/v1/qrcodes", QrcodeController)
		router.GET("/uploads/:path", UploadsController)
		router.GET("/api/v1/addresses", AddressesController)
		router.POST("/api/v1/texts", TextsController)
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

func TextsController(c *gin.Context) {

	var json struct {
		Raw string `json:"raw"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		exe, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exe)
		if err != nil {
			log.Fatal(err)
		}
		filename := uuid.New().String()
		uploads := filepath.Join(dir, "uploads")
		err = os.MkdirAll(uploads, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		fullpath := path.Join("uploads", filename+".txt")
		err = ioutil.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})
	}
}

func AddressesController(c *gin.Context) {
	addrs, _ := net.InterfaceAddrs()
	var result []string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"addresses": result})
}

func GetUploadsDir() (uploads string) {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe)
	uploads = filepath.Join(dir, "uploads")
	return
}

func UploadsController(c *gin.Context) {
	if path := c.Param("path"); path != "" {
		target := filepath.Join(GetUploadsDir(), path)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename"+path)
		c.Header("Content-Type", "application/octet-stream")
		c.File(target)
	} else {
		c.Status(http.StatusNotFound)
	}
}

func QrcodeController(c *gin.Context) {
	if content := c.Query("content"); content != "" {
		png, err := qrcode.Encode(content, qrcode.Medium, 256)
		if err != nil {
			log.Fatal(err)
		}
		c.Data(http.StatusOK, "imge/png", png)
	} else {
		c.Status(http.StatusBadRequest)
	}
}

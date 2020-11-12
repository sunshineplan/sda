package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/httpsvr"
	"github.com/vharitonsky/iniflags"
)

var logPath string
var server httpsvr.Server

func main() {
	self, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&server.Unix, "unix", "", "UNIX-domain Socket")
	flag.StringVar(&server.Host, "host", "0.0.0.0", "Server Host")
	flag.StringVar(&server.Port, "port", "12345", "Server Port")
	//flag.StringVar(&logPath, "log", "/var/log/app/sda-go.log", "Log Path")
	flag.StringVar(&logPath, "log", "", "Log Path")
	iniflags.SetConfigFile(filepath.Join(filepath.Dir(self), "config.ini"))
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.Parse()

	if logPath != "" {
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatalln("Failed to open log file:", err)
		}
		gin.DefaultWriter = f
		gin.DefaultErrorWriter = f
		log.SetOutput(f)
	}

	router := gin.Default()
	server.Handler = router
	router.StaticFS("/static", http.Dir(filepath.Join(filepath.Dir(self), "static")))
	router.LoadHTMLGlob(filepath.Join(filepath.Dir(self), "templates/*"))

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	router.POST("/analyze", analyze)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

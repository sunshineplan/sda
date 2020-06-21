package main

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/vharitonsky/iniflags"
)

func main() {
	self, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	unix := flag.String("unix", "", "Server Host")
	host := flag.String("host", "127.0.0.1", "Server Host")
	port := flag.String("port", "12345", "Server Port")
	logPath := flag.String("log", "/var/log/app/sda-go.log", "Log Path")
	iniflags.SetConfigFile(filepath.Join(filepath.Dir(self), "config.ini"))
	iniflags.SetAllowMissingConfigFile(true)
	iniflags.Parse()

	f, _ := os.OpenFile(*logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	gin.DefaultWriter = io.MultiWriter(f)

	router := gin.Default()
	router.StaticFS("/static", http.Dir(filepath.Join(filepath.Dir(self), "static")))
	router.LoadHTMLGlob(filepath.Join(filepath.Dir(self), "templates/*"))

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.POST("/analysis", handler)

	if *unix != "" && runtime.GOOS == "linux" {
		if _, err = os.Stat(*unix); err == nil {
			err = os.Remove(*unix)
			if err != nil {
				log.Fatal(err)
			}
		}

		listener, err := net.Listen("unix", *unix)
		if err != nil {
			log.Fatal(err)
		}

		idleConnsClosed := make(chan struct{})
		go func() {
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit

			if err := listener.Close(); err != nil {
				log.Printf("HTTP Listener close: %v", err)
			}
			if err := os.Remove(*unix); err != nil {
				log.Printf("Remove socket file: %v", err)
			}
			close(idleConnsClosed)
		}()

		if err = os.Chmod(*unix, 0666); err != nil {
			log.Fatal(err)
		}

		http.Serve(listener, router)
		<-idleConnsClosed
	} else {
		router.Run(*host + ":" + *port)
	}
}

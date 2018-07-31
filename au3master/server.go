package au3master

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server is a http server that relays commands to an
// autoit executable.
type Server struct {
	r         *gin.Engine
	tosend    chan *Command
	toreceive map[string]chan *Result
	hostaddr  string
	shutdown0 chan bool
	shutdown1 chan bool
	htps      *http.Server
}

// New autoit command server relay
func New(path string, devmode bool) *Server {
	if !devmode {
		gin.SetMode(gin.ReleaseMode)
	}
	s := &Server{
		tosend:    make(chan *Command, 32),
		toreceive: make(map[string]chan *Result),
		hostaddr:  path,
		shutdown0: make(chan bool, 1),
		shutdown1: make(chan bool, 1),
	}
	r := gin.Default()
	s.r = r
	setup(s)
	return s
}

// NewProduction autoit command server relay
func NewProduction(path string) *Server {
	return New(path, false)
}

func setup(s *Server) {
	r := s.r
	// the route that the autoit program accesses to
	// see if there are any commands to be run
	r.GET("/sync", func(c *gin.Context) {
		select {
		case <-s.shutdown0:
			c.JSON(http.StatusOK, gin.H{
				"action": "shutdown",
			})
			s.shutdown1 <- true
		case cmd := <-s.tosend:
			c.JSON(http.StatusOK, gin.H{
				"action":  "command",
				"command": cmd,
			})
		case <-time.After(time.Millisecond * 100):
			c.String(http.StatusOK, "0")
		}
	})
	r.POST("/sync", func(c *gin.Context) {
		jd := &au3resp{}
		if err := c.BindJSON(jd); err != nil {
			fmt.Println("JSON error", err.Error())
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if jd.Type == "command" {
			if chan0, ok := s.toreceive[jd.CommandID]; ok {
				chan0 <- &Result{
					CommandID: jd.CommandID,
					Value:     jd.Value,
				}
				c.String(http.StatusOK, "1")
			} else {
				fmt.Println("POST /sync invalid ID", jd.CommandID)
				c.AbortWithStatus(http.StatusBadRequest)
			}
			return
		}
		c.AbortWithStatus(http.StatusNotFound)
	})
	r.GET("/health_check", func(c *gin.Context) {
		c.String(http.StatusOK, "1")
	})
}

func (s *Server) wait(id string) *Result {
	s.toreceive[id] = make(chan *Result)
	res := <-s.toreceive[id]
	return res
}

// RunHTTP server (blocks)
func (s *Server) RunHTTP() error {
	s.htps = &http.Server{
		Addr:    s.hostaddr,
		Handler: s.r,
	}
	return s.htps.ListenAndServe()
}

// WinGetTitle Retrieves the full title from a window.
//   title - The title/hWnd/class of the window to get the title.
//           See Title special definition:
//           https://www.autoitscript.com/autoit3/docs/intro/windowsadvanced.htm
//   text  - [optional] The text of the window to get the title.
//           Default is an empty string. See Text special definition:
//           https://www.autoitscript.com/autoit3/docs/intro/windowsbasic.htm#specialtext
func (s *Server) WinGetTitle(title, text string) string {
	cmd := newCommand("WinGetTitle")
	cmd.SetParams(title, text)
	s.tosend <- cmd
	result := s.wait(cmd.ID)
	rr := ""
	json.Unmarshal(result.Value, &rr)
	return rr
}

// Shutdown the http server
func (s *Server) Shutdown() {
	s.shutdown0 <- true
	<-s.shutdown1
	time.Sleep(time.Millisecond * 10)
	if s.htps != nil {
		s.htps.Shutdown(context.Background())
	}
}

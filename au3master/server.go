package au3master

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	trayindex int
	traychan  map[int]chan bool
	Stdout    io.Writer
}

// new autoit command server relay
func new(path string, devmode bool) *Server {
	if !devmode {
		gin.SetMode(gin.ReleaseMode)
	}
	s := &Server{
		tosend:    make(chan *Command, 64),
		toreceive: make(map[string]chan *Result),
		hostaddr:  path,
		shutdown0: make(chan bool, 1),
		shutdown1: make(chan bool, 1),
		traychan:  make(map[int]chan bool),
	}
	r := gin.New()
	r.Use(gin.Recovery())
	if devmode {
		r.Use(gin.Logger())
	}
	s.r = r
	setup(s)
	return s
}

// NewProduction autoit command server relay
func NewProduction(path string) *Server {
	return new(path, false)
}

// NewDevelopment autoit command server relay
func NewDevelopment(path string) *Server {
	return new(path, true)
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
	r.GET("/tray/:id", func(c *gin.Context) {
		if s.traychan == nil {
			c.String(http.StatusOK, "1")
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		if s.traychan[id] != nil {
			select {
			case s.traychan[id] <- true:
			default:
			}
		}
		c.String(http.StatusOK, "1")
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

func (s *Server) waitchan(id string) <-chan *Result {
	s.toreceive[id] = make(chan *Result)
	return s.toreceive[id]
}

// RunHTTP server (blocks)
func (s *Server) RunHTTP() error {
	s.htps = &http.Server{
		Addr:    s.hostaddr,
		Handler: s.r,
	}
	return s.htps.ListenAndServe()
}

// TestConnection tests the connection with the localhost http server
func (s *Server) TestConnection() bool {
	clcl := &http.Client{
		Timeout: time.Second * 2,
	}
	resp, err := clcl.Get(fmt.Sprintf("http://%s/health_check", s.hostaddr))
	if err != nil {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
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

package au3master

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server is a http server that relays commands to an
// autoit executable.
type Server struct {
	r         *gin.Engine
	tosend    chan *Command
	toreceive map[string]chan *Result
}

// New autoit command server relay
func New(path string, devmode bool) *Server {
	if !devmode {
		gin.SetMode(gin.ReleaseMode)
	}
	s := &Server{
		tosend:    make(chan *Command, 32),
		toreceive: make(map[string]chan *Result),
	}
	r := gin.Default()
	setup(r)
	s.r = r
	return s
}

func setup(r *gin.Engine) {
	// the route that the autoit program accesses to
	// see if there are any commands to be run
	r.GET("/sync", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"commands": make([]int, 0),
		})
	})
}

func (s *Server) wait(id string) *Result {
	s.toreceive[id] = make(chan *Result)
	res := <-s.toreceive[id]
	return res
}

// WinGetTitle Retrieves the full title from a window.
//   title - The title/hWnd/class of the window to get the title.
//           See Title special definition:
//           https://www.autoitscript.com/autoit3/docs/intro/windowsadvanced.htm
//   text  - [optional] The text of the window to get the title.
//           Default is an empty string. See Text special definition:
//           https://www.autoitscript.com/autoit3/docs/intro/windowsbasic.htm#specialtext
func (s *Server) WinGetTitle(title, text string) *Result {
	cmd := newCommand("WinGetTitle")
	cmd.SetParams(title, text)
	s.tosend <- cmd
	return s.wait(cmd.ID)
}

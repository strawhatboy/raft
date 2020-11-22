package core

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	logger  *log.Entry
	httpSrv *gin.Engine
	store   *Store
	config  *Config
}

func CreateServer(c *Config) (*Server, error) {


	s := &Server{
		logger:  GetLogger("server"),
		config:  c,
		httpSrv: gin.Default(),
	}
	store, err := NewStore(*c.SingleNode, c.ID, c.RaftAddr, c.RaftPath)
	if err != nil {
		s.logger.Error("failed to create store ", err)
		return nil, err
	}

	s.store = store

	s.logger.Debug("creating the server with config: ", c)
	s.httpSrv.GET("/:key", func(c *gin.Context) {
		k := c.Param("key")

		var v string
		var err error
		if v = s.store.Get(k); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":   100,
				"msg":    fmt.Sprintf("Failed to get value from key: %v", k),
				"result": nil,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":   0,
				"msg":    "Ok",
				"result": v,
			})
		}
	})

	s.httpSrv.PUT("/:key/:value", func(c *gin.Context) {
		k := c.Param("key")
		v := c.Param("value")

		if err := s.store.Put(k, v); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":   100,
				"msg":    fmt.Sprintf("Failed to put value: %v to key: %v", v, k),
				"result": nil,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":   0,
				"msg":    "Ok",
				"result": v,
			})
		}
	})

	s.httpSrv.POST("/join", func(c *gin.Context) {
		id := c.Param("id")
		addr := c.Param("addr")

		if err := s.store.Join(id, addr); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 101,
				"msg": fmt.Sprintf("failed to handle join request with %v %v, err: %v", id, addr, err),
				"result": nil,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg": "Ok",
				"result": nil,
			})
		}

	})

	s.logger.Debug("server created")
	return s, nil
}

func (s *Server) Run() {
	s.httpSrv.Run(s.config.HttpAddr)
}

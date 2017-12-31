package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperpilotio/snap-plugin-collector-ddagent/pkg/dogstatsd/message"
	log "github.com/sirupsen/logrus"
)

const HasBeenInitialized = "Server has been initialized"

// FIXME rewrite here, add flag
func init() {
	log.SetLevel(log.DebugLevel)
}

type Server struct {
	*http.Server
	isInitialized bool
}

func NewServer() *Server {
	return &Server{
		Server: &http.Server{
			Addr:    ":8000",
			Handler: router(),
		},
		isInitialized: false,
	}
}

func (srv *Server) Run(errCh chan<- error) {
	if srv.isInitialized {
		log.Debug(HasBeenInitialized)
		errCh <- errors.New(HasBeenInitialized)
		return
	}
	srv.isInitialized = true
	if err := srv.ListenAndServe(); err != nil {
		log.Error(err.Error())
		errCh <- err
		// } else {
		// close(errCh)
	}
}

func (srv *Server) Stop() (err error) {
	// FIXME srv.shutdown()
	err = srv.Close()
	srv.isInitialized = false
	return
}

func router() (r *gin.Engine) {
	r = gin.Default()
	v1 := r.Group("/v1")
	{
		v1.POST("/dogstatsd", dogstatsdHandler)
		// FIXME handle `check`
		// v1.POST("/checkMetrics", checkHandler)
		// FIXME handle `machine metrics`
		// v1.POST("/machineMetrics", machineHandler)
	}
	return
}

func dogstatsdHandler(c *gin.Context) {
	if rawData, err := c.GetRawData(); err == nil {
		message.Push(rawData)
		c.JSON(http.StatusOK, gin.H{"response": "received"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

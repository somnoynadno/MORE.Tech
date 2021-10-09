package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func handleOK(c *gin.Context, response interface{}) {
	j, _ := json.Marshal(response)
	log.Info(string(j))
	c.JSON(http.StatusOK, response)
}

func handleBadRequest(c *gin.Context, err error) {
	log.Warn(err.Error())
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func handleInternalError(c *gin.Context, err error) {
	log.Error(err.Error())
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

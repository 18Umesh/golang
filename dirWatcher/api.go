package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type API struct {
	config  *Config
	db      *gorm.DB
	watcher *Watcher
}

func NewAPI(config *Config, db *gorm.DB, watcher *Watcher) *API {
	return &API{
		config:  config,
		db:      db,
		watcher: watcher,
	}
}

func (api *API) RegisterRoutes(router *gin.Engine) {
	router.GET("/task-run", api.getTaskRun)
	router.POST("/configure", api.configure)
	router.POST("/start", api.start)
	router.POST("/stop", api.stop)
}

func (api *API) getTaskRun(c *gin.Context) {

	api.watcher.processFile(api.config.Directory)
	// Retrieve the latest task run details from the database
	var taskRun TaskRuns
	api.db.Last(&taskRun)

	runtimeInSeconds := int64(time.Since(taskRun.StartTime).Seconds())

    // Save the updated runtime to the database
    api.db.Model(&taskRun).Update("runtime", runtimeInSeconds)

	c.JSON(http.StatusOK, gin.H{
		"startTime":       taskRun.StartTime,
		"endTime":         taskRun.EndTime,
		"runtime":         runtimeInSeconds,
		"filesAdded":      taskRun.FilesAdded,
		"filesDeleted":    taskRun.FilesDeleted,
		"magicStringHits": taskRun.MagicStringHits,
		"status":          taskRun.Status,
	})
}

func (api *API) configure(c *gin.Context) {
	var request struct {
		Directory    string        `json:"directory"`
		TimeInterval string `json:"timeInterval"`
		MagicString  string        `json:"magicString"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update the configuration
	api.config.Directory = request.Directory
	api.config.TimeInterval = request.TimeInterval
	api.config.MagicString = request.MagicString

	c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully."})
}

func (api *API) start(c *gin.Context) {
	// Manually start the background task
	api.watcher.Start()

	c.JSON(http.StatusOK, gin.H{"message": "Task started successfully."})
}

func (api *API) stop(c *gin.Context) {
	// Manually stop the background task
	api.watcher.Stop()

	c.JSON(http.StatusOK, gin.H{"message": "Task stopped successfully."})
}

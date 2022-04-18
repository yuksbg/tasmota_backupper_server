package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	ginlogrus "github.com/toorop/gin-logrus"
	"log"
	"tasmota_backup/helpers"
	"tasmota_backup/modules/devices"
	"tasmota_backup/modules/download"
	"tasmota_backup/modules/files"
	"tasmota_backup/tasks/backup"
)

//go:embed fe_ui
var FeUI embed.FS

func init() {
	helpers.PrepareTables()
	go backup.Run()
}
func main() {
	helpers.GetLogger().Info("App is starting")
	r := gin.New()
	r.Use(ginlogrus.Logger(helpers.GetLogger()), gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		ind, _ := FeUI.ReadFile("fe_ui/index.html")
		c.Data(200, "text/html", ind)
	})

	r.GET("/css/*any", func(c *gin.Context) {
		log.Println(c.Request.RequestURI)
		cont, er := FeUI.ReadFile(fmt.Sprintf("fe_ui%s", c.Request.RequestURI))
		if er != nil {
			c.JSON(404, gin.H{"error": er.Error()})
			return
		}
		c.Data(200, "text/css", cont)
	})
	r.GET("/js/*any", func(c *gin.Context) {
		cont, er := FeUI.ReadFile(fmt.Sprintf("fe_ui%s", c.Request.RequestURI))
		if er != nil {
			c.JSON(404, gin.H{"error": er.Error()})
			return
		}
		c.Data(200, "application/javascript; charset=UTF-8", cont)
	})

	api := r.Group("/api")
	{
		api.GET("/discover_devices", devices.DiscoverHandle)
		api.GET("/get_devices", devices.ListHandle)

		api.POST("devices/delete/:mac", devices.DeleteHandle)
		api.POST("devices/add/by_ip/:ip", devices.AddDeviceByIp)
		api.POST("devices/clear_all", devices.DeleteAll)

		api.POST("backup/manual/:mac", devices.HandleManualRunBackup)
		api.POST("backup/all", devices.HandleBackupAll)
		api.GET("backup/get_backups/:mac", devices.HandleGetBackups)
		api.GET("backup/download/:mac/:filename", download.HandleDownload)
		api.POST("backup/delete/:mac/:filename", files.HandleDelete)
	}

	if err := r.Run(); err != nil {
		helpers.GetLogger().WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Starting app error")
	}

}

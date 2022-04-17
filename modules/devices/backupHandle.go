package devices

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
	"tasmota_backup/helpers"
	"tasmota_backup/tasks/backup"
)

type backupFilesR struct {
	FileName   string `json:"file_name,omitempty"`
	FileSize   int64  `json:"file_size,omitempty"`
	BackupTime string `json:"backup_time,omitempty"`
}

func HandleManualRunBackup(c *gin.Context) {

	mac := c.Param("mac")
	if mac == "" {
		c.JSON(200, gin.H{"result": false})
		return
	}
	//
	var deviceIP string
	sql := `select device_ip from tasmota_devices where mac_address = ? limit 1`
	helpers.GetDb().Get(&deviceIP, sql, mac)
	if deviceIP == "" {
		c.JSON(200, gin.H{"result": false})
		return
	}
	backup.RunBackup(deviceIP, mac)

	c.JSON(200, gin.H{"result": true})
}

func HandleGetBackups(c *gin.Context) {
	mac := c.Param("mac")
	if mac == "" {
		c.JSON(200, gin.H{"result": false})
		return
	}
	folder := fmt.Sprintf("backup_data/%s", mac)
	fContent, _ := ioutil.ReadDir(folder)

	var data []backupFilesR

	for _, info := range fContent {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".dmp") {
			data = append(data, backupFilesR{
				FileName:   info.Name(),
				FileSize:   info.Size(),
				BackupTime: strings.ReplaceAll(info.Name(), ".dmp", ""),
			})
		}
	}

	c.JSON(200, gin.H{
		"files": data,
	})

}

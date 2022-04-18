package devices

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
	"sync"
	"tasmota_backup/helpers"
	"tasmota_backup/tasks/backup"
)

type backupFilesR struct {
	FileName   string `json:"file_name,omitempty"`
	FileSize   int64  `json:"file_size,omitempty"`
	BackupTime string `json:"backup_time,omitempty"`
}

func HandleBackupAll(c *gin.Context) {
	type backupMe struct {
		Ip  string `db:"ip"`
		Mac string `db:"mac"`
	}

	var devices []backupMe

	sql := `select device_ip as ip, mac_address as mac from tasmota_devices`
	if er := helpers.GetDb().Select(&devices, sql); er != nil {
		helpers.GetLogger().WithFields(map[string]interface{}{
			"error": er.Error(),
		}).Error("HandleBackupAll")
	}
	var wg sync.WaitGroup
	helpers.GetLogger().WithFields(map[string]interface{}{
		"devices": devices,
	}).Info("DEBU")
	for _, device := range devices {
		wg.Add(1)
		go func(device backupMe) {
			defer wg.Done()
			backup.RunBackup(device.Ip, device.Mac)
		}(device)
	}
	wg.Wait()

	c.JSON(200, gin.H{"result": true})
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

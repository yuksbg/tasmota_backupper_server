package devices

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"tasmota_backup/helpers"
	"tasmota_backup/modules/device_info"
	"tasmota_backup/modules/discover_ip"
)

func DiscoverHandle(c *gin.Context) {

	cidr, _ := c.GetQuery("cidr")

	helpers.GetLogger().Debug("Start discovery process")
	devices := discover_ip.Discover(cidr)

	helpers.GetLogger().WithFields(map[string]interface{}{
		"devices": devices,
	}).Info("Found devices")

	// filter devices
	var existDevices []interface{}
	for _, device := range devices {
		var exist int64
		sql := `select count(mac_address) as broy from tasmota_devices where mac_address = ?`
		helpers.GetDb().Get(&exist, sql, device.Mac)
		if exist != 0 {
			existDevices = append(existDevices, map[string]string{"mac": device.Mac})
		}
	}

	c.JSON(200, gin.H{
		"devices":  devices,
		"existing": existDevices,
	})

}

func ListHandle(c *gin.Context) {
	sql := `SELECT 
       				mac_address as  mac, 
       				device_ip as ip, 
       				device_name as name 
			FROM tasmota_devices ORDER BY mac_address DESC`

	var devices []device_info.TasmoDevice
	if er := helpers.GetDb().Select(&devices, sql); er != nil {
		log.Println(er.Error())
	}

	backups := map[string]interface{}{}
	backupFolder, _ := ioutil.ReadDir("backup_data")
	for _, info := range backupFolder {
		if info.IsDir() {
			bf, _ := ioutil.ReadDir(fmt.Sprintf("%s/%s", "backup_data", info.Name()))
			lastBackup := ""
			if len(bf) != 0 {
				log.Println("BBBB", bf[0].Name())
				sort.SliceStable(bf, func(i, j int) bool {
					return bf[i].Name() > bf[j].Name()
				})
				lastBackup = strings.ReplaceAll(bf[0].Name(), ".dmp", "")
			}
			backups[info.Name()] = map[string]interface{}{
				"file_count":  len(bf),
				"last_backup": lastBackup,
			}

		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"devices": devices,
		"backups": backups,
	})
}

func DeleteAll(c *gin.Context) {
	sql := `delete from tasmota_devices`
	_, err := helpers.GetDb().Exec(sql)
	if err != nil {
		helpers.GetLogger().WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Error in deleting devices")
		c.JSON(200, gin.H{"result": false, "error": err.Error()})
		return
	}

	// delete backup files
	os.RemoveAll("backup_data/")
	os.MkdirAll("backup_data", os.ModePerm)

	c.JSON(200, gin.H{"result": true})
}

func DeleteHandle(c *gin.Context) {
	mac := c.Param("mac")
	sql := `delete from tasmota_devices where mac_address = ?`
	helpers.GetDb().Exec(sql, mac)
	c.JSON(200, gin.H{"result": true})
}

func AddDeviceByIp(c *gin.Context) {
	ip := c.Param("ip")

	dev, er := device_info.GetDeviceInfo(ip)
	if er != nil {
		c.JSON(200, gin.H{
			"result": false,
			"error":  er.Error(),
		})
		return
	}
	sql := `insert into tasmota_devices (mac_address,device_ip,device_name,tasmota_version,last_check) values(?,?,?,?,datetime('now'))`
	_, er = helpers.GetDb().Exec(sql, dev.Mac, dev.IP, dev.Name, dev.FirmwareVersion)
	if er != nil {
		helpers.GetLogger().Error(er.Error())
		c.JSON(200, gin.H{"result": false, "error": er.Error()})
		return
	}
	c.JSON(200, gin.H{"result": true})
}

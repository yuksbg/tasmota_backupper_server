package backup

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"tasmota_backup/helpers"
	"time"
)

func Run() {
	//runBackup("192.168.88.19", "9C:9C:1F:47:96:87")
}

func RunBackup(deviceIP string, mac string) {
	url := fmt.Sprintf("http://%s/dl?", deviceIP)

	resp, err := http.Get(url)
	if err != nil {
		helpers.GetLogger().WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Downloading error")
		return
	}
	defer resp.Body.Close()

	//@TODO Check for later fix
	fileName := fmt.Sprintf("backup_data/%s/%d.dmp", mac, time.Now().Unix())

	if err = os.MkdirAll(fmt.Sprintf("backup_data/%s", mac), os.ModePerm); err != nil {
		helpers.GetLogger().WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Error in Creating SubDir")
		return
	}

	f, err := os.Create(fileName)
	if err != nil {
		helpers.GetLogger().WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Error in Creating file")
		return
	}
	defer f.Close()

	io.Copy(f, resp.Body)
	fsize, _ := os.Stat(f.Name())
	helpers.GetLogger().WithFields(map[string]interface{}{
		"device_ip": deviceIP,
		"fileSize":  fsize.Size(),
		"fileName":  f.Name(),
	}).Info("Backup successful")

}

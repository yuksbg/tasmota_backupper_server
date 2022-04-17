package device_info

import (
	"errors"
	"github.com/tidwall/gjson"
	_ "github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var parseFirmwareRegex = regexp.MustCompile(`(.*)\((.*)\)`)

func parseFirmwareVersion(v string) (string, string, error) {
	res := parseFirmwareRegex.FindAllStringSubmatch(v, 1)
	if len(res) != 1 {
		return "", "", errors.New("Regex parser failed\n" + v)
	}
	return res[0][1], res[0][2], nil
}

func GetDeviceInfo(ip string) (device TasmoDevice, err error) {

	url := "http://" + ip + "/cm?cmnd=Status%200" // build url

	info, err := getAPIInfo(url)
	if err == nil && info != "" {
		// validate json
		if valid := gjson.Valid(info); !valid {
			err = errors.New("not valid json")
			return
		}
		// get version
		fw := gjson.Get(info, "StatusFWR.Version").String()
		// extract details
		version, variant, er := parseFirmwareVersion(fw)
		if er != nil {
			err = er
			return
		}
		// set data for returning
		device.IP = ip
		device.FirmwareVersion = version
		device.FirmwareType = variant
		device.Name = gjson.Get(info, "Status.DeviceName").String()
		device.Mac = gjson.Get(info, "StatusNET.Mac").String()
	}

	return
}

func getAPIInfo(url string) (string, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", errors.New("JSON download failed")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

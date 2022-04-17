package device_info

type TasmoDevice struct {
	Name            string `json:"name,omitempty" db:"name"`
	FirmwareVersion string `json:"firmware_version,omitempty" db:"firmware_version"`
	FirmwareType    string `json:"firmware_type,omitempty" db:"firmware_type"`
	Outdated        bool   `json:"outdated,omitempty" db:"outdated"`
	IP              string `json:"ip,omitempty" db:"ip"`
	Mac             string `json:"mac,omitempty" db:"mac"`
}

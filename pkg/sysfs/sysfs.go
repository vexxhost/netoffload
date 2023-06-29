package sysfs

import (
	"fmt"
	"strconv"

	"github.com/spf13/afero"
)

func readInt(fs afero.Fs, deviceName string, attrName string) (int, error) {
	path := fmt.Sprintf("/sys/class/net/%s/%s", deviceName, attrName)
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(data))
}

func writeInt(fs afero.Fs, deviceName string, attrName string, value int) error {
	path := fmt.Sprintf("/sys/class/net/%s/%s", deviceName, attrName)
	data := []byte(strconv.Itoa(value))
	return afero.WriteFile(fs, path, data, 0644)
}

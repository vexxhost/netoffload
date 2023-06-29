package sysfs

import (
	"fmt"

	"github.com/spf13/afero"
)

type InterfaceConfig struct {
	Name       string
	Filesystem afero.Fs
}

func Get(deviceName string) (*InterfaceConfig, error) {
	fs := afero.NewOsFs()
	return GetWithFs(fs, deviceName)
}

func GetWithFs(fs afero.Fs, deviceName string) (*InterfaceConfig, error) {
	return &InterfaceConfig{
		Name:       deviceName,
		Filesystem: fs,
	}, nil
}

func (c *InterfaceConfig) GetVfsCount() (int, error) {
	return readInt(c.Filesystem, c.Name, "sriov_numvfs")
}

func (c *InterfaceConfig) SetVfsCount(count int) error {
	vfsCount, err := c.GetVfsCount()
	if err != nil {
		return err
	}

	if count == vfsCount {
		return nil
	}

	maxVfs, err := c.GetVfsMax()
	if err != nil {
		return err
	}

	if count > maxVfs {
		return fmt.Errorf("sysfs: cannot set VFs count to %d, max is %d", count, maxVfs)
	}

	return writeInt(c.Filesystem, c.Name, "sriov_numvfs", count)
}

func (c *InterfaceConfig) GetVfsMax() (int, error) {
	return readInt(c.Filesystem, c.Name, "sriov_totalvfs")
}

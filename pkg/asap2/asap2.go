package asap2

import (
	"fmt"

	"github.com/vexxhost/netoffload/pkg/cmdline"
	"github.com/vexxhost/netoffload/pkg/mstconfig"
)

func Enable(devName string, vfsCount int) error {
	err := cmdline.IsIommuEnabled()
	if err != nil {
		return err
	}

	pciAddress := "detect me"

	cardConfig, err := mstconfig.Query(pciAddress)
	if err != nil {
		return err
	}

	if !cardConfig.SrIov.Enabled || cardConfig.SrIov.VfsCount != vfsCount {
		// TODO: enable sr-iov and update vfs count (and request reboot?)
		fmt.Println("blah")
	}

	// update numvfs in sysfs
	// check if device is in switchdev mode and hw-tc-offload is enabled
	// if it is, do nothing
	// if it is not, unbind/switchdev/enabled hw-tc-offload/bind
	// enable hwoffload in ovs

	return nil
}

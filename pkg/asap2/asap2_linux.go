package asap2

import (
	"fmt"

	"github.com/k8snetworkplumbingwg/sriovnet"
	"github.com/vishvananda/netlink"

	"github.com/vexxhost/netoffload/pkg/cmdline"
	"github.com/vexxhost/netoffload/pkg/mstconfig"
	"github.com/vexxhost/netoffload/pkg/sysfs"
)

func Enable(devName string, vfsCount int) error {
	err := cmdline.IsIommuEnabled()
	if err != nil {
		return err
	}

	pciAddress, err := sriovnet.GetPciFromNetDevice(devName)
	if err != nil {
		return err
	}

	cardConfig, err := mstconfig.Query(pciAddress)
	if err != nil {
		return err
	}

	if !cardConfig.SrIov.Enabled || cardConfig.SrIov.VfsCount != vfsCount {
		// TOOD(mnaser): Enable SR-IOV and set VFs count using `mstconfig`
		return fmt.Errorf("sr-iov is not enabled or vfs count is not %d", vfsCount)
	}

	iface, err := sysfs.Get(devName)
	if err != nil {
		return err
	}

	err = iface.SetVfsCount(vfsCount)
	if err != nil {
		return err
	}

	eswitchMode, err := GetDevlinkEswitchMode(pciAddress)
	if err != nil {
		return err
	}

	fmt.Println(eswitchMode)

	// check if device is in switchdev mode and hw-tc-offload is enabled
	// if it is, do nothing
	// if it is not, unbind/switchdev/enabled hw-tc-offload/bind
	// enable hwoffload in ovs

	return nil
}

func GetDevlinkEswitchMode(pciAddress string) (string, error) {
	devLink, err := netlink.DevLinkGetDeviceByName("pci", pciAddress)
	if err != nil {
		return "", err
	}

	return devLink.Attrs.Eswitch.Mode, nil
}

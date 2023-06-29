package asap2

import (
	"context"
	"fmt"
	"log"

	"github.com/k8snetworkplumbingwg/sriovnet"
	"github.com/ovn-org/libovsdb/client"
	"github.com/ovn-org/libovsdb/model"
	"github.com/safchain/ethtool"
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

	log.Printf("asap2: detected pci address %s", pciAddress)

	err = applyMstConfig(pciAddress, vfsCount)
	if err != nil {
		return err
	}

	err = applyVfsCount(devName, vfsCount)
	if err != nil {
		return err
	}

	err = applyDevlinkEswitchMode(devName, pciAddress)
	if err != nil {
		return err
	}

	// TODO: check iommu groups

	err = applyHwTcOffload(devName)
	if err != nil {
		return err
	}

	err = applyOvsHwOffload()
	if err != nil {
		return err
	}

	return nil
}

func applyMstConfig(pciAddress string, vfsCount int) error {
	cardConfig, err := mstconfig.Query(pciAddress)
	if err != nil {
		return err
	}

	log.Printf("mstconfig: SRIOV_EN = %t (expected: true)", cardConfig.SrIov.Enabled)
	log.Printf("mstconfig: NUM_OF_VFS = %d (expected: %d)", cardConfig.SrIov.VfsCount, vfsCount)

	if !cardConfig.SrIov.Enabled || cardConfig.SrIov.VfsCount != vfsCount {
		return fmt.Errorf("mstconfig: card misconfigured")
	}

	return nil
}

func applyVfsCount(devName string, vfsCount int) error {
	iface, err := sysfs.Get(devName)
	if err != nil {
		return err
	}

	currentVfsCount, err := iface.GetVfsCount()
	if err != nil {
		return err
	}

	if currentVfsCount != vfsCount {
		log.Printf("sriov_numvfs: changing to %d", vfsCount)

		return iface.SetVfsCount(vfsCount)
	}

	log.Printf("sriov_numvfs: already set to %d", vfsCount)

	return nil
}

func applyDevlinkEswitchMode(devName string, pciAddress string) error {
	devLink, err := netlink.DevLinkGetDeviceByName("pci", pciAddress)
	if err != nil {
		return err
	}

	netdev, err := sriovnet.GetPfNetdevHandle(devName)
	if err != nil {
		return err
	}

	if devLink.Attrs.Eswitch.Mode != "switchdev" {
		log.Println("devlink: setting eswitch mode to switchdev")

		// TODO: detect if anything is allocated

		for _, vf := range netdev.List {
			if !vf.Bound {
				continue
			}

			err = sriovnet.UnbindVf(netdev, vf)
			if err != nil {
				return fmt.Errorf("devlink: failed to unbind VF %d: %w", vf.Index, err)
			}
		}

		err = netlink.DevLinkSetEswitchMode(devLink, "switchdev")
		if err != nil {
			return fmt.Errorf("devlink: failed to set eswitch mode: %w", err)
		}

		for _, vf := range netdev.List {
			err = sriovnet.BindVf(netdev, vf)
			if err != nil {
				return fmt.Errorf("devlink: failed to bind VF %d: %w", vf.Index, err)
			}
		}

		return nil
	}

	log.Println("devlink: eswitch mode is already set to switchdev")

	return nil
}

func applyHwTcOffload(devName string) error {
	e, err := ethtool.NewEthtool()
	if err != nil {
		return err
	}
	defer e.Close()

	features, err := e.Features(devName)
	if err != nil {
		return err
	}

	if _, ok := features["hw-tc-offload"]; !ok {
		return fmt.Errorf("hw-tc-offload: device not reporting feature")
	}

	if !features["hw-tc-offload"] {
		log.Println("hw-tc-offload: enabling")

		return e.Change(devName, map[string]bool{
			"hw-tc-offload": true,
		})
	}

	log.Println("hw-tc-offload: already enabled")

	return nil
}

type ovsType struct {
	UUID        string            `ovsdb:"_uuid"`
	OtherConfig map[string]string `ovsdb:"other_config"`
}

func applyOvsHwOffload() error {
	ctx := context.Background()

	clientDBModel, err := model.NewClientDBModel("Open_vSwitch", map[string]model.Model{
		"Open_vSwitch": &ovsType{},
	})
	if err != nil {
		return err
	}

	ovs, err := client.NewOVSDBClient(clientDBModel, client.WithEndpoint("unix:/run/openvswitch/db.sock"))
	if err != nil {
		return err
	}

	if err = ovs.Connect(ctx); err != nil {
		return err
	}
	defer ovs.Disconnect()

	if _, err := ovs.MonitorAll(ctx); err != nil {
		return err
	}

	var ovsConfig *ovsType
	for _, data := range ovs.Cache().Table("Open_vSwitch").Rows() {
		ovsConfig = data.(*ovsType)
	}

	offloadEnabled, ok := ovsConfig.OtherConfig["hw-offload"]
	if !ok || offloadEnabled != "true" {
		// TODO: enable and restart ovs somehow?
		return fmt.Errorf("ovs: other-config:hw-offload not enabled")
	}

	log.Println("ovs: hw-offload already enabled")

	return nil
}

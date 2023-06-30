package sysfs

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGetVfsCount(t *testing.T) {
	tests := map[string]int{
		"0": 0,
		"8": 8,
	}

	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			err := writeInt(fs, "eth0", "sriov_numvfs", value)
			assert.NoError(t, err)

			config, err := GetWithFs(fs, "eth0")
			assert.NoError(t, err)

			vfsCount, err := config.GetVfsCount()
			assert.NoError(t, err)
			assert.Equal(t, value, vfsCount)
		})
	}
}

func TestSetVfsMax(t *testing.T) {
	tests := map[string]struct {
		maxVfs int
		vfs    int
		err    error
	}{
		"max vfs": {
			maxVfs: 8,
			vfs:    8,
			err:    nil,
		},
		"max vfs - 1": {
			maxVfs: 8,
			vfs:    7,
			err:    nil,
		},
		"max vfs + 1": {
			maxVfs: 8,
			vfs:    9,
			err:    fmt.Errorf("sysfs: cannot set VFs count to 9, max is 8"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			err := writeInt(fs, "eth0", "sriov_numvfs", 0)
			assert.NoError(t, err)
			err = writeInt(fs, "eth0", "sriov_totalvfs", test.maxVfs)
			assert.NoError(t, err)

			config, err := GetWithFs(fs, "eth0")
			assert.NoError(t, err)

			err = config.SetVfsCount(test.vfs)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestGetVfsMax(t *testing.T) {
	tests := map[string]int{
		"0": 0,
		"8": 8,
	}

	for name, value := range tests {
		t.Run(name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			writeInt(fs, "eth0", "sriov_totalvfs", value)

			config, err := GetWithFs(fs, "eth0")
			assert.NoError(t, err)

			vfsCount, err := config.GetVfsMax()
			assert.NoError(t, err)
			assert.Equal(t, value, vfsCount)
		})
	}
}

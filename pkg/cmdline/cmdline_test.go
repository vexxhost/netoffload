package cmdline

import (
	"io"
	"testing"

	// NOTE(mnaser): This is a hack to get around the fact that the cmdline package
	//               is not exported.  We are only using this for tests to make sure
	//               we can detect the iommu flag.
	_ "unsafe"

	"github.com/klauspost/cpuid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/u-root/u-root/pkg/cmdline"
)

//go:linkname parse cmdline.parse
func parse(cmdlineReader io.Reader) *cmdline.CmdLine

func TestIsIommuEnabledForConfig(t *testing.T) {
	tests := map[cpuid.Vendor]map[string]string{
		cpuid.Intel: {
			"iommu":       "pt",
			"intel_iommu": "on",
		},
		cpuid.AMD: {
			"iommu":     "pt",
			"amd_iommu": "on",
		},
	}

	for vendor, cmdLineMap := range tests {
		config := &Config{
			Vendor: vendor,
			CmdLine: &cmdline.CmdLine{
				AsMap: cmdLineMap,
			},
		}

		err := IsIommuEnabledForConfig(config)
		assert.NoError(t, err)
	}
}

func TestIsIommuEnabledForConfigNegative(t *testing.T) {
	tests := map[cpuid.Vendor]map[string]string{
		cpuid.KVM: {
			"iommu": "pt",
		},
		cpuid.Intel: {
			"iommu":     "pt",
			"amd_iommu": "on",
		},
		cpuid.AMD: {
			"iommu":       "pt",
			"intel_iommu": "on",
		},
	}

	for vendor, cmdLineMap := range tests {
		config := &Config{
			Vendor: vendor,
			CmdLine: &cmdline.CmdLine{
				AsMap: cmdLineMap,
			},
		}

		err := IsIommuEnabledForConfig(config)
		assert.Error(t, err)
	}
}

func TestAssertValues(t *testing.T) {
	c := &cmdline.CmdLine{
		AsMap: map[string]string{
			"foo": "bar",
		},
	}

	err := assertValue(c, "foo", "bar")
	require.NoError(t, err)
}

func TestAssertValuesNegative(t *testing.T) {
	c := &cmdline.CmdLine{
		AsMap: map[string]string{
			"foo": "baz",
		},
	}

	err := assertValue(c, "foo", "bar")
	assert.Error(t, err)
}

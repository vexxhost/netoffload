package cmdline

import (
	"errors"
	"fmt"

	"github.com/klauspost/cpuid/v2"
	"github.com/u-root/u-root/pkg/cmdline"
)

type Config struct {
	cpuid.Vendor
	*cmdline.CmdLine
}

// IsIommuEnabled is identical to IsIommuEnabledFromCmdline but uses the
// default cmdline reader.
func IsIommuEnabled() error {
	return IsIommuEnabledForConfig(
		&Config{
			Vendor:  cpuid.CPU.VendorID,
			CmdLine: cmdline.NewCmdLine(),
		},
	)
}

// IsIommuEnabledForConfig returns no error if the kernel command line contains
// the "iommu=pt" options and the vendor-specific "intel_iommu" or "amd_iommu"
// options.
func IsIommuEnabledForConfig(config *Config) error {
	var iommuErr, vendorIommuErr error

	iommuErr = assertValue(config.CmdLine, "iommu", "pt")

	switch config.Vendor {
	case cpuid.Intel:
		vendorIommuErr = assertValue(config.CmdLine, "intel_iommu", "on")
	case cpuid.AMD:
		vendorIommuErr = assertValue(config.CmdLine, "amd_iommu", "on")
	default:
		vendorIommuErr = fmt.Errorf("cmdline: unknown CPU vendor %q", config.Vendor)
	}

	return errors.Join(iommuErr, vendorIommuErr)
}

func assertValue(cmd *cmdline.CmdLine, flag string, expected string) error {
	val, ok := cmd.Flag(flag)

	if !ok {
		return fmt.Errorf("cmdline: %s flag not found (expected: %q)", flag, expected)
	}

	if val != expected {
		return fmt.Errorf("cmdline: %s flag is not %q (got: %q)", flag, expected, val)
	}

	return nil
}

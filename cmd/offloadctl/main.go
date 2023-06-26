package main

import (
	"fmt"

	"github.com/vexxhost/netoffload/pkg/cmdline"
)

func main() {
	err := cmdline.IsIommuEnabled()
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Check if IOMMU groups for VFs are not the same

	fmt.Println("Hello, playground")
}

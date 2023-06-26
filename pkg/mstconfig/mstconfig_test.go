package mstconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	f, err := os.Open("testdata/query-sriov-enabled-with-16-vfs.txt")
	assert.NoError(t, err)

	config, err := parseQuery(f)
	assert.NoError(t, err)

	assert.Equal(t, &Config{
		SrIov: SrIovConfig{
			Enabled:  true,
			VfsCount: 16,
		},
	}, config)
}

func TestParseConfig(t *testing.T) {
	tests := map[string]struct {
		line string
		key  string
		val  string
	}{
		"SRIOV_EN": {
			line: "         SRIOV_EN                            True(1)",
			key:  "SRIOV_EN",
			val:  "True(1)",
		},
		"NUM_OF_VFS": {
			line: "         NUM_OF_VFS                          16",
			key:  "NUM_OF_VFS",
			val:  "16",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			key, val := parseConfig(test.line)
			assert.Equal(t, test.key, key)
			assert.Equal(t, test.val, val)
		})
	}
}

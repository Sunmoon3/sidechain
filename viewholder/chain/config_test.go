package chain

import (
	"fmt"
	"github.com/test-go/testify/assert"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configpath := "../config.yaml"
	config, err := LoadConfig(configpath)
	fmt.Println(config)
	assert.NoError(t, err)
}

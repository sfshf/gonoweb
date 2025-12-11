package config_test

import (
	"context"
	"os"
	"testing"

	"github.com/sfshf/gonoweb/internal/config"
)

// go test -v -count=1 -timeout 30s -run ^TestInitAppConfig$ github.com/sfshf/gonoweb/config
func TestInitAppConfig(t *testing.T) {
	pwd, _ := os.Getwd()
	t.Logf("pwd: %s", pwd)
	if err := config.InitAppConfig(context.TODO(), "../../config/dev/gono.toml"); err != nil {
		t.Fatal(err)
	}
	t.Log("Config is OK !")
}

package repo_test

import (
	"context"
	"testing"

	"github.com/sfshf/gonoweb/internal/config"
	"github.com/sfshf/gonoweb/internal/repo"
)

// go test -v -count=1 -timeout 30s -run ^TestInitGorm$ github.com/sfshf/gonoweb/internal/repo
func TestInitGorm(t *testing.T) {
	ctx := context.TODO()
	if err := config.InitAppConfig(ctx, "../../config/dev/gono.toml"); err != nil {
		t.Fatal(err)
	}
	if err := repo.InitGorm(ctx); err != nil {
		t.Fatal(err)
	}
	t.Log("Gorm is OK !")
}

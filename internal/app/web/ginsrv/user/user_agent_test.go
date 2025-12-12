package user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sfshf/gonoweb/internal/app/web/ginsrv"
	"github.com/sfshf/gonoweb/internal/config"
	"github.com/sfshf/gonoweb/internal/repo"
	"github.com/stretchr/testify/assert"
)

// go test -v -count=1 -timeout 30s -run ^TestVisit$ github.com/sfshf/gonoweb/internal/app/web/ginsrv/user
func TestVisit(t *testing.T) {
	ctx := context.TODO()
	if err := config.InitAppConfig(ctx, "../../../../../config/dev/gono.toml"); err != nil {
		t.Fatal(err)
	}
	r, err := ginsrv.InitGin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.InitGorm(ctx); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/visit", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

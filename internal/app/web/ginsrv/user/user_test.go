package user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sfshf/gonoweb/internal/app/web/ginsrv"
	"github.com/sfshf/gonoweb/internal/config"
	"github.com/sfshf/gonoweb/internal/repo"
	"github.com/sfshf/gonoweb/internal/service/casbin"
	"github.com/sfshf/gonoweb/internal/service/user"
	"github.com/stretchr/testify/assert"
)

// go test -v -count=1 -timeout 30s -run ^TestSignInOut$ github.com/sfshf/gonoweb/internal/app/web/ginsrv/user
func TestSignInOut(t *testing.T) {
	ctx := context.TODO()
	if err := config.InitAppConfig(
		ctx,
		"../../../../../config/dev/gono.toml",
		"../../../../../config/dev/rbac_with_domains_model.conf",
	); err != nil {
		t.Fatal(err)
	}
	if err := repo.InitGorm(ctx); err != nil {
		t.Fatal(err)
	}
	// 初始化前置服务
	//   - 初始化root账号
	clear, err := user.Launch()
	if err != nil {
		t.Fatal(err)
	}
	defer clear()
	//   - casbin service
	clear, err = casbin.Launch()
	if err != nil {
		t.Fatal(err)
	}
	defer clear()
	r, err := ginsrv.InitGin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	signInParams := `{
		"email": "root@gono.com",
		"password": "9eae7836cb0d4145e15321d018814196"
	}`
	signInW := httptest.NewRecorder()
	signInReq, _ := http.NewRequest(http.MethodPost, "/api/v1/user/signIn", strings.NewReader(signInParams))
	r.ServeHTTP(signInW, signInReq)
	assert.Equal(t, http.StatusOK, signInW.Code)

	token := signInW.Header().Get("Authorization")
	t.Logf("Response Header [Authorization]: %s\n", token)

	signOutW := httptest.NewRecorder()
	signOutReq, _ := http.NewRequest(http.MethodPost, "/api/v1/user/signOut", nil)
	signOutReq.Header.Set("Authorization", token)
	r.ServeHTTP(signOutW, signOutReq)
	assert.Equal(t, http.StatusOK, signOutW.Code)
	t.Logf("Response Body: %s\n", signOutW.Body.String())
}

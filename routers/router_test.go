package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wangyi1310/mycloud-disk/conf"
)

func TestPing(t *testing.T) {
	t.Log("test ping")
	asserts := assert.New(t)
	router := InitMaster()
	rsp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v3/site/ping", nil)
	router.ServeHTTP(rsp, req)
	asserts.Equal(rsp.Code, 200)
	asserts.Contains(rsp.Body.String(), conf.BackendVersion)
}

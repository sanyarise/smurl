package delivery

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/sanyarise/smurl/internal/infrastructure/api/handler"
	"github.com/sanyarise/smurl/internal/infrastructure/db/mockdb"
	"github.com/sanyarise/smurl/internal/usecases/repos/smurlrepo"

	"go.uber.org/zap"
)

func TestPostCreate(t *testing.T) {
	mdb := mockdb.NewMockdb("", zap.L())
	srs := smurlrepo.NewSmurlStorage(mdb, zap.L())
	hs := handler.NewHandlers(srs, zap.L())
	roa := NewRouterOpenAPI(hs, zap.L(), "")
	srv := httptest.NewServer(roa)
	buffer := new(bytes.Buffer)
	params := url.Values{}
	params.Set("long_url", "http://vk.com")
	buffer.WriteString(params.Encode())

	r, _ := http.NewRequest("POST", srv.URL+"/create", buffer)
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	cli := srv.Client()

	resp, err := cli.Do(r)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected statuscode %d, got statuscode %d", http.StatusCreated, resp.StatusCode)
	}
	resp.Body.Close()
	buffer = new(bytes.Buffer)
	params = url.Values{}
	params.Set("long_url", "vk.com")
	buffer.WriteString(params.Encode())

	r, _ = http.NewRequest("POST", srv.URL+"/create", buffer)
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	cli = srv.Client()

	resp, err = cli.Do(r)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected statuscode %d, got statuscode %d", http.StatusBadRequest, resp.StatusCode)
	}
	defer resp.Body.Close()

	buffer = new(bytes.Buffer)
	params = url.Values{}
	params.Set("long_url", "")
	buffer.WriteString(params.Encode())

	r, _ = http.NewRequest("POST", srv.URL+"/create", buffer)
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	cli = srv.Client()

	resp, err = cli.Do(r)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected statuscode %d, got statuscode %d", http.StatusBadRequest, resp.StatusCode)
	}
	defer resp.Body.Close()
}

func TestPostStat(t *testing.T) {
	mdb := mockdb.NewMockdb("", zap.L())
	srs := smurlrepo.NewSmurlStorage(mdb, zap.L())
	hs := handler.NewHandlers(srs, zap.L())
	roa := NewRouterOpenAPI(hs, zap.L(), "")
	srv := httptest.NewServer(roa)

	buffer := new(bytes.Buffer)
	params := url.Values{}
	params.Set("admin_url", "zxcvbnma")
	buffer.WriteString(params.Encode())
	r, _ := http.NewRequest("POST", srv.URL+"/stat", buffer)
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	cli := srv.Client()

	resp, err := cli.Do(r)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected statuscode %d, got statuscode %d", http.StatusOK, resp.StatusCode)
	}
	resp.Body.Close()
	buffer = new(bytes.Buffer)
	params = url.Values{}
	params.Set("admin_url", "")
	buffer.WriteString(params.Encode())
	r, _ = http.NewRequest("POST", srv.URL+"/stat", buffer)
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	cli = srv.Client()

	resp, err = cli.Do(r)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected statuscode %d, got statuscode %d", http.StatusBadRequest, resp.StatusCode)
	}
	resp.Body.Close()
	buffer = new(bytes.Buffer)
	params = url.Values{}
	params.Set("admin_url", "fasdfies")
	buffer.WriteString(params.Encode())
	r, _ = http.NewRequest("POST", srv.URL+"/stat", buffer)
	r.Header.Set("content-type", "application/x-www-form-urlencoded")

	cli = srv.Client()

	resp, err = cli.Do(r)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected statuscode %d, got statuscode %d", http.StatusBadRequest, resp.StatusCode)
	}
	resp.Body.Close()
}

func TestGetSmallUrl(t *testing.T) {
	mdb := mockdb.NewMockdb("", zap.L())
	srs := smurlrepo.NewSmurlStorage(mdb, zap.L())
	hs := handler.NewHandlers(srs, zap.L())
	roa := NewRouterOpenAPI(hs, zap.L(), "")
	srv := httptest.NewServer(roa)
	r, _ := http.NewRequest("GET", srv.URL+"/abcdefgh", nil)
	cli := srv.Client()

	resp, err := cli.Do(r)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected statuscode %d, got statuscode %d", http.StatusOK, resp.StatusCode)
	}

	resp.Body.Close()

	r, _ = http.NewRequest("GET", srv.URL+"/ ", nil)
	resp, err = cli.Do(r)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected statuscode %d, got statuscode %d", http.StatusBadRequest, resp.StatusCode)
	}
	resp.Body.Close()
	r, _ = http.NewRequest("GET", srv.URL+"/fsdfasdf", nil)
	resp, err = cli.Do(r)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected statuscode %d, got statuscode %d", http.StatusBadRequest, resp.StatusCode)
	}
	resp.Body.Close()
}

func TestGet(t *testing.T) {
	mdb := mockdb.NewMockdb("", zap.L())
	srs := smurlrepo.NewSmurlStorage(mdb, zap.L())
	hs := handler.NewHandlers(srs, zap.L())
	roa := NewRouterOpenAPI(hs, zap.L(), "")
	srv := httptest.NewServer(roa)

	r, _ := http.NewRequest("GET", srv.URL+"/", nil)
	cli := srv.Client()

	resp, err := cli.Do(r)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected statuscode %d, got statuscode %d", http.StatusOK, resp.StatusCode)
	}
	resp.Body.Close()
}

package delivery

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	helpers "github.com/sanyarise/smurl/internal/helpers/mocks"
	"github.com/sanyarise/smurl/internal/models"
	"github.com/sanyarise/smurl/internal/usecase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	ctx       = context.Background()
	testLong  = "http://vk.com"
	err       = errors.New("error")
	testSmurl = &models.Smurl{
		SmallURL: "test",
		AdminURL: "test",
	}
	testSmurlUpd = &models.Smurl{
		SmallURL: "test",
		AdminURL: "test",
		Count: 1,
		IPInfo: []string{"testIpInfo"},
	}
)

type TestStatement struct {
	helpers *helpers.MockHelper
	logger  *zap.Logger
	usecase *mocks.MockUsecase
	router  *Router
}

func NewTestStatement(ctrl *gomock.Controller) *TestStatement {
	helpers := helpers.NewMockHelper(ctrl)
	usecase := mocks.NewMockUsecase(ctrl)
	logger := zap.L()
	router := NewRouter(usecase, helpers, logger, "testUrl")
	return &TestStatement{
		helpers: helpers,
		logger:  logger,
		usecase: usecase,
		router:  router,
	}
}

func GetRequest(longUrl string, serverUrl string, method string) *http.Request {
	params := url.Values{}
	params.Set("long_url", longUrl)
	buffer := new(bytes.Buffer)
	buffer.WriteString(params.Encode())
	r, _ := http.NewRequest("POST", serverUrl+"/create", buffer)
	r.Header.Set("content-type", "application/x-www-form-urlencoded")
	return r
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := NewTestStatement(ctrl)
	server := httptest.NewServer(s.router)

	r := GetRequest(testLong, server.URL, "POST")
	client := server.Client()
	s.helpers.EXPECT().CheckURL(testLong).Return(false)
	resp, err := client.Do(r)
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, 400, resp.StatusCode)
	resp.Body.Close()
}

func TestCreate2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := NewTestStatement(ctrl)
	server := httptest.NewServer(s.router)

	r := GetRequest(testLong, server.URL, "POST")
	client := server.Client()
	s.helpers.EXPECT().CheckURL(testLong).Return(true)
	s.usecase.EXPECT().Create(ctx, testLong).Return(nil, err)
	resp, err := client.Do(r)
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, 500, resp.StatusCode)
}

func TestCreate3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := NewTestStatement(ctrl)
	server := httptest.NewServer(s.router)

	r := GetRequest(testLong, server.URL, "POST")
	client := server.Client()
	s.helpers.EXPECT().CheckURL(testLong).Return(true)
	s.usecase.EXPECT().Create(ctx, testLong).Return(testSmurl, nil)
	resp, err := client.Do(r)
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, 201, resp.StatusCode)
}

func TestRedirect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := NewTestStatement(ctrl)
	server := httptest.NewServer(s.router)

	r, _ := http.NewRequest("GET", server.URL+"/r/testSmallUrl", nil)
	client := server.Client()
	s.usecase.EXPECT().FindURL(ctx, "testSmallUrl").Return(nil, models.ErrNotFound)
	resp, err := client.Do(r)
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, 400, resp.StatusCode)
	resp.Body.Close()
}

func TestRedirect2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := NewTestStatement(ctrl)
	server := httptest.NewServer(s.router)

	r, _ := http.NewRequest("GET", server.URL+"/r/testSmallUrl", nil)
	client := server.Client()
	s.usecase.EXPECT().FindURL(ctx, "testSmallUrl").Return(nil, err)
	resp, err := client.Do(r)
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, 500, resp.StatusCode)
	resp.Body.Close()
}

/*func TestRedirect3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := NewTestStatement(ctrl)
	server := httptest.NewServer(s.router)

	r, _ := http.NewRequest("GET", server.URL+"/r/testSmallUrl", nil)
	client := server.Client()
	s.usecase.EXPECT().FindURL(ctx, "testSmallUrl").Return(testSmurl, nil)
	s.helpers.EXPECT().GetIP(r.).Return("testIpInfo")
	s.usecase.EXPECT().UpdateStat(ctx, testSmurlUpd).Return(err)
	resp, err := client.Do(r)
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, 500, resp.StatusCode)
	resp.Body.Close()

}*/
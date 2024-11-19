package endpoints

import (
	"bytes"
	"context"
	internalmock "emailn/internal/test/internal-mock"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
)

var (
	service *internalmock.CampaignServiceMock
	handler = Handler{}
)

func setup() {
	service = new(internalmock.CampaignServiceMock)
	handler.CampaignService = service
}

func newHttpTest(requestMethod string, path string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {

	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest(requestMethod, path, &buf)
	rr := httptest.NewRecorder()

	return req, rr
}

func addParameter(req *http.Request, keyParameter string, value string) *http.Request {
	chiContext := chi.NewRouteContext()
	chiContext.URLParams.Add(keyParameter, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiContext))
}

func addContext(req *http.Request, keyParameter string, value string) *http.Request {
	ctx := context.WithValue(req.Context(), keyParameter, value)
	return req.WithContext(ctx)
}

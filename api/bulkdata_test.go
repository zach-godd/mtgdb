package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"mtgdb/data"
)

// fakeDataRetriever is a DataRetriever backed by test-supplied funcs.
type fakeDataRetriever struct {
	fetch    func() (*data.BulkDataList, error)
	download func(data.BulkDataItem, string) error
}

func (f fakeDataRetriever) FetchBulkData() (*data.BulkDataList, error) {
	return f.fetch()
}

func (f fakeDataRetriever) DownloadBulkData(item data.BulkDataItem, destPath string) error {
	return f.download(item, destPath)
}

// newTestAPI builds an API wired to a fakeDataRetriever and redirects
// bulkDataDir to a temp directory for the duration of the test.
func newTestAPI(t *testing.T, fetch func() (*data.BulkDataList, error), download func(data.BulkDataItem, string) error) API {
	t.Helper()

	origDir := bulkDataDir
	bulkDataDir = t.TempDir()
	t.Cleanup(func() {
		bulkDataDir = origDir
	})

	return NewAPI(fakeDataRetriever{fetch: fetch, download: download})
}

func TestHandleFetchBulkData_MethodNotAllowed(t *testing.T) {
	a := newTestAPI(t,
		func() (*data.BulkDataList, error) {
			t.Fatal("FetchBulkData should not be called for disallowed method")
			return nil, nil
		},
		func(data.BulkDataItem, string) error {
			t.Fatal("DownloadBulkData should not be called for disallowed method")
			return nil
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/bulk-data", nil)
	rec := httptest.NewRecorder()

	a.HandleFetchBulkData(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestHandleFetchBulkData_FetchError(t *testing.T) {
	a := newTestAPI(t,
		func() (*data.BulkDataList, error) {
			return nil, errors.New("boom")
		},
		func(data.BulkDataItem, string) error {
			t.Fatal("DownloadBulkData should not be called when fetch fails")
			return nil
		},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/bulk-data", nil)
	rec := httptest.NewRecorder()

	a.HandleFetchBulkData(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, rec.Code)
	}
}

func TestHandleFetchBulkData_TypeNotFound(t *testing.T) {
	a := newTestAPI(t,
		func() (*data.BulkDataList, error) {
			return &data.BulkDataList{Data: []data.BulkDataItem{{Type: "rulings"}}}, nil
		},
		func(data.BulkDataItem, string) error {
			t.Fatal("DownloadBulkData should not be called when type is not found")
			return nil
		},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/bulk-data?type=default_cards", nil)
	rec := httptest.NewRecorder()

	a.HandleFetchBulkData(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHandleFetchBulkData_DownloadError(t *testing.T) {
	a := newTestAPI(t,
		func() (*data.BulkDataList, error) {
			return &data.BulkDataList{Data: []data.BulkDataItem{{Type: defaultBulkDataType}}}, nil
		},
		func(data.BulkDataItem, string) error {
			return errors.New("network down")
		},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/bulk-data", nil)
	rec := httptest.NewRecorder()

	a.HandleFetchBulkData(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, rec.Code)
	}
}

func TestHandleFetchBulkData_Success(t *testing.T) {
	var downloadedItem data.BulkDataItem
	var downloadedPath string

	a := newTestAPI(t,
		func() (*data.BulkDataList, error) {
			return &data.BulkDataList{Data: []data.BulkDataItem{
				{Type: "rulings", DownloadURI: "https://example.com/rulings.json"},
				{Type: defaultBulkDataType, DownloadURI: "https://example.com/default-cards.json"},
			}}, nil
		},
		func(item data.BulkDataItem, destPath string) error {
			downloadedItem = item
			downloadedPath = destPath
			return os.WriteFile(destPath, []byte("[]"), 0o644)
		},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/bulk-data", nil)
	rec := httptest.NewRecorder()

	a.HandleFetchBulkData(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if downloadedItem.Type != defaultBulkDataType {
		t.Errorf("expected downloaded item type %q, got %q", defaultBulkDataType, downloadedItem.Type)
	}
	wantPath := filepath.Join(bulkDataDir, defaultBulkDataType+".json")
	if downloadedPath != wantPath {
		t.Errorf("expected dest path %q, got %q", wantPath, downloadedPath)
	}

	gotContentType := rec.Header().Get("Content-Type")
	if gotContentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", gotContentType)
	}

	if _, err := os.Stat(wantPath); err != nil {
		t.Errorf("expected bulk data directory to contain downloaded file: %v", err)
	}
}
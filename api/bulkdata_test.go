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

// withFakes swaps fetchBulkData/downloadBulkData/bulkDataDir for the
// duration of a test and restores the originals afterwards.
func withFakes(t *testing.T, fetch func() (*data.BulkDataList, error), download func(data.BulkDataItem, string) error) {
	t.Helper()

	origFetch, origDownload, origDir := fetchBulkData, downloadBulkData, bulkDataDir
	fetchBulkData = fetch
	downloadBulkData = download
	bulkDataDir = t.TempDir()

	t.Cleanup(func() {
		fetchBulkData = origFetch
		downloadBulkData = origDownload
		bulkDataDir = origDir
	})
}

func TestHandleFetchBulkData_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/bulk-data", nil)
	rec := httptest.NewRecorder()

	HandleFetchBulkData(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestHandleFetchBulkData_FetchError(t *testing.T) {
	withFakes(t,
		func() (*data.BulkDataList, error) {
			return nil, errors.New("boom")
		},
		func(data.BulkDataItem, string) error {
			t.Fatal("downloadBulkData should not be called when fetch fails")
			return nil
		},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/bulk-data", nil)
	rec := httptest.NewRecorder()

	HandleFetchBulkData(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, rec.Code)
	}
}

func TestHandleFetchBulkData_TypeNotFound(t *testing.T) {
	withFakes(t,
		func() (*data.BulkDataList, error) {
			return &data.BulkDataList{Data: []data.BulkDataItem{{Type: "rulings"}}}, nil
		},
		func(data.BulkDataItem, string) error {
			t.Fatal("downloadBulkData should not be called when type is not found")
			return nil
		},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/bulk-data?type=default_cards", nil)
	rec := httptest.NewRecorder()

	HandleFetchBulkData(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHandleFetchBulkData_DownloadError(t *testing.T) {
	withFakes(t,
		func() (*data.BulkDataList, error) {
			return &data.BulkDataList{Data: []data.BulkDataItem{{Type: defaultBulkDataType}}}, nil
		},
		func(data.BulkDataItem, string) error {
			return errors.New("network down")
		},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/bulk-data", nil)
	rec := httptest.NewRecorder()

	HandleFetchBulkData(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, rec.Code)
	}
}

func TestHandleFetchBulkData_Success(t *testing.T) {
	var downloadedItem data.BulkDataItem
	var downloadedPath string

	withFakes(t,
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

	HandleFetchBulkData(rec, req)

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
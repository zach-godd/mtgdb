package data

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"mtgdb/types"
)

func TestFetchBulkData(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	list, err := FetchBulkData()
	if err != nil {
		t.Fatalf("FetchBulkData: %v", err)
	}
	if len(list.Data) == 0 {
		t.Fatal("expected at least one bulk data item")
	}
	for _, item := range list.Data {
		if item.ID == "" {
			t.Errorf("item %q missing ID", item.Type)
		}
		if item.DownloadURI == "" {
			t.Errorf("item %q missing DownloadURI", item.Type)
		}
	}
}

func TestDownloadBulkData(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cards := []types.Card{
		{ID: "test-1", Object: "card", Name: "Lightning Bolt"},
		{ID: "test-2", Object: "card", Name: "Counterspell"},
	}
	body, err := json.Marshal(cards)
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}))
	defer srv.Close()

	item := BulkDataItem{
		ID:          "test-item",
		Name:        "Test Bulk Data",
		DownloadURI: srv.URL + "/bulk.json",
		ContentType: "application/json",
	}

	destPath := filepath.Join(t.TempDir(), "bulk.json")
	if err := DownloadBulkData(item, destPath); err != nil {
		t.Fatalf("DownloadBulkData: %v", err)
	}

	info, err := os.Stat(destPath)
	if err != nil {
		t.Fatalf("stat dest file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("downloaded file is empty")
	}
}

func TestImportBulkData(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := Connect(mongoURI, "mtgdb_test")
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer func() {
		client.DB.Drop(context.Background())
		client.Close(context.Background())
	}()

	cards := []types.Card{
		{ID: "abc-001", Object: "card", Name: "Lightning Bolt"},
		{ID: "abc-002", Object: "card", Name: "Counterspell"},
		{ID: "abc-003", Object: "card", Name: "Dark Ritual"},
	}

	f, err := os.CreateTemp(t.TempDir(), "cards-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(cards); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()

	if err := client.ImportBulkData(ctx, f.Name(), "cards"); err != nil {
		t.Fatalf("ImportBulkData: %v", err)
	}

	coll := client.DB.Collection("cards")
	count, err := coll.CountDocuments(ctx, bson.D{})
	if err != nil {
		t.Fatalf("CountDocuments: %v", err)
	}
	if count != int64(len(cards)) {
		t.Errorf("expected %d documents, got %d", len(cards), count)
	}

	// re-importing must not duplicate documents (upsert semantics)
	if err := client.ImportBulkData(ctx, f.Name(), "cards"); err != nil {
		t.Fatalf("ImportBulkData re-run: %v", err)
	}
	count, err = coll.CountDocuments(ctx, bson.D{})
	if err != nil {
		t.Fatalf("CountDocuments after re-import: %v", err)
	}
	if count != int64(len(cards)) {
		t.Errorf("expected %d documents after re-import, got %d", len(cards), count)
	}
}

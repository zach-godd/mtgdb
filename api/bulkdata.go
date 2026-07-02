package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"mtgdb/data"
)

// bulkDataDir is where downloaded Scryfall bulk data files are stored.
const bulkDataDir = "data/bulk"

// defaultBulkDataType is used when the caller does not specify a type.
// "default_cards" has one entry per print, matching the fields in types.Card.
const defaultBulkDataType = "default_cards"

// HandleFetchBulkData fetches the Scryfall bulk-data listing, selects the
// item matching the "type" query parameter (default_cards if unset), and
// downloads it to a local file under data/bulk/.
func HandleFetchBulkData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bulkType := r.URL.Query().Get("type")
	if bulkType == "" {
		bulkType = defaultBulkDataType
	}

	list, err := data.FetchBulkData()
	if err != nil {
		http.Error(w, fmt.Sprintf("fetching bulk data listing: %v", err), http.StatusBadGateway)
		return
	}

	var item *data.BulkDataItem
	for i := range list.Data {
		if list.Data[i].Type == bulkType {
			item = &list.Data[i]
			break
		}
	}
	if item == nil {
		http.Error(w, fmt.Sprintf("no bulk data item found for type %q", bulkType), http.StatusNotFound)
		return
	}

	if err := os.MkdirAll(bulkDataDir, 0o755); err != nil {
		http.Error(w, fmt.Sprintf("creating bulk data directory: %v", err), http.StatusInternalServerError)
		return
	}
	destPath := filepath.Join(bulkDataDir, item.Type+".json")

	if err := data.DownloadBulkData(*item, destPath); err != nil {
		http.Error(w, fmt.Sprintf("downloading bulk data: %v", err), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"type": item.Type,
		"path": destPath,
	})
}

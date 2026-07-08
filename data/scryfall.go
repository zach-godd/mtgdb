package data

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"mtgdb/types"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// importBatchSize number of documents to upload at once
const importBatchSize = 500

type BulkDataItem struct {
	Object          string    `json:"object"`
	ID              string    `json:"id"`
	Type            string    `json:"type"`
	UpdatedAt       time.Time `json:"updated_at"`
	URI             string    `json:"uri"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	CompressedSize  int64     `json:"compressed_size"`
	DownloadURI     string    `json:"download_uri"`
	ContentType     string    `json:"content_type"`
	ContentEncoding string    `json:"content_encoding"`
}

type BulkDataList struct {
	Object  string         `json:"object"`
	HasMore bool           `json:"has_more"`
	Data    []BulkDataItem `json:"data"`
}

// ScryfallClient is a DataRetriever backed by the real Scryfall API.
type ScryfallClient struct{}

func (ScryfallClient) FetchBulkData() (*BulkDataList, error) {
	return FetchBulkData()
}

func (ScryfallClient) DownloadBulkData(item BulkDataItem, destPath string) error {
	return DownloadBulkData(item, destPath)
}

// FetchBulkData fetches bulk data from sryfall.com
// Note: See https://scryfall.com/docs/api/bulk-data/all for Scryfall documentation
func FetchBulkData() (*BulkDataList, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.scryfall.com/bulk-data", nil)
	if err != nil {
		return nil, fmt.Errorf("building bulk data request: %w", err)
	}
	req.Header.Set("User-Agent", "mtgdb/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching bulk data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var result BulkDataList
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding bulk data: %w", err)
	}
	return &result, nil
}

// DownloadBulkData streams the bulk data file to destPath.
// The file is written as-is (compressed if the item uses gzip encoding).
func DownloadBulkData(item BulkDataItem, destPath string) error {
	req, err := http.NewRequest(http.MethodGet, item.DownloadURI, nil)
	if err != nil {
		return fmt.Errorf("building bulk data download request: %w", err)
	}
	req.Header.Set("User-Agent", "mtgdb/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("downloading bulk data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("writing bulk data: %w", err)
	}
	return nil
}

// ImportBulkData reads a Scryfall bulk data JSON file (plain or gzip-compressed)
// and upserts all cards into the named MongoDB collection in batches.
func (c *Client) ImportBulkData(ctx context.Context, filePath, collectionName string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	var r io.Reader = f
	if strings.HasSuffix(filePath, ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			return fmt.Errorf("creating gzip reader: %w", err)
		}
		defer gz.Close()
		r = gz
	}

	dec := json.NewDecoder(r)

	// consume opening '['
	if _, err := dec.Token(); err != nil {
		return fmt.Errorf("reading array start: %w", err)
	}

	coll := c.DB.Collection(collectionName)
	var models []mongo.WriteModel
	var total int

	for dec.More() {
		var card types.Card
		if err := dec.Decode(&card); err != nil {
			return fmt.Errorf("decoding card: %w", err)
		}
		models = append(models, mongo.NewReplaceOneModel().
			SetFilter(bson.D{{Key: "_id", Value: card.ID}}).
			SetReplacement(card).
			SetUpsert(true))

		if len(models) >= importBatchSize {
			if _, err := coll.BulkWrite(ctx, models); err != nil {
				return fmt.Errorf("bulk write: %w", err)
			}
			total += len(models)
			models = models[:0]
		}
	}

	if len(models) > 0 {
		if _, err := coll.BulkWrite(ctx, models); err != nil {
			return fmt.Errorf("bulk write: %w", err)
		}
		total += len(models)
	}

	_ = total
	return nil
}

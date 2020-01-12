package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
)

type GcsEvent struct {
	Bucket         string    `json:"bucket"`
	Name           string    `json:"name"`
	Metageneration string    `json:"metageneration"`
	ResourceState  string    `json:"resourceState"`
	TimeCreated    time.Time `json:"timeCreated"`
	Updated        time.Time `json:"updated"`
}

type NotificationEvent struct {
	Name      string `json:"email"`
	Event     string `json:"event"`
	Timestamp int64  `json:"timestamp"`
}

// Entry Point
func GcfEventHandler(ctx context.Context, e GcsEvent) error {
	// todo: metadata を使ってなにか処理をしたい場合には有効にする
	// _, err := metadata.FromContext(ctx)
	// if err != nil {
	// 	return fmt.Errorf("metadata.FromContext: %v", err)
	// }

	n, _ := getObjectFromGcs(ctx, e.Bucket, e.Name)
	writeToBigQuery(ctx, n)

	return nil
}

// GCS からオブジェクトを取得して JSON を構造体に合わせてパースする
func getObjectFromGcs(ctx context.Context, bucket string, object string) ([]NotificationEvent, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to initialize storage client : %v", err)
		return nil, err
	}

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		log.Printf("Failed to get object : %v", err)
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Printf("Failed to read object : %v", err)
		return nil, err
	}

	var n []NotificationEvent
	err = json.Unmarshal(data, &n)
	if err != nil {
		log.Printf("Failed to JSON parse : %v", err)
		return nil, err
	}

	return n, nil
}

// BigQuery に書き込む
func writeToBigQuery(ctx context.Context, n []NotificationEvent) error {
	ctx = context.Background()
	client, err := bigquery.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Printf("Failed to initialize bigquery client : %v", err)
		return err
	}
	defer client.Close()

	u := client.Dataset(os.Getenv("BIGQUERY_DATASET_NAME")).Table(os.Getenv("BIGQUERY_TABLE_NAME")).Uploader()
	// todo: range で回して突っ込むした方法がないのかな...
	for _, d := range n {
		item := []NotificationEvent{d}
		err = u.Put(ctx, item)
		if err != nil {
			log.Printf("Failed to write data to bigquery : %v", err)
			return err
		}
	}

	return nil
}

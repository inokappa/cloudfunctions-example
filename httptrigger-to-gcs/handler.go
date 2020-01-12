package handler

import (
	"bytes"
	"cloud.google.com/go/storage"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type NotificationEvent struct {
	Email string `json:"name"`
}

// Entry Point
func GcfWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Request Body を読み取って, そのまま GCS bucket に書き込む
	body := new(bytes.Buffer)
	body.ReadFrom(r.Body)
	strBody := body.String() + "\n"
	err := writeStorage(r, strBody)
	if err != nil {
		log.Printf("Failed to webhook request : %v", err)
		fmt.Fprintf(w, "error")
	} else {
		fmt.Fprintf(w, "ok")
	}
}

// GCS に保存するオブジェクト名を生成する
func genObjectName() string {
	unixTimeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	u, _ := uuid.NewRandom()
	uuid := u.String()
	objectName := unixTimeStamp + "-" + uuid + ".json"
	return objectName
}

func IsJson(str string) bool {
	var n []NotificationEvent
	err := json.Unmarshal([]byte(str), &n)
	if err != nil {
		log.Printf("Failed to JSON parse : %v", err)
		return false
	}
	return true
}

// GCS に POST された JSON をそのまま書き込む
// Refer to: https://blog.a-know.me/entry/2018/05/13/110707
// func writeStorage(r *http.Request, strBody string) error {
// 以下のような手続きオブジェクトとして書くことが出来る
var writeStorage = func(r *http.Request, strBody string) error {
	bucketName := os.Getenv("DEST_BUCKET")
	objectName := genObjectName()

	if !IsJson(strBody) {
		log.Printf("Request Body is not JSON : %s", strBody)
		return errors.New("Request Body is not JSON.")
	}

	ctx := r.Context()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create gcs client : %v", err)
		return err
	}

	writer := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	writer.ContentType = "text/plain"

	if _, err := writer.Write(([]byte)(strBody)); err != nil {
		log.Printf("Failed to write object body : %v", err)
		return err
	}

	if err := writer.Close(); err != nil {
		log.Printf("Failed to close gcs writer : %v", err)
		return err
	}
	return nil
}

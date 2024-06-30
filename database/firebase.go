package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

// InitFirebase initializes the Firebase Admin SDK
func InitFirebase() {
	// Get the private key from environment variable and format it correctly
	privateKey := os.Getenv("FIREBASE_PRIVATE_KEY")
	privateKey = strings.Replace(privateKey, "\\n", "\n", -1)

	// Construct the service account configuration using environment variables
	serviceAccount := map[string]string{
		"type":                        "service_account",
		"project_id":                  os.Getenv("FIREBASE_PROJECT_ID"),
		"private_key_id":              os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		"private_key":                 privateKey,
		"client_email":                os.Getenv("FIREBASE_CLIENT_EMAIL"),
		"client_id":                   os.Getenv("FIREBASE_CLIENT_ID"),
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        os.Getenv("FIREBASE_CLIENT_CERT_URL"),
	}

	// Convert the service account to JSON
	sa, err := json.Marshal(serviceAccount)
	if err != nil {
		log.Fatalf("error marshaling service account: %v", err)
	}

	// Initialize the Firebase app
	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(sa))
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}
	FirebaseApp = app
}

// GetFirebaseApp returns the Firebase app instance
func GetFirebaseApp() *firebase.App {
	return FirebaseApp
}

// GetFileFromBucket retrieves a file from Firebase Storage
func GetFileFromBucket(bucketName, filePath, destPath string) error {
	ctx := context.Background()
	client, err := FirebaseApp.Storage(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %v", err)
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket: %v", err)
	}

	object := bucket.Object(filePath)
	r, err := object.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create object reader: %v", err)
	}
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read object data: %v", err)
	}

	if err := ioutil.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}

// UploadFileToBucket uploads a file to Firebase Storage
func UploadFileToBucket(bucketName, srcPath, destPath string) error {
	ctx := context.Background()
	client, err := FirebaseApp.Storage(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %v", err)
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket: %v", err)
	}

	f, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	object := bucket.Object(destPath)
	w := object.NewWriter(ctx)
	if _, err := io.Copy(w, f); err != nil {
		return fmt.Errorf("failed to copy file to bucket: %v", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	return nil
}

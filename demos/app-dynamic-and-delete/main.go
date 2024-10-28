package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	secretValue []byte
	mu          sync.RWMutex // To ensure thread safety when accessing the secret value
)

func main() {
	// Path to the secret file
	secretFilePath := "/etc/secret/my-secret"

	// Start a background goroutine to refresh the secret every 10 seconds
	go func() {
		for {
			refreshSecret(secretFilePath)
			time.Sleep(10 * time.Second)
		}
	}()

	// Setup HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read the secret value safely
		mu.RLock()
		defer mu.RUnlock()

		// Display the secret value on the web page with some fun HTML styling
		fmt.Fprintf(w, `
            <!DOCTYPE html>
            <html lang="en">
            <head>
                <meta charset="UTF-8">
                <meta name="viewport" content="width=device-width, initial-scale=1.0">
                <title>Secret Value</title>
                <style>
                    body {
                        background-color: #2e3440;
                        font-family: Arial, sans-serif;
                        text-align: center;
                        padding-top: 50px;
                    }
                    .container {
                        display: inline-block;
                        background-color: #4c566a;
                        border-radius: 15px;
                        padding: 30px;
                        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
                    }
                    .secret-title {
                        color: #4682b4;
                        font-size: 24px;
                        font-weight: bold;
                    }
                    .secret-value {
                        color: #d2691e;
                        font-size: 32px;
                        margin-top: 10px;
                    }
                </style>
            </head>
            <body>
                <div class="container">
                    <div class="secret-title">Your Dynamic Secret Value is:</div>
                    <div class="secret-value">%q</div>
                </div>
            </body>
            </html>
        `, secretValue)
	})

	// Start the server
	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// refreshSecret reads the secret from the file, updates the in-memory value, and deletes the file content
func refreshSecret(secretFilePath string) {
	file, err := os.Open(secretFilePath)
	defer file.Close()
	if err != nil {
		log.Printf("Failed to open secret from file: %v Reading from env var", err)
		secretValue = []byte(getEnv("MY_SECRET", "Super Secret Default Secret"))
		return
	}

	r := bufio.NewReader(file)
	newSecretValue, _, err := r.ReadLine()
	if err != nil {
		log.Printf("Failed to open secret from file: %v Reading from env var", err)
		secretValue = []byte(getEnv("MY_SECRET", "Super Secret Default Secret"))
		return
	}

	//// Clear the file content after reading the secret
	//err = os.WriteFile(secretFilePath, []byte{}, 0644)
	//if err != nil {
	//	log.Printf("Failed to clear secret file: %v", err)
	//	return
	//}

	if !bytes.Equal(secretValue, newSecretValue) {
		log.Printf("NEW SECRET VALUE FOUND")
	}

	// Update the secret value in memory safely
	mu.Lock()
	secretValue = newSecretValue
	mu.Unlock()

	log.Printf("Secret value reloaded successfully.")
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

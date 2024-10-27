package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	//Read the secret value from a file (e.g., /etc/secret/my-secret)
	secretFilePath := "/etc/secret/my-secret"
	var secretValue []byte

	file, err := os.Open(secretFilePath)
	defer file.Close()
	if err == nil {
		r := bufio.NewReader(file)
		secretValue, _, err = r.ReadLine()
	}

	if err != nil {
		log.Printf("Failed to read secret from file: %v\nReading from env var", err)
		secretValue = []byte(getEnv("MY_SECRET", "Super Secret Default Secret"))
	}

	// Setup HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
                        background-color: #f0f8ff;
                        font-family: Arial, sans-serif;
                        text-align: center;
                        padding-top: 50px;
                    }
                    .container {
                        display: inline-block;
                        background-color: #ffebcd;
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
                    <div class="secret-title">Your Secret Value is:</div>
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

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

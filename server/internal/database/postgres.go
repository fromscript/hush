package database

import (
	"database/sql"
	"github.com/fromscript/hush/internal/crypto"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := "user=wisper_user dbname=wisper password=supersecret sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func SaveMessage(encryptedMsg []byte) {
	// Encrypt message with server master key before storage
	masterKey := []byte("demo_master_key") // Load from secure source in production
	doubleEncrypted, _ := crypto.Encrypt(encryptedMsg, masterKey)

	_, err := DB.Exec("INSERT INTO messages (content) VALUES ($1)", doubleEncrypted)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
	}
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Connection string - point to the SQLens proxy (5433)
	connStr := "postgres://user:password@localhost:5433/demo?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("🚀 Starting Realistic SQLens Benchmark...")

	// 1. Setup Table and Seed Data
	fmt.Println("\n--- 🛠️ Setup: Creating 'users' table ---")
	_, err = db.Exec(`
		DROP TABLE IF EXISTS users;
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("Inserting sample users...")
	for i := 1; i <= 20; i++ {
		_, err = db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", 
			fmt.Sprintf("User %d", i), 
			fmt.Sprintf("user%d@example.com", i))
		if err != nil {
			log.Printf("Insert error: %v", err)
		}
	}

	// 2. Realistic N+1 Simulation
	// Scenario: Fetching user IDs, then fetching full details for each one individually
	fmt.Println("\n--- 🕵️ Scenario: N+1 Detection (Fetching user details one by one) ---")
	rows, err := db.Query("SELECT id FROM users LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}
	
	var ids []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}
	rows.Close()

	for _, id := range ids {
		start := time.Now()
		var name, email string
		// SQLens will normalize this to 'SELECT * FROM users WHERE id = ?'
		err := db.QueryRow("SELECT name, email FROM users WHERE id = $1", id).Scan(&name, &email)
		if err != nil {
			log.Printf("Query error for ID %d: %v", id, err)
			continue
		}
		fmt.Printf("Fetched %s (latency: %v)\n", name, time.Since(start))
		time.Sleep(100 * time.Millisecond) // Slow it down to see it live
	}

	// 3. Slow Query Simulation
	fmt.Println("\n--- 🐢 Scenario: Slow Query Detection (Complex search) ---")
	start := time.Now()
	var count int
	// Using pg_sleep to force a slow query report in SQLens
	err = db.QueryRow("SELECT count(*), pg_sleep(1.5) FROM users WHERE email LIKE '%example%'").Scan(&count, new(string))
	if err != nil {
		log.Printf("Slow query error: %v", err)
	}
	fmt.Printf("Slow query finished (latency: %v)\n", time.Since(start))

	// 4. Batch Updates
	fmt.Println("\n--- ⚡ Scenario: Batch Updates ---")
	for i := 0; i < 50; i++ {
		_, _ = db.Exec("UPDATE users SET created_at = NOW() WHERE id = $1", (i%20)+1)
	}
	fmt.Println("Finished 50 quick updates.")

	fmt.Println("\n✅ Benchmark Complete!")
	fmt.Println("Check the SQLens Dashboard (http://localhost:8080) for N+1 alerts and latency maps.")
}

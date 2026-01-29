package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type Content struct {
	Content string `json:"content"`
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get Method Allowed", http.StatusMethodNotAllowed)
		return
	}

	var contents []Content
	query := `SELECT content FROM crud`

	rows, err := DB.Query(query)
	if err != nil {
		http.Error(w, "Invalid query to the database", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var content Content
		if err := rows.Scan(&content.Content); err != nil {
			http.Error(w, "Invalid database scan", http.StatusInternalServerError)
			return
		}

		contents = append(contents, content)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(contents); err != nil {
		http.Error(w, "Invalid response", http.StatusBadRequest)
		return
	}

}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post Method Allowed", http.StatusMethodNotAllowed)
		return
	}

	query := `INSERT INTO crud (content) VALUES (?)`

	var content Content
	err := json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = DB.Exec(query, content.Content)
	if err != nil {
		http.Error(w, "Invalid query to the database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Content created successfully"})
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only Delete Method Allowed", http.StatusMethodNotAllowed)
		return
	}
	query := `DELETE FROM crud WHERE content = ?`

	var content Content
	err := json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = DB.Exec(query, content.Content)
	if err != nil {
		http.Error(w, "Invalid query to the database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Content deleted successfully"})
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only Put Method Allowed", http.StatusMethodNotAllowed)
		return
	}
	query := `UPDATE crud SET content = ? WHERE content = ?`

	var contents struct {
		OldContent string `json:"old_content"`
		NewContent string `json:"new_content"`
	}
	err := json.NewDecoder(r.Body).Decode(&contents)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = DB.Exec(query, contents.NewContent, contents.OldContent)
	if err != nil {
		http.Error(w, "Invalid query to the database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Content updated successfully"})
}

func TotalNumberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get Method Allowed", http.StatusMethodNotAllowed)
		return
	}

	var total int
	query := `SELECT COUNT(*) FROM crud`

	err := DB.QueryRow(query).Scan(&total)
	if err != nil {
		http.Error(w, "Invalid query to the database", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := struct {
		Total int `json:"total"`
	}{
		Total: total,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Invalid response", http.StatusBadRequest)
		return
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h(w, r)
	}
}

func main() {
	var err error
	dsn := "root:29112003@tcp(localhost:3306)/crud"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer DB.Close()

	err = DB.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		return
	}
	fmt.Println("Connected to the database successfully.")

	http.HandleFunc("/", withCORS(Health))
	http.HandleFunc("/get", withCORS(GetHandler))
	http.HandleFunc("/post", withCORS(PostHandler))
	http.HandleFunc("/delete", withCORS(DeleteHandler))
	http.HandleFunc("/update", withCORS(UpdateHandler))
	http.HandleFunc("/totalnum", withCORS(TotalNumberHandler))

	fmt.Println("Server is running on port 3030")
	err = http.ListenAndServe(":3030", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}

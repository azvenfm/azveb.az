
package main
import (
    "fmt"
    "database/sql"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    _ "github.com/lib/pq"
)
func main() {
    db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/autosmm?sslmode=disable")
    if err != nil { panic(err) }
    defer db.Close()

    password := "Admin123!"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

    id := uuid.New()
    _, err = db.Exec("INSERT INTO users (id, username, email, password_hash, full_name, role, plan) VALUES ($1, $2, $3, $4, $5, $6, $7)",
        id, "admin", "admin@autosmm.az", string(hash), "AutoSMM Admin", "admin", "premium")
    
    if err != nil {
        fmt.Println("Error creating admin:", err)
    } else {
        fmt.Println("Admin account created successfully!")
    }
}

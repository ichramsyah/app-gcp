package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func registerHandler(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var userID int
	err = db.QueryRow(context.Background(),
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
		creds.Username, string(hashedPassword)).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username already exists or database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": userID})
}

func loginHandler(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var user User
	err := db.QueryRow(context.Background(),
		"SELECT id, username, password_hash FROM users WHERE username=$1", creds.Username).
		Scan(&user.ID, &user.Username, &user.PasswordHash)
		
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func getNotesHandler(c *gin.Context) {
	userID, _ := c.Get("userID")

	rows, err := db.Query(context.Background(),
		"SELECT id, title, content, created_at FROM notes WHERE user_id=$1 ORDER BY created_at DESC", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}
	defer rows.Close()

	var notes []Note
	// Initialize slice to empty rather than nil so that empty lists return [] instead of null
	notes = make([]Note, 0)
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing notes"})
			return
		}
		notes = append(notes, n)
	}

	c.JSON(http.StatusOK, notes)
}

func createNoteHandler(c *gin.Context) {
	userID, _ := c.Get("userID")

	var note Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	err := db.QueryRow(context.Background(),
		"INSERT INTO notes (user_id, title, content) VALUES ($1, $2, $3) RETURNING id, created_at",
		userID, note.Title, note.Content).Scan(&note.ID, &note.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}
	note.UserID = userID.(int)

	c.JSON(http.StatusCreated, note)
}

func updateNoteHandler(c *gin.Context) {
	userID, _ := c.Get("userID")
	noteID := c.Param("id")

	var note Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	tag, err := db.Exec(context.Background(),
		"UPDATE notes SET title=$1, content=$2 WHERE id=$3 AND user_id=$4",
		note.Title, note.Content, noteID, userID)

	if err != nil || tag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note updated"})
}

func deleteNoteHandler(c *gin.Context) {
	userID, _ := c.Get("userID")
	noteID := c.Param("id")

	tag, err := db.Exec(context.Background(),
		"DELETE FROM notes WHERE id=$1 AND user_id=$2", noteID, userID)

	if err != nil || tag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}

package server

import "github.com/gin-gonic/gin"

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) RegisterHandler(c *gin.Context) {
	var body RegisterRequest
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	email := body.Email
	password := body.Password

	createdUser, err := s.db.CreateUser(email, password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	c.JSON(201, gin.H{"message": "user created", "id": createdUser.Id})
}

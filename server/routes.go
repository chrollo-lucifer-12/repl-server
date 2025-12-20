package server

import "github.com/gin-gonic/gin"

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateProjectHandleRequest struct {
	Slug   string `json:"slug"`
	UserId uint   `json:"userId"`
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

func (s *Server) CreateProjectHandler(c *gin.Context) {
	var body CreateProjectHandleRequest
	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	slug := body.Slug
	userId := body.UserId

	createdProject, err := s.db.CreateProject(slug, userId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	c.JSON(201, gin.H{"message": "project created", "id": createdProject.Id})
}

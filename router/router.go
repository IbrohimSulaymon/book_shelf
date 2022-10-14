package router

import (
	"book_shelf/server"
	"github.com/gin-gonic/gin"
)

func InitRouter(s server.Server) *gin.Engine {
	r := gin.Default()

	r.POST("/signup", s.CreateUser)
	r.GET("/myself", s.GetUserInfo)
	r.POST("/books", s.CreateBook)
	r.GET("/books", s.GetAllBooks)
	r.PATCH("/books/:id", s.EditBook)
	r.DELETE("/books/:id", s.DeleteBook)

	return r
}

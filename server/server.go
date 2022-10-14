package server

import (
	"book_shelf/domain"
	"book_shelf/pkg"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	repo Repository
}

func New(repo Repository) Server {
	return Server{
		repo: repo,
	}
}

func (s Server) Authorise(c *gin.Context, signString string) error {
	key := c.GetHeader("key")
	sign := c.GetHeader("sign")
	secret, err := s.repo.Check(c, key)
	if err != nil {
		return err
	}

	mySign := md5.Sum([]byte(signString + secret))
	if fmt.Sprintf("%x", mySign) != sign {
		return fmt.Errorf("secret key is not valid")
	}

	return nil
}

func (s Server) CreateUser(c *gin.Context) {
	var request CreateNewUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := s.repo.CreateUser(c.Request.Context(), domain.User{
		Name:   request.Name,
		Key:    request.Key,
		Secret: request.Secret,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error")
		return
	}

	c.JSON(http.StatusOK, CreateNewUserResponse{
		Data: dataUser{
			ID:     id,
			Name:   request.Name,
			Key:    request.Key,
			Secret: request.Secret,
		},
		IsOk:    true,
		Message: "ok",
	})
}

func (s Server) GetUserInfo(c *gin.Context) {
	//if err := s.Authorise(c, pkg.GetUserInfoSign); err != nil {
	//	c.JSON(http.StatusBadRequest, GetUserInfoResponse{
	//		IsOk:    false,
	//		Message: err.Error(),
	//	})
	//
	//	return
	//}
	user, err := s.repo.GetUserInfo(c.Request.Context(), c.GetHeader("key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetUserInfoResponse{
			IsOk:    false,
			Message: fmt.Sprintf("%s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, GetUserInfoResponse{
		Data: dataUser{
			ID:     user.ID,
			Name:   user.Name,
			Key:    user.Key,
			Secret: user.Secret,
		},
		IsOk:    true,
		Message: "ok",
	})
}

func (s Server) CreateBook(c *gin.Context) {
	var (
		body        CreateNewBookRequest
		book        GetBookWithISBNResponse
		author      GetAuthorWithISBNResponse
		authorNames string
	)
	//if err := s.Authorise(c, pkg.CreateBookSign); err != nil {
	//	c.JSON(http.StatusBadRequest, CreateNewBookResponse{
	//		IsOk:    false,
	//		Message: err.Error(),
	//	})
	//
	//	return
	//}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, CreateNewBookResponse{
			IsOk:    false,
			Message: err.Error(),
		})

		return
	}
	url := pkg.GetBookRequestURL + body.ISBN + ".json"
	response, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, CreateNewBookResponse{
			IsOk:    false,
			Message: err.Error(),
		})

		return
	}

	if err := json.NewDecoder(response.Body).Decode(&book); err != nil {
		c.JSON(http.StatusBadRequest, CreateNewBookResponse{
			IsOk:    false,
			Message: err.Error(),
		})

		return
	}
	for i, val := range book.Author {
		url = pkg.GetAuthorRequestURL + val.Key + ".json"
		response, err = http.Get(url)
		if err != nil {
			c.JSON(http.StatusBadRequest, CreateNewBookResponse{
				IsOk:    false,
				Message: err.Error(),
			})

			return
		}

		if err := json.NewDecoder(response.Body).Decode(&author); err != nil {
			c.JSON(http.StatusBadRequest, CreateNewBookResponse{
				IsOk:    false,
				Message: err.Error(),
			})

			return
		}
		if i == len(book.Author)-1 {
			authorNames += author.Name
		} else {
			authorNames += author.Name + ", "
		}
	}

	pDate, err := strconv.Atoi(book.PublishDate[len(book.PublishDate)-4:])
	if err != nil {
		c.JSON(http.StatusInternalServerError, CreateNewBookResponse{
			IsOk:    false,
			Message: err.Error(),
		})

		return
	}
	id, err := s.repo.CreateBook(c.Request.Context(), domain.Book{
		ISBN:      body.ISBN,
		Title:     book.Title,
		Author:    authorNames,
		Published: pDate,
		Pages:     book.NumberOfPages,
		Status:    0,
	})
	if err != nil {
		if err != nil {
			c.JSON(http.StatusInternalServerError, CreateNewBookResponse{
				IsOk:    false,
				Message: err.Error(),
			})

			return
		}
	}

	c.JSON(http.StatusOK, CreateNewBookResponse{
		Data: createBookData{
			Book: dataBook{
				ID:        id,
				ISBN:      body.ISBN,
				Title:     book.Title,
				Author:    authorNames,
				Published: pDate,
				Page:      book.NumberOfPages,
			},
			Status: 0,
		},
		IsOk:    true,
		Message: "ok",
	})
}

func (s Server) EditBook(c *gin.Context) {
	body := EditBookRequest{}
	//if err := s.Authorise(c, pkg.EditBookSign); err != nil {
	//	c.JSON(http.StatusBadRequest, EditBookResponse{
	//		IsOk:    false,
	//		Message: err.Error(),
	//	})
	//
	//  return
	//}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, EditBookResponse{
			IsOk:    false,
			Message: err.Error(),
		})

		return
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, EditBookResponse{
			IsOk:    false,
			Message: err.Error(),
		})

		return
	}

	b, err := s.repo.EditBook(c, id, body.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, EditBookResponse{
			IsOk:    false,
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, EditBookResponse{
		Data: createBookData{
			Book: dataBook{
				ID:        b.ID,
				ISBN:      b.ISBN,
				Title:     b.Title,
				Author:    b.Author,
				Published: b.Published,
				Page:      b.Pages,
			},
			Status: b.Status,
		},
		IsOk:    true,
		Message: "ok",
	})
}

func (s Server) GetAllBooks(c *gin.Context) {
	res := GetAllBooksResponse{}
	//if err := s.Authorise(c, pkg.GetAllBooksSign); err != nil {
	//	c.JSON(http.StatusBadRequest, GetAllBooksResponse{
	//		IsOk:    false,
	//		Message: err.Error(),
	//	})
	//
	//	return
	//}

	books, err := s.repo.GetAllBooks(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetAllBooksResponse{
			IsOk:    false,
			Message: err.Error(),
		})

		return
	}

	for _, val := range books {
		res.Data = append(res.Data, createBookData{
			Book: dataBook{
				ID:        val.ID,
				ISBN:      val.ISBN,
				Title:     val.Title,
				Author:    val.Author,
				Published: val.Published,
				Page:      val.Pages,
			},
			Status: val.Status,
		})
	}
	res.IsOk = true
	res.Message = "ok"

	c.JSON(http.StatusOK, res)
}

func (s Server) DeleteBook(c *gin.Context) {
	//if err := s.Authorise(c, pkg.GetAllBooksSign); err != nil {
	//	c.JSON(http.StatusBadRequest, DeleteBookResponse{
	//		IsOk:    false,
	//		Message: err.Error(),
	//	})
	//
	//	return
	//}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, DeleteBookResponse{
			IsOk:    false,
			Message: err.Error(),
		})

		return
	}
	if err := s.repo.DeleteBook(c, id); err != nil {
		c.JSON(http.StatusBadRequest, DeleteBookResponse{
			IsOk:    false,
			Message: err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, DeleteBookResponse{
		Data:    "Successfully deleted",
		IsOk:    true,
		Message: "ok",
	})
}

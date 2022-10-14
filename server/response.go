package server

type dataUser struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type CreateNewUserResponse struct {
	Data    dataUser `json:"data"`
	IsOk    bool     `json:"is_ok"`
	Message string   `json:"message"`
}

type dataBook struct {
	ID        int    `json:"id"`
	ISBN      string `json:"isbn"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Published int    `json:"published"`
	Page      int    `json:"page"`
}

type CreateNewBookResponse struct {
	Data    createBookData `json:"data"`
	IsOk    bool           `json:"isOk"`
	Message string         `json:"message"`
}

type createBookData struct {
	Book   dataBook `json:"book"`
	Status int      `json:"status"`
}

type EditBookResponse struct {
	Data    createBookData `json:"data"`
	IsOk    bool           `json:"isOk"`
	Message string         `json:"message"`
}

type GetUserInfoResponse struct {
	Data    dataUser `json:"data"`
	IsOk    bool     `json:"isOk"`
	Message string   `json:"message"`
}

type GetBookWithISBNResponse struct {
	Title         string    `json:"title"`
	NumberOfPages int       `json:"number_of_pages"`
	PublishDate   string    `json:"publish_date"`
	Author        []authors `json:"authors"`
}

type GetAuthorWithISBNResponse struct {
	Name string `json:"personal_name"`
}

type authors struct {
	Key string `json:"key"`
}

type GetAllBooksResponse struct {
	Data    []createBookData `json:"data"`
	IsOk    bool             `json:"isOk"`
	Message string           `json:"message"`
}

type DeleteBookResponse struct {
	Data    string `json:"data"`
	IsOk    bool   `json:"isOk"`
	Message string `json:"message"`
}

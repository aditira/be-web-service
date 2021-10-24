package model

type Book struct {
	BookID   string `json:"bookid"`
	BookName string `json:"bookname"`
}

type JsonResponse struct {
	Type    string `json:"type"`
	Data    []Book `json:"data"`
	Message string `json:"message"`
}

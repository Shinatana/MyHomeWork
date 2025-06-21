package models

type GetUserResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type PostUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type CreateUserResponse struct {
	ID int64 `json:"id"`
}

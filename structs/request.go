package structs

type PostUserRequest struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Age      uint    `json:"age"`
}

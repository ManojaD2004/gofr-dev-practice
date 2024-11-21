package types

type UserRequestType struct {
	Name string `json:"name"`
}

type UserResponseType struct {
	Username string `json:"username"`
	RollNo int `json:"roll no"`
}

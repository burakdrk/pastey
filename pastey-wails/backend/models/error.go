package models

type Error struct {
	Message string `json:"error"`
}

func GetDefaultError() Error {
	return Error{Message: "There was an error. Please try again later."}
}

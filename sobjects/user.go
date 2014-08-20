package sobjects

type User struct {
	BaseSObject
	Email         string `force:",omitempty"`
	FirstName     string `force:",omitempty"`
	LastName      string `force:",omitempty"`
	SmallPhotoUrl string `force:",omitempty"`
	FullPhotoUrl  string `force:",omitempty"`
}

func (t *User) ApiName() string {
	return "User"
}

type UserQueryResponse struct {
	BaseQuery
	Records []User `json:"Records" force:"records"`
}

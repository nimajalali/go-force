package sobjects

type User struct {
	BaseSObject
	Alias             string `force:",omitempty"`
	CommunityNickname string `force:",omitempty"`
	Email             string `force:",omitempty"`
	EmailEncodingKey  string `force:",omitempty"`
	FirstName         string `force:",omitempty"`
	FullPhotoUrl      string `force:",omitempty"`
	LanguageLocaleKey string `force:",omitempty"`
	LastName          string `force:",omitempty"`
	LocaleSidKey      string `force:",omitempty"`
	ProfileId         string `force:",omitempty"`
	SmallPhotoUrl     string `force:",omitempty"`
	TimeZoneSidKey    string `force:",omitempty"`
	Username          string `force:",omitempty"`
}

func (t *User) ApiName() string {
	return "User"
}

type UserQueryResponse struct {
	BaseQuery
	Records []User `json:"Records" force:"records"`
}

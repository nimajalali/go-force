package sobjects

type Profile struct {
	BaseSObject
	Description               string `force:",omitempty"`
	IsSsoEnabled              bool   `force:",omitempty"`
	LastReferencedDate        string `force:",omitempty"`
	LastViewedDate            string `force:",omitempty"`
	Name                      string `force:",omitempty"`
	PermissionsPermissionName bool   `force:",omitempty"`
	UserLicenseId             string `force:",omitempty"`
	UserType                  string `force:",omitempty"`
}

func (t *Profile) APIName() string {
	return "Profile"
}

type ProfileQueryResponse struct {
	BaseQuery
	Records []Profile `json:"Records" force:"records"`
}

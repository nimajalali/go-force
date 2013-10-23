package sobjects

type Lead struct {
	BaseSObject
}

func (l *Lead) ApiName() string {
	return "Lead"
}

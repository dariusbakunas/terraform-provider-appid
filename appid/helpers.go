package appid

func getBoolPtr(b bool) *bool {
	val := b
	return &val
}

func getStringPtr(s string) *string {
	val := s
	return &val
}

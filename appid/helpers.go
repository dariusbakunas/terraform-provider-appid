package appid

func getBoolPtr(b bool) *bool {
	val := b
	return &val
}

func getStringPtr(s string) *string {
	val := s
	return &val
}

func flattenStringList(list []string) []interface{} {
	vs := make([]interface{}, len(list))
	for i, v := range list {
		vs[i] = v
	}
	return vs
}

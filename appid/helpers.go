package appid

import "encoding/json"

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

func expandStringList(list []interface{}) []string {
	vs := make([]string, 0, len(list))
	for _, v := range list {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
}

func dbgPrint(data interface{}) string {
	dataJSON, _ := json.MarshalIndent(data, "", "  ")
	return string(dataJSON)
}

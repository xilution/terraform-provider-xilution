package provider

import "strings"

func getIdFromLocationUrl(location *string) *string {
	index := strings.LastIndex(*location, "/")
	id := string((*location)[(index + 1):])

	return &id
}

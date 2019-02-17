package main

import (
	"strconv"
)

func addToRegistry(value string) string {
	if (value == "") {
		return ""
	}

	index, exists := registry[value]

	if exists {
		return strconv.Itoa(index)
	} else {
		index = len(registry)
		registry[value] = index
		return strconv.Itoa(index)
	}
}

func getFromRegistry(value string) string {
	if value == "" {
		return ""
	}

	z, _ := strconv.Atoi(value)

	for k, v := range registry {
		if  z == v {
			return k
		}
	}

	return ""
}
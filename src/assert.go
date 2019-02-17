package main

import (
	"strings"
)

func validateData(fields map[string]interface{}, account Account) bool {
	for k, _ := range fields {
		switch k {
		case "joined":
			if !assertDate(account.Joined) {
				return false
			}
		case "birth":
			if !assertDate(account.Birth) {
				return false
			}
		case "email":
			if (!assertEmail(account.Email)) {
				return false
			}
		case "phone":
			if !assertPhone(account.Phone) {
				return false
			}
		case "sex":
			if !assertSex(account.Sex) {
				return false
			}
		case "status":
			if !assertStatus(account.Status) {
				return false
			}
		case "premium":
			if account.Premium.Start <= 0 || account.Premium.Finish <= 0 {
				return false
			}
		case "likes":
			if (!assertLikes(account.Likes)) {
				return false
			}
		}
	}

	return true
}

func assertDate(date int) bool {
	if date > 0 {
		return true
	}

	return false
}

func assertEmail(email string) bool {
	if strings.Index(email, "@") == -1 {
		return false
	}

	_, exists := emails[email]

	if exists {
		return false
	}

	return true
}

func assertPhone(phone string) bool {
	_, exists := phones[phone]

	if exists {
		return false
	}

	return true
}

func assertSex(sex string) bool {
	if sex == "m" {
		return true
	}

	if sex == "f" {
		return true
	}

	return false
}

func assertStatus(status string) bool {
	if status == "свободны" {
		return true
	}

	if status == "заняты" {
		return true
	}

	if status == "всё сложно" {
		return true
	}

	return false
}

func assertLikes(likes []Like) bool {
	for _, like := range likes {
		_, exists := accountsData[like.Id]

		if !exists {
			return false
		}

		if like.Ts <= 0 {
			return false
		}
	}

	return true
}
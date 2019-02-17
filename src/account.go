package main

type Accounts struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	Id int `json:"id"`
	Interests []string `json:"interests"`
	Fname string `json:"fname"`
	Sname string `json:"sname"`
	Status string `json:"status"`
	Premium Premium `json:"premium"`
	Email string `json:"email"`
	Sex string `json:"sex"`
	Phone string `json:"phone"`
	Birth int `json:"birth"`
	Joined int `json:"joined"`
	City string `json:"city"`
	Country string `json:"country"`
	Likes []Like `json:"likes"`
}

type Premium struct {
	Start int `json:"start"`
	Finish int `json:"finish"`
}

type Like struct {
	Id int `json:"id"`
	Ts int `json:"ts"`
}
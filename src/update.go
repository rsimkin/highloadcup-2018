package main

import (
	"fmt"
	"time"
	"github.com/valyala/fasthttp"
	"strings"
	"encoding/json"
	"strconv"
)

func handleUpdate(ctx *fasthttp.RequestCtx) {
	parts := strings.Split(string(ctx.Path()), "/")

	if len(parts) < 3 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(parts[2])
	_, exists := accountsData[id]

	if !exists || id <= 0 {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	var account Account
	tmp := make(map[string]interface{})
	json.Unmarshal(ctx.PostBody(), &account)
	json.Unmarshal(ctx.PostBody(), &tmp)

	if !validateData(tmp, account) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	//Удалить из всех индексов фильтра
	for fieldName, indexInfo := range indexes {
		if fieldName != "interests" {
			for _, set := range indexInfo {
				_, exists = set.data[uint32(id)]

				if exists {
					//fmt.Println("delete from", valueName)
					delete(set.data, uint32(id))
				}
			}
		}
	}

	//Получить текущий аккаунт
	currentAccount := getAccountObjectFromData(id)

	birth := strconv.Itoa(time.Unix(int64(currentAccount.Birth), 0).Year())
	joined := strconv.Itoa(time.Unix(int64(currentAccount.Joined), 0).Year())
	addAccountToGroupIndex("_", "_", currentAccount, -1)
	addAccountToGroupIndex("fname", currentAccount.Fname, currentAccount, -1)
	addAccountToGroupIndex("sname", currentAccount.Sname, currentAccount, -1)
	addAccountToGroupIndex("city", currentAccount.City, currentAccount, -1)
	addAccountToGroupIndex("country", currentAccount.Country, currentAccount, -1)
	addAccountToGroupIndex("status", currentAccount.Status, currentAccount, -1)
	addAccountToGroupIndex("sex", currentAccount.Sex, currentAccount, -1)
	addAccountToGroupIndex("birth", birth, currentAccount, -1)
	addAccountToGroupIndex("joined", joined, currentAccount, -1)
	
	addAccountToGroupIndex("birth|city", birth + currentAccount.City, currentAccount, -1)
	addAccountToGroupIndex("birth|status", birth + currentAccount.Status, currentAccount, -1)
	addAccountToGroupIndex("joined|sex", joined + currentAccount.Sex, currentAccount, -1)
	addAccountToGroupIndex("joined|status", joined + currentAccount.Status, currentAccount, -1)
	addAccountToGroupIndex("birth|sex", birth + currentAccount.Sex, currentAccount, -1)
	addAccountToGroupIndex("city|joined", currentAccount.City + joined, currentAccount, -1)
	addAccountToGroupIndex("country|joined", currentAccount.Country + joined, currentAccount, -1)
	addAccountToGroupIndex("birth|country", birth + currentAccount.Country, currentAccount, -1)

	/*if currentAccount.Email != "" {
		delete(emails, currentAccount.Email)
	}*/

	/*if (currentAccount.Phone != "") {
		delete(emails, currentAccount.phones)
	}*/

	json.Unmarshal(ctx.PostBody(), &currentAccount)
	//fmt.Println(currentAccount)
	//fmt.Println(currentAccount)
	processAccount(currentAccount)

	//TODO
	/*
	selectFields := []string{"fname", "sname", "status", "email", "country", "birth", "city"}
	currentAccount := getAccountFromData(id, selectFields)

	for field, value := range tmp {
		switch field {
		case "fname", "sname", "status", "email", "phone":
			//Удалить из индекса
			deleteIdFromIndex(field, currentAccount[field].(string), id)
			addIdToIndex(field, value.(string), id)
			currentAccount[field] = value.(string)
		}
	}*/

	//updateAccountInData

	ctx.SetStatusCode(fasthttp.StatusAccepted)
	fmt.Fprintf(ctx, "{}")
}
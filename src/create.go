package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

func handleCreate(ctx *fasthttp.RequestCtx) {
	var account Account
	tmp := make(map[string]interface{})
	json.Unmarshal(ctx.PostBody(), &account)
	json.Unmarshal(ctx.PostBody(), &tmp)

	if !validateData(tmp, account) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	requiredFields := []string{"email", "sex", "birth", "joined", "status"}

	for _, field := range requiredFields {
		_, exists := tmp[field]

		if !exists {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	fmt.Fprintf(ctx, "{}")
	processAccount(account)
}
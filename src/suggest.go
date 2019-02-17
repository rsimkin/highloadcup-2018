package main

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/valyala/fasthttp"
)

func handleSuggest(ctx *fasthttp.RequestCtx) {
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

	fmt.Fprintf(ctx, "{\"accounts\":[]}")
}
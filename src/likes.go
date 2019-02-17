package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

func handleLikes(ctx *fasthttp.RequestCtx) {

	var likes LikesExt
	json.Unmarshal(ctx.PostBody(), &likes)

	for _, like := range likes.Likes {
		if like.Ts <= 0 {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		_, exists := accountsData[like.Liker]

		if !exists || like.Liker <= 0 {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		_, exists = accountsData[like.Likee]

		if !exists || like.Likee <= 0{
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}

	//TODO
	ctx.SetStatusCode(fasthttp.StatusAccepted)
	fmt.Fprintf(ctx, "{}")
}

type LikesExt struct {
	Likes []LikeExt `json:"likes"`
}

type LikeExt struct {
	Ts int `json:"Ts"`
	Liker int `json:"liker"`
	Likee int `json:"likee"`
}
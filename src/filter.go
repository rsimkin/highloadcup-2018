package main

import (
	"fmt"
	"strconv"
	"strings"
	"encoding/json"
	"sort"
	"net/url"
    "github.com/valyala/fasthttp"
)

func handleFilter(ctx *fasthttp.RequestCtx) {
	queryParameters, _ := url.ParseQuery(string(ctx.QueryArgs().QueryString()))

	var buckets []Bucket
	var filterFunctions []func (id int) bool
	selectFields := []string{"id", "email"}
	limit := 0
	
	for k, values := range queryParameters {
		v := values[0]
		switch k {
		case "sex_eq":
			buckets = append(buckets, createBucket(getSetFromIndex("sex", v)))
			selectFields = append(selectFields, "sex")
		case "status_eq":
			buckets = append(buckets, createBucket(getSetFromIndex("status", v)))
			selectFields = append(selectFields, "status")
		case "status_neq":
			statuses := []string{"свободны", "заняты", "всё сложно"}
			b := Bucket{}

			for _, item := range statuses {
				if item != v {
					b.add(getSetFromIndex("status", item))
				}
			}

			buckets = append(buckets, b)
			selectFields = append(selectFields, "status")
		case "email_domain":
			buckets = append(buckets, createBucket(getSetFromIndex("email_domain", v)))
			selectFields = append(selectFields, "email")
		case "email_lt":
			if disableSeqScan {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}

			filterFunctions = append(filterFunctions, func (id int) bool {
				account := getAccountFromData(id, []string{"email"})
				return account["email"].(string) < v
			})
			selectFields = append(selectFields, "email")
		case "email_gt":
			if disableSeqScan {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}

			filterFunctions = append(filterFunctions, func (id int) bool {
				account := getAccountFromData(id, []string{"email"})
				return account["email"].(string) > v
			})
			selectFields = append(selectFields, "email")
		case "fname_eq":
			buckets = append(buckets, createBucket(getSetFromIndex("fname", v)))
			selectFields = append(selectFields, "fname")
		case "fname_any":
			b := Bucket{}

			for _, item := range strings.Split(v, ",") {
				b.add(getSetFromIndex("fname", item))
			}

			buckets = append(buckets, b)
			selectFields = append(selectFields, "fname")
		case "fname_null":
			handleNullFilter(&selectFields, &buckets, "fname", v)
		case "sname_eq":
			buckets = append(buckets, createBucket(getSetFromIndex("sname", v)))
			selectFields = append(selectFields, "sname")
		case "sname_null":
			handleNullFilter(&selectFields, &buckets, "sname", v)
		case "country_eq":
			buckets = append(buckets, createBucket(getSetFromIndex("country", v)))
			selectFields = append(selectFields, "country")
		case "country_null":
			handleNullFilter(&selectFields, &buckets, "country", v)
		case "city_eq":
			buckets = append(buckets, createBucket(getSetFromIndex("city", v)))
			selectFields = append(selectFields, "city")
		case "city_any":
			b := Bucket{}

			for _, item := range strings.Split(v, ",") {
				b.add(getSetFromIndex("city", item))
			}

			buckets = append(buckets, b)
			selectFields = append(selectFields, "city")
		case "city_null":
			handleNullFilter(&selectFields, &buckets, "city", v)
		case "phone_code":
			buckets = append(buckets, createBucket(getSetFromIndex("phone_code", v)))
			selectFields = append(selectFields, "phone")
		case "phone_null":
			handleNullFilter(&selectFields, &buckets, "phone", v)
		case "premium_now":
			buckets = append(buckets, createBucket(getSetFromIndex("premium_now", "1")))
			selectFields = append(selectFields, "premium")
		case "premium_null":
			handleNullFilter(&selectFields, &buckets, "premium", v)
		case "sname_starts":
			b := Bucket{}

			for key, set := range indexes["sname"] {
				if strings.Index(key, v) == 0 {
					b.add(set)
				}
			}

			buckets = append(buckets, b)
			selectFields = append(selectFields, "sname")
		case "birth_lt":
			if disableSeqScan {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}

			k, _ := strconv.Atoi(v)
			filterFunctions = append(filterFunctions, func (id int) bool {
				account := getAccountFromData(id, []string{"birth"})
				return account["birth"].(int) < k
			})
			selectFields = append(selectFields, "birth")
		case "birth_gt":
			if disableSeqScan {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}

			k, _ := strconv.Atoi(v)
			filterFunctions = append(filterFunctions, func (id int) bool {
				account := getAccountFromData(id, []string{"birth"})
				return account["birth"].(int) > k
			})
			selectFields = append(selectFields, "birth")
		case "birth_year":
			buckets = append(buckets, createBucket(getSetFromIndex("birth_year", v)))
			selectFields = append(selectFields, "birth")
		case "interests_any":
			if disableSeqScan {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}

			b := Bucket{}

			for _, item := range strings.Split(v, ",") {
				b.add(getSetFromIndex("interests", item))
			}

			buckets = append(buckets, b)
		case "interests_contains":
			for _, item := range strings.Split(v, ",") {
				buckets = append(buckets, createBucket(getSetFromIndex("interests", item)))
			}
		case "likes_contains":
			for _, item := range strings.Split(v, ",") {
				id, _ := strconv.Atoi(item)
				ids, exists := likes[uint32(id)]
				s := createSet()

				if exists {
					for _, id := range ids {
						s.add(uint32(id))
					}
				}

				buckets = append(buckets, createBucket(s))
			}
		case "limit":
			limit, _ = strconv.Atoi(v)
		case "query_id":
			;
		default:
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}

	if limit <= 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if len(buckets) == 0 {
		b := Bucket{}
		b.add(getSetFromIndex("sex", "m"))
		b.add(getSetFromIndex("sex", "f"))
		buckets = append(buckets, b)
	}

	ids := getIds(buckets, filterFunctions, limit)
	var accounts []interface{}

	for _, id := range ids {
		accounts = append(accounts, getAccountFromData(id, selectFields))
	}

	writeResult(ctx, accounts)
}

func getIds(buckets []Bucket, filterFunctions []func (int) bool, limit int) []int {
	var ids []int

	if (len(buckets) == 0) {
		return ids
	}

	sort.Sort(byCount(buckets))
	tmp := make(map[int]bool)

	for !buckets[0].atEnd() && len(ids) < limit {
		contains := true
		id := buckets[0].getMaxAndMove()

		for _, bucket := range buckets {
			if !bucket.contains(id) {
				contains = false
				break
			}
		}

		for _, function := range filterFunctions {
			if !function(int(id)) {
				contains = false
				break
			}
		}
		_, exists := tmp[int(id)]
		if contains && ! exists {
			ids = append(ids, int(id))
			tmp[int(id)] = true
		}
	}

	return ids
}

func writeResult(ctx *fasthttp.RequestCtx, accounts []interface{}) {
	result := make(map[string]interface{})

	if len(accounts) > 0 {
		result["accounts"] = accounts
	} else {
		result["accounts"] = make([]int, 0)
	}

	output, _ := json.Marshal(result)
	fmt.Fprintf(ctx, string(output))
}

func getSetFromIndex(fieldName string, fieldValue string) *Set {
	set, ok := indexes[fieldName][fieldValue]	

	if ok {
		return set
	} else {
		return createSet()
	}
}

func handleNullFilter(selectFields *[]string, buckets *[]Bucket, indexName string, v string) {
	nullIndexName := fmt.Sprintf("%s_null", indexName)
	if v == "1" {
		*buckets = append(*buckets, createBucket(indexes[nullIndexName][""]))
	} else {
		*selectFields = append(*selectFields, indexName)
		*buckets = append(*buckets, createBucket(indexes[nullIndexName]["1"]))
	}
}
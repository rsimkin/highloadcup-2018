package main

import (
	"fmt"
	"sort"
	"strings"
	"strconv"
//	"time"
	"github.com/valyala/fasthttp"
	"net/url"
)

func handleRecommend(ctx *fasthttp.RequestCtx) {
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

	if phase > 1 {
		fmt.Fprintf(ctx, "{\"accounts\":[]}")
		return
	}

	var country string
	var city string
	var limit int

	queryParameters, _ := url.ParseQuery(string(ctx.QueryArgs().QueryString()))

	for k, values := range queryParameters {
		v := values[0]

		switch k {
		case "country":
			if v == "" {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			country = v
		case "city":
			if v == "" {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			city = v
		case "limit":
			limit, _ = strconv.Atoi(v)
		case "query_id":
			;
		default:
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}

	account := getAccountFromData(id, []string{"interests", "sex", "birth", "status", "birth_year", "country"})
	//fmt.Println("==============")
	//fmt.Println(account)
	//fmt.Println("==============")
	var resultIds []int
	resultMap := make(map[int]bool)

	if limit <= 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	for _, premium := range []string{"1", ""} {
		for _, status := range []string{"свободны", "всё сложно", "заняты"} {
			c := combination(account["interests"].([]string))
			for _, interestsBucket := range c {
				if len(interestsBucket) == 0 {
					continue
				}

				var currentIds []int
				//fmt.Println("===========")
				for _, interests := range interestsBucket {
					if len(interests) == 0 || interests[0] == "" {
						continue
					}
					//fmt.Println(len(interests), interests)
					////fmt.Print(fmt.Sprintf("premium: %s, ", premium))
					////fmt.Print(fmt.Sprintf("status: %s, ", status))
					var buckets []Bucket;
					if account["sex"].(string) == "m" {
						buckets = append(buckets, createBucket(getSetFromIndex("sex", "f")))
					} else {
						buckets = append(buckets, createBucket(getSetFromIndex("sex", "m")))
					}

					if country != "" {
						buckets = append(buckets, createBucket(getSetFromIndex("country", country)))
						////fmt.Print(fmt.Sprintf("country: '%s', ", country))
					}

					if city != "" {
						buckets = append(buckets, createBucket(getSetFromIndex("city", city)))
						////fmt.Print(fmt.Sprintf("city: '%s', ", city))
					}

					buckets = append(buckets, createBucket(getSetFromIndex("premium_now", premium)))
					buckets = append(buckets, createBucket(getSetFromIndex("status", status)))

					for _, k := range interests {
						buckets = append(buckets, createBucket(getSetFromIndex("interests", k)))
					}
					////fmt.Println("")

					ids := getIds(buckets, []func (int) bool{}, 10000)

					for _, currentId := range ids {
						currentIds = append(currentIds, currentId)
					}
				}

				//fmt.Println(currentIds)
				currentIds = getSortedIds(currentIds, account["birth"].(int))

				for _, id := range currentIds {
					_, exists := resultMap[id]

					if exists {
						continue
					}

					//fmt.Println(id)

					resultIds = append(resultIds, id)
					resultMap[id] = true

					if len(resultIds) == limit {
						//resultIds = getSortedIds(resultIds, account["birth"].(int))
						goto End
					}
				}
			}
		}
	}
	End:

	var accounts []interface{}
	selectFields := []string{"id", "email", "fname", "sname", "premium", "status", "birth"}

	for _, id := range resultIds {
		a := getAccountFromData(id, selectFields)
		if a["premium"].(map[string]int)["start"] == 0 {
			delete(a, "premium")
		}

		if a["sname"].(string) == "" {
			delete(a, "sname")
		}

		if a["fname"].(string) == "" {
			delete(a, "fname")
		}
		/*if account["birth"].(int) > a["birth"].(int) {
			a["diff"] = account["birth"].(int) - a["birth"].(int)
		} else {
			a["diff"] = a["birth"].(int) - account["birth"].(int)
		}*/
		accounts = append(accounts, a)
	}

	writeResult(ctx, accounts)
}

func getSortedIds(ids []int, birth int) []int {
	var data [][]int
	for _, id := range ids {
		account := getAccountFromData(id, []string{"birth"})
		currentBirth := account["birth"].(int)
		diff := 0

		if currentBirth > birth {
			diff = currentBirth - birth
		} else {
			diff = birth - currentBirth
		}

		data = append(data, []int{id, diff})
	}

	sort.Slice(data, func(i, j int) bool {
		if data[i][1] == data[j][1] {
			return data[i][0] < data[j][0]
		} else {
			return data[i][1] < data[j][1]
		}
	})

	var result []int

	for _, i := range data {
		result = append(result, i[0])
	}

	//fmt.Println(data)

	return result
}

func combination(array []string) [][][]string {
    var results [][]string
    results = append(results, make([]string, 0))

    for _, element := range array {
    	for _, combination := range results {
    		tmp := append(combination, element)
    		results = append(results, tmp)
    	}
    }

    hTmp := make(map[int][][]string)

    for _, g := range results {
    	_, exists := hTmp[len(g)]

    	if !exists {
    		hTmp[len(g)] = make([][]string, 0)
    	}

    	hTmp[len(g)] = append(hTmp[len(g)], g)
    }

    var keys []int
    for k := range hTmp {
        keys = append(keys, k)
    }
    sort.Ints(keys)

    var totalResult [][][]string

    for i := len(keys) - 1; i >= 0; i-- {
		totalResult = append(totalResult, hTmp[keys[i]])    	
    }

    //////fmt.Println(totalResult)

    return totalResult
}

func inArray(element string, data []string) bool {
	for _, i := range data {
		if i == element {
			return true
		}
	}

	return false
}
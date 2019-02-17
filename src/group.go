package main

import (
	"fmt"
	"sort"
	"net/url"
	"encoding/json"
	"strings"
	"strconv"
	"github.com/valyala/fasthttp"
)

func handleGroup(ctx *fasthttp.RequestCtx) {
	/*if phase > 2 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}*/
	//ctx.SetStatusCode(fasthttp.StatusBadRequest)

	queryParameters, _ := url.ParseQuery(string(ctx.QueryArgs().QueryString()))

	var set map[string]int
	var groupKeys []string
	var order string
	limit := 0
	set = groupIndexes["_"]["_"]

	queryKeys := make([]string, 0, len(queryParameters))

	for name := range queryParameters {
		queryKeys = append(queryKeys, name)
	}
	sort.Strings(queryKeys)

	var whereKeys []string
	var whereValues []string
	
	for _, k := range queryKeys {
		values := queryParameters[k]
		v := values[0]

		switch k {
		case "sex", "fname", "sname", "status", "country", "city", "birth", "joined":
			whereKeys = append(whereKeys, k)
			whereValues = append(whereValues, v)
		case "keys":
			for _, g := range strings.Split(v, ",") {
				switch g {
				case "sex", "country", "city", "status":
					groupKeys = append(groupKeys, g);
				default:
					ctx.SetStatusCode(fasthttp.StatusBadRequest)
					return;
				}
			}
		case "order":
			order = v
		case "limit":
			limit, _ = strconv.Atoi(v)
		case "query_id":
			;
		default:
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}

	if limit <= 0 || order == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if len(whereKeys) > 0 {
		whereKey := strings.Join(whereKeys, "|")
		whereValue := strings.Join(whereValues, "")
		_, exists := groupIndexes[whereKey]
		
		if (exists) {
			q, exists := groupIndexes[whereKey][whereValue]

			if !exists {
				set = nil
			} else {
				set = q
			}
		} else {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}

	resultTable := make(map[string]int)

	for key, cnt := range set {
		groupKey := getGroupKey(key, groupKeys)

		_, exists := resultTable[groupKey]

		if !exists {
			resultTable[groupKey] = 0
		}

		resultTable[groupKey] += cnt
	}

	var groupPairs []GroupPair

	for k, v := range resultTable {
		p := GroupPair{}
		p.Key = k
		p.Cnt = v
		groupPairs = append(groupPairs, p)
	}

	sort.Sort(byCnt(groupPairs))

	if order == "-1" {
		for i, j := 0, len(groupPairs)-1; i < j; i, j = i+1, j-1 {
	        groupPairs[i], groupPairs[j] = groupPairs[j], groupPairs[i]
	    }
	}

	var result []map[string]interface{}

	for _, p := range groupPairs {
		row := make(map[string]interface{})
		row["count"] = p.Cnt

		for i, k := range strings.Split(p.Key, "|") {
			if k != "" {
				row[groupKeys[i]] = k
			}
		}
		if p.Cnt > 0 {
			result = append(result, row)
		}

		if len(result) == limit {
			break
		}
	}

	writeGroupResult(ctx, result)
}

func getGroupKey(storedKey string, groupKeys []string) string {
	values := strings.Split(storedKey, "|")
	var tmp []string

	for _, f := range groupKeys {
		switch f {
			case "sex":
				tmp = append(tmp, values[0])
			case "status":
				tmp = append(tmp, values[1])
			case "country":
				tmp = append(tmp, values[2])
			case "city":
				tmp = append(tmp, values[3])
		}
	}

	return strings.Join(tmp, "|")
}

func writeGroupResult(ctx *fasthttp.RequestCtx, rows []map[string]interface{}) {
	result := make(map[string]interface{})

	if len(rows) > 0 {
		result["groups"] = rows
	} else {
		result["groups"] = make([]int, 0)
	}

	output, _ := json.Marshal(result)
	fmt.Fprintf(ctx, string(output))
}
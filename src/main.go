package main

import (
	"fmt"
    "github.com/valyala/fasthttp"
    "strings"
    "sync"
    "time"
    "runtime"
)

var indexes map[string]map[string]*Set
var groupIndexes map[string]map[string]map[string]int
var accountsData map[int]string
var registry map[string]int
var emails map[string]bool
var phones map[string]bool
var likes map[uint32][]uint32
var timestamp uint
var mode uint
var postCount = 0
var getCount = 0
var mx sync.RWMutex
var phase = 0
var disableSeqScan = false

func main() {
    fmt.Println("Start")

    indexes = make(map[string]map[string]*Set);
    groupIndexes = make(map[string]map[string]map[string]int)
    registry = make(map[string]int);
    emails = make(map[string]bool);
	phones = make(map[string]bool);
	accountsData = make(map[int]string);
	likes = make(map[uint32][]uint32)

	runtime.GOMAXPROCS(runtime.NumCPU())

    loadOptions("/tmp/data/options.txt")
    loadFiles("/tmp/data/data.zip")
    sortIndexes()
    sortLikes()

    PrintMemUsage()

    ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
	    for {
	       select {
	        case <- ticker.C:
	            //PrintMemUsage()
	        case <- quit:
	            ticker.Stop()
	            return
	        }
	    }
	 }()

    lastMethod := ""
    resorted := false
    limit := 10000 - 1

    if mode == 1 {
        limit = 90000-1 //88200 //
    }

    requestHandler := func(ctx *fasthttp.RequestCtx) {
        if string(ctx.Method()) == "POST" {
            postCount++
        }
		
		if string(ctx.Method()) == "GET" {
            getCount++

            if !disableSeqScan && getCount > 27000 + 60000 / 4 {
            	//disableSeqScan = true
            	//fmt.Println("Seq scan disabled")
            }
        }

        if lastMethod == "" || lastMethod != string(ctx.Method()) {
            phase++
            fmt.Println(fmt.Sprintf("Got first reqest %s phase %d", string(ctx.Method()), phase))
            fmt.Println(fmt.Sprintf("Post count %d", postCount))
            fmt.Println(fmt.Sprintf("Get count %d", getCount))
            PrintMemUsage()
            lastMethod = string(ctx.Method())
        }

        ctx.SetContentType("application/json")

        switch string(ctx.Path()) {
        case "/accounts/filter/":
            mx.RLock()
            handleFilter(ctx)
            mx.RUnlock()
		case "/accounts/group/":
            mx.RLock()
            handleGroup(ctx)
            mx.RUnlock()
        case "/accounts/new/":
        	mx.Lock()
            handleCreate(ctx)
            mx.Unlock()
        case "/accounts/likes/":
            mx.RLock()
            handleLikes(ctx)
            mx.RUnlock()

        default:
            if string(ctx.Method()) == "GET" && strings.Index(string(ctx.Path()), "/recommend/") != -1 {
                handleRecommend(ctx)
            } else if string(ctx.Method()) == "GET" && strings.Index(string(ctx.Path()), "/suggest/") != -1 {
                handleSuggest(ctx)
            } else if string(ctx.Method()) == "POST" && strings.Index(string(ctx.Path()), "/accounts/") == 0 {
                mx.Lock()
                handleUpdate(ctx)
                mx.Unlock()
            } else {
                ctx.Error("Unsupported path", fasthttp.StatusNotFound)
            }
        }

        if !resorted && postCount == limit {
            resorted = true
            time.Sleep(2 * time.Second)
            mx.Lock()
            sortIndexes()
            mx.Unlock()
        }
    }
    
    fasthttp.ListenAndServe(":80", requestHandler)
}
package main

import (
	"fmt"
	"time"
	"io/ioutil"
	"strconv"
	"strings"
	"archive/zip"
	"encoding/json"
	"runtime"
	"sort"
)

func loadOptions(filePath string) {
	b, err := ioutil.ReadFile(filePath)
    
    if err != nil {
        panic(err)
    }
	
	strings := strings.Split(string(b), "\n")
	tInt32, _ := strconv.Atoi(strings[0])
	timestamp = uint(tInt32)
	fmt.Println("time from file is", timestamp)

	tInt32, _ = strconv.Atoi(strings[1])
	mode = uint(tInt32)
	fmt.Println("mode from file is", mode)	
}

func loadFiles(filePath string) {
    start := time.Now()
	r, err := zip.OpenReader(filePath)
	
	if err != nil {
		panic(err)
	}
	
	defer r.Close()

	for _, f := range r.File {
		var accounts Accounts
		//fmt.Printf("%d from %d. Loading %s:\n", z, len(r.File), f.Name)
		jsonFile, err := f.Open()
		
		if err != nil {
			panic(err)
		}

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		json.Unmarshal(byteValue, &accounts)
		jsonFile.Close()

		processAccounts(accounts.Accounts);
		runtime.GC()

		accounts = Accounts{}
		//PrintMemUsage()
	}

	elapsed := time.Since(start)
    fmt.Println("Loading took ", elapsed)
}

func processAccounts(accounts []Account) {
	for _, account := range accounts {
		processAccount(account)
		//processAccountLikes(account)
	}
}

func processAccount(account Account) {
	emailDomain := ""
	parts := strings.Split(account.Email, "@")

	if (len(parts) > 0) {
		emailDomain = parts[1]
	}

	premium := ""

	if account.Premium.Start > 0 && account.Premium.Finish > 0 {
		premium = "1"
	}

	phoneCode := ""

	if len(account.Phone) > 5 {
		phoneCode = account.Phone[2:5]
	}

	premiumNow := ""

	if account.Premium.Start <= int(timestamp) && account.Premium.Finish >= int(timestamp) {
		premiumNow = "1"
	}

	computed := make(map[string]string)
	computed["fname"] = account.Fname
	computed["sname"] = account.Sname
	computed["status"] = account.Status
	computed["sex"] = account.Sex
	computed["country"] = account.Country
	computed["city"] = account.City
	computed["email_domain"] = emailDomain
	computed["phone_code"] = phoneCode
	computed["premium_now"] = premiumNow
	computed["birth_year"] = strconv.Itoa(time.Unix(int64(account.Birth), 0).Year())
	computed["fname_null"] = getNullValue(account.Fname)
	computed["sname_null"] = getNullValue(account.Sname)
	computed["city_null"] = getNullValue(account.City)
	computed["country_null"] = getNullValue(account.Country)
	computed["phone_null"] = getNullValue(account.Phone)
	computed["premium_null"] = premium

	for indexName, indexValue := range computed {
		addIdToIndex(indexName, indexValue, int(account.Id))
	}

	for _, interest := range account.Interests {
		addIdToIndex("interests", interest, int(account.Id))
	}

	addAccountToData(account)

	birth := strconv.Itoa(time.Unix(int64(account.Birth), 0).Year())
	joined := strconv.Itoa(time.Unix(int64(account.Joined), 0).Year())

	addAccountToGroupIndex("_", "_", account, 1)
	addAccountToGroupIndex("fname", account.Fname, account, 1)
	addAccountToGroupIndex("sname", account.Sname, account, 1)
	addAccountToGroupIndex("city", account.City, account, 1)
	addAccountToGroupIndex("country", account.Country, account, 1)
	addAccountToGroupIndex("status", account.Status, account, 1)
	addAccountToGroupIndex("sex", account.Sex, account, 1)
	addAccountToGroupIndex("birth", birth, account, 1)
	addAccountToGroupIndex("joined", joined, account, 1)
	
	addAccountToGroupIndex("birth|city", birth + account.City, account, 1)
	addAccountToGroupIndex("birth|status", birth + account.Status, account, 1)
	addAccountToGroupIndex("joined|sex", joined + account.Sex, account, 1)
	addAccountToGroupIndex("joined|status", joined + account.Status, account, 1)
	addAccountToGroupIndex("birth|sex", birth + account.Sex, account, 1)
	addAccountToGroupIndex("city|joined", account.City + joined, account, 1)
	addAccountToGroupIndex("country|joined", account.Country + joined, account, 1)
	addAccountToGroupIndex("birth|country", birth + account.Country, account, 1)
}

func addAccountToGroupIndex(fieldName string, value string, account Account, valueCnt int) {
	var tmp []string
	var exists bool
	tmp = append(tmp, account.Sex)
	tmp = append(tmp, account.Status)
	tmp = append(tmp, account.Country)
	tmp = append(tmp, account.City)
	groupKey := strings.Join(tmp, "|")

	_, exists = groupIndexes[fieldName]

	if !exists {
		groupIndexes[fieldName] = make(map[string]map[string]int)
	}

	_, exists = groupIndexes[fieldName][value]

	if !exists {
		groupIndexes[fieldName][value] = make(map[string]int)
	}

	_, exists = groupIndexes[fieldName][value][groupKey]

	if !exists {
		groupIndexes[fieldName][value][groupKey] = 0
	}

	groupIndexes[fieldName][value][groupKey] += valueCnt
}

func getNullValue(value string) string {
	if value == "" {
		return ""
	} else {
		return "1"
	}
}

func addIdToIndex(indexName string, value string, id int) {
	_, ok := indexes[indexName]

	if !ok {
		indexes[indexName] = make(map[string]*Set)
	}

	_, ok = indexes[indexName][value]

	if !ok {
		indexes[indexName][value] = createSet()
	}

	_, exists := indexes[indexName][value].data[uint32(id)]

	if !exists {
		indexes[indexName][value].add(uint32(id))
	}
}

func deleteIdFromIndex(indexName string, value string, id int) {
	delete(indexes[indexName][value].data, uint32(id))
}

func addAccountToData(account Account) {
	var interests []string

	if phase < 2 {
		for _, i := range account.Interests {
			interests = append(interests, addToRegistry(i))
		}
	}

	h := []string{
		addToRegistry(account.Fname),
		addToRegistry(account.Sname),
		addToRegistry(account.Status),
		strconv.Itoa(int(account.Premium.Start)),
		strconv.Itoa(int(account.Premium.Finish)),
		account.Email,
		account.Sex,
		account.Phone,
		strconv.Itoa(int(account.Birth)),
		addToRegistry(account.City),
		addToRegistry(account.Country),
		strconv.Itoa(int(account.Joined)),
		strings.Join(interests, ","),
	}
	accountsData[int(account.Id)] = strings.Join(h, "|")
	emails[account.Email] = true

	if account.Phone != "" {
		phones[account.Phone] = true
	}
}

func getAccountFromData(id int, selectFields []string) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = id
	f := accountsData[id]
	h := strings.Split(f, "|")

	for _, field := range selectFields {
		switch field {
		case "fname":
			data["fname"] = getFromRegistry(h[0])
		case "sname":
			data["sname"] = getFromRegistry(h[1])
		case "status":
			data["status"] = getFromRegistry(h[2])
		case "premium":
			tmp := make(map[string]int)
			s, _ := strconv.Atoi(h[3])
			f, _ := strconv.Atoi(h[4])
			tmp["start"] = s
			tmp["finish"] = f
			data["premium"] = tmp
		case "email":
			data["email"] = h[5]
		case "sex":
			data["sex"] = h[6]
		case "phone":
			data["phone"] = h[7]
		case "birth":
			s, _ := strconv.Atoi(h[8])
			data["birth"] = s
		case "city":
			data["city"] = getFromRegistry(h[9])
		case "country":
			data["country"] = getFromRegistry(h[10])
		case "joined":
			s, _ := strconv.Atoi(h[11])
			data["joined"] = s
		case "interests":
			var interests []string

			for _, i := range strings.Split(h[12], ",") {
				interests = append(interests, getFromRegistry(i))
			}

			data["interests"] = interests
		}
	}

	return data
}

func getAccountObjectFromData(id int) Account {
	var account Account
	data := getAccountFromData(id, []string{"fname", "sname", "status", "premium", "email", "sex", "phone", "birth", "city", "country", "joined"})
	account.Id = id
	account.Fname = data["fname"].(string)
	account.Sname = data["sname"].(string)
	account.Status = data["status"].(string)
	account.Premium.Start = data["premium"].(map[string]int)["start"]
	account.Premium.Finish = data["premium"].(map[string]int)["finish"]
	account.Email = data["email"].(string)
	account.Sex = data["sex"].(string)
	account.Phone = data["phone"].(string)
	account.Birth = data["birth"].(int)
	account.City = data["city"].(string)
	account.Country = data["country"].(string)
	account.Joined = data["joined"].(int)

	return account
}

func processAccountLikes(account Account) {
	for _, like := range account.Likes {
		_, exists := likes[uint32(like.Id)]

		if !exists {
			likes[uint32(like.Id)] = []uint32{uint32(account.Id)}
		} else {
			likes[uint32(like.Id)] = append(likes[uint32(like.Id)], uint32(account.Id))
		}
	}
}

func sortIndexes() {
	start := time.Now()

	for _, indexInfo := range indexes {
		for _, set := range indexInfo {
			set.sort()
		}
	}

	elapsed := time.Since(start)
    fmt.Println("Sorting took ", elapsed)
}

func sortLikes() {
	start := time.Now()

	for _, likesBucket := range likes {
		sort.Slice(likesBucket, func(i, j int) bool { return likesBucket[i] < likesBucket[j] })
	}

	elapsed := time.Since(start)
    fmt.Println("Sorting likes took ", elapsed)
}
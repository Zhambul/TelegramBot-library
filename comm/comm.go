package comm

import (
	"strconv"
	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"bytes"
	"fmt"
	"io"
)

var token string

func Init(tkn string) {
	token = tkn
}

//TODO HTTP URL PARAMS
func GetUpdates(offset int) *Updates {
	u := &Updates{}
	updateIds := make([]int, 0)
	log.Println("RunUpdateLoop START")
	params := make(map[string]string)
	if offset != 0 {
		params["offset"] = strconv.Itoa(offset)
	}
	params["timeout"] = strconv.Itoa(1000)
	resp, err := get("getUpdates", params)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(data, &objmap)
	if err != nil {
		log.Fatal(err)
	}

	results := make([]*MessageInfo, 0)
	err = json.Unmarshal(*objmap["result"], &results)
	if err != nil {
		panic(err)
	}

	for _, res := range results {
		updateIds = append(updateIds, res.UpdateId)

		log.Printf("updateId - %v\n", res.UpdateId)
		if res.Message != nil {
			res.Message.UpdateId = res.UpdateId
			log.Printf("got message - %+v", res.Message)
			u.Messages = append(u.Messages, res.Message)
		}

		if res.Callback != nil {
			res.Callback.UpdateId = res.UpdateId
			log.Printf("got callbackQuery - %+v", res.Callback)
			u.Callbacks = append(u.Callbacks, res.Callback)
		}

		if res.Inline != nil {
			res.Inline.UpdateId = res.UpdateId
			log.Printf("got inlineQuery - %+v", res.Inline)
			u.Inlines = append(u.Inlines, res.Inline)
		}
	}
	log.Println("RunUpdateLoop END")
	u.NextUpdateId = getNextUpdateId(updateIds)
	return u
}

func getNextUpdateId(updateIds []int) int {
	updatesLen := len(updateIds)
	if updatesLen == 0 {
		return 0
	}
	return updateIds[updatesLen-1] + 1
}

type MessageIdResp struct {
	Result struct {
		MessageId int `json:"message_id"`
	}    `json:"result"`
}

func Update(url string, body interface{}) error {
	bytes, _ := json.Marshal(body)
	_, err := post(url, bytes)
	return err
}

func Send(url string, body interface{}) (int, error) {
	bytes, _ := json.Marshal(body)
	resp, err := post(url, bytes)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	messageIdResp := MessageIdResp{}
	respBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respBytes, &messageIdResp)
	return messageIdResp.Result.MessageId, err
}

func post(methodName string, body []byte) (*http.Response, error) {
	url := requestUrl(methodName)

	log.Printf("HTTP POST - %v\n", url)
	log.Println(string(body))
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if resp.StatusCode < 199 || resp.StatusCode > 299 {
		data := bodyToString(resp.Body)

		return resp, fmt.Errorf("HTTP ERROR: url - %v\n, status code - %v body - %v\n", url, resp.StatusCode, data)
	}
	return resp, err
}

func get(methodName string, params map[string]string) (*http.Response, error) {
	url := requestUrl(methodName)
	if params != nil && len(params) != 0 {
		url = url + "?"
		for key, value := range params {
			url = url + key + "=" + value + "&"
		}
	}

	log.Printf("HTTP GET - %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 199 || resp.StatusCode > 299 {
		return resp, fmt.Errorf("url - %v\n, status code - %v\n", url, resp.StatusCode)
	}
	return resp, err
}

func requestUrl(methodName string) string {
	return "https://api.telegram.org/bot" + token + "/" + methodName
}

func bodyToString(body io.ReadCloser) string {
	defer body.Close()
	bodyBytes, _ := ioutil.ReadAll(body)
	return string(bodyBytes)
}

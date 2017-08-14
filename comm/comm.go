package comm

import (
	"strconv"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"fmt"
	"io"
	"errors"
	"bytes"
	"time"
	"log"
)

var token string

func Init(tkn string) {
	log.Println("comm::Init")
	token = tkn
}

func GetUpdates(offset int) (*Updates, error) {
	log.Println("Comm::GetUpdates START")
	resp, err := getUpdatesQuery(offset)
	if err != nil {
		log.Printf("Comm::GetUpdates END. Cannot get update, %v\n", err)
		return nil, errors.New("cannot get update")
	}

	data, err := ioutil.ReadAll(resp.Body)

	var objmap map[string]*json.RawMessage
	json.Unmarshal(data, &objmap)
	results := make([]*MessageInfo, 0)
	json.Unmarshal(*objmap["result"], &results)

	u := &Updates{}
	updateIds := make([]int, 0)
	for _, res := range results {
		updateIds = append(updateIds, res.UpdateId)

		if res.Message != nil {
			res.Message.UpdateId = res.UpdateId
			log.Printf("Comm::GetUpdates Message")
			u.Messages = append(u.Messages, res.Message)
		}

		if res.Callback != nil {
			res.Callback.UpdateId = res.UpdateId
			log.Printf("Comm::GetUpdates Callback")
			u.Callbacks = append(u.Callbacks, res.Callback)
		}

		if res.Inline != nil {
			res.Inline.UpdateId = res.UpdateId
			log.Printf("Comm::GetUpdates Inline")
			u.Inlines = append(u.Inlines, res.Inline)
		}
	}
	u.NextUpdateId = getNextUpdateId(updateIds)
	log.Println("Comm::GetUpdates END")
	return u, nil
}

func UpdateMessage(reply *Reply) error {
	log.Println("Comm::UpdateMessage")
	return update("editMessageText", reply)
}

func SendMessage(reply *Reply) (int, error) {
	log.Println("Comm::SendMessage")
	return send("sendMessage", reply)
}

func AnswerInlineQuery(a *InlineQueryAnswer) error {
	log.Println("Comm::AnswerInlineQuery")
	return update("answerInlineQuery", a)
}

func DeleteMessage(d *DeleteMsg) error {
	log.Println("Comm::DeleteMessage")
	return update("deleteMessage", d)
}

func update(url string, body interface{}) error {
	log.Println("Comm::update START")
	bytes, _ := json.Marshal(body)
	_, err := post(url, bytes)
	log.Println("Comm::update END")
	return err
}

func send(url string, body interface{}) (int, error) {
	log.Println("Comm::send START")
	bytes, _ := json.Marshal(body)
	resp, err := post(url, bytes)
	if err != nil {
		log.Printf("Comm::send END. ERROR - %v\n", err)
		return 0, err
	}
	defer resp.Body.Close()
	messageIdResp := MessageIdResp{}
	respBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respBytes, &messageIdResp)
	log.Println("Comm::send END")
	return messageIdResp.Result.MessageId, err
}

func getNextUpdateId(updateIds []int) int {
	log.Println("Comm::getNextUpdateId")
	updatesLen := len(updateIds)
	if updatesLen == 0 {
		return 0
	}
	return updateIds[updatesLen-1] + 1
}

func post(methodName string, body []byte) (*http.Response, error) {
	log.Println("Comm::post")
	return request("POST", methodName, nil, bytes.NewReader(body))
}

func getUpdatesQuery(offset int) (*http.Response, error) {
	log.Println("Comm::getUpdatesQuery")
	params := make(map[string]string)
	if offset != 0 {
		params["offset"] = strconv.Itoa(offset)
	}
	params["timeout"] = strconv.Itoa(1000)
	return request("GET", "getUpdates", params, nil)
}

var client = &http.Client{
	Timeout: time.Second * 100,
}

func request(method string, telegramMethod string, params map[string]string, body io.Reader) (*http.Response, error) {
	log.Println("Comm::request START")
	req, err := http.NewRequest(method, requestUrl(telegramMethod), body)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	if params != nil && len(params) != 0 {
		for key, value := range params {
			q.Add(key, value)
		}
	}
	req.URL.RawQuery = q.Encode()
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode < 199 || resp.StatusCode > 299 {
		return resp, fmt.Errorf("url - %v, status code - %v, body - %v\n", req.URL.String(),
			resp.StatusCode, bodyToString(resp.Body))
	}
	log.Println("Comm::request END")
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

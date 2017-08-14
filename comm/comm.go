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
)

var token string

func Init(tkn string) {
	token = tkn
}

func GetUpdates(offset int) (*Updates, error) {
	resp, err := getUpdatesQuery(offset)
	if err != nil {
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
			u.Messages = append(u.Messages, res.Message)
		}

		if res.Callback != nil {
			res.Callback.UpdateId = res.UpdateId
			u.Callbacks = append(u.Callbacks, res.Callback)
		}

		if res.Inline != nil {
			res.Inline.UpdateId = res.UpdateId
			u.Inlines = append(u.Inlines, res.Inline)
		}
	}
	u.NextUpdateId = getNextUpdateId(updateIds)
	return u, nil
}

func UpdateMessage(reply *Reply) error {
	return update("editMessageText", reply)
}

func SendMessage(reply *Reply) (int, error) {
	return send("sendMessage", reply)
}

func AnswerInlineQuery(a *InlineQueryAnswer) error {
	return update("answerInlineQuery", a)
}

func DeleteMessage(d *DeleteMsg) error {
	return update("deleteMessage", d)
}

func update(url string, body interface{}) error {
	bytes, _ := json.Marshal(body)
	_, err := post(url, bytes)
	return err
}

func send(url string, body interface{}) (int, error) {
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

func getNextUpdateId(updateIds []int) int {
	updatesLen := len(updateIds)
	if updatesLen == 0 {
		return 0
	}
	return updateIds[updatesLen-1] + 1
}

func post(methodName string, body []byte) (*http.Response, error) {
	return request("POST", methodName, nil, bytes.NewReader(body))
}

func getUpdatesQuery(offset int) (*http.Response, error) {
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
		return resp, fmt.Errorf("url - %v, status code - %v\n", req.URL.String(), resp.StatusCode)
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

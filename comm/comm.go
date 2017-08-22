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

var offset int

var webhookEnabled bool
var webhookUpdatesChan chan (*Updates)

func EnableWebhook(host string) error {
	webhookEnabled = true
	url := host + "/webhook"
	log.Printf("Setting webhook to url %v\n", url)
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Webhook START")
		print(bodyToString(r.Body))
		log.Println("Webhook END")
		w.Write([]byte("Hi!"))
	})
	log.Fatal(http.ListenAndServe(":80", nil))
	r, err := post("setWebhook", webhook{
		url: url,
	})
	log.Println(bodyToString(r.Body))
	return err
}

func GetUpdates() (*Updates, error) {
	log.Println("Comm::GetUpdates START")
	defer func() {
		log.Println("Comm::GetUpdates END")
	}()

	if !webhookEnabled {
		return pullUpdates()
	} else {
		log.Println("Comm::GetUpdates. Waiting for update")
		return <-webhookUpdatesChan, nil
	}
}

func pullUpdates() (*Updates, error) {
	log.Println("Comm::pullUpdates START")
	resp, err := getUpdatesQuery(offset)
	if err != nil {
		log.Printf("Comm::pullUpdates END. Cannot get update, %v\n", err)
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
			log.Printf("Comm::pullUpdates Message")
			u.Messages = append(u.Messages, res.Message)
		}

		if res.Callback != nil {
			res.Callback.UpdateId = res.UpdateId
			log.Printf("Comm::pullUpdates Callback")
			u.Callbacks = append(u.Callbacks, res.Callback)
		}

		if res.Inline != nil {
			res.Inline.UpdateId = res.UpdateId
			log.Printf("Comm::pullUpdates Inline")
			u.Inlines = append(u.Inlines, res.Inline)
		}
	}
	offset = getNextUpdateId(updateIds)
	log.Println("Comm::pullUpdates END")
	return u, nil
}

func UpdateMessage(reply *Reply) error {
	log.Println("Comm::UpdateMessage")
	_, err := post("editMessageText", reply)
	return err
}

func SendMessage(reply *Reply) (int, error) {
	log.Println("Comm::SendMessage")
	resp, err := post("sendMessage", reply)
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

func AnswerInlineQuery(a *InlineQueryAnswer) error {
	log.Println("Comm::AnswerInlineQuery")
	_, err := post("answerInlineQuery", a)
	return err
}

func DeleteMessage(d *DeleteMsg) error {
	log.Println("Comm::DeleteMessage")
	_, err := post("deleteMessage", d)
	return err
}

func getNextUpdateId(updateIds []int) int {
	log.Println("Comm::getNextUpdateId")
	updatesLen := len(updateIds)
	if updatesLen == 0 {
		return 0
	}
	return updateIds[updatesLen-1] + 1
}

func getUpdatesQuery(offset int) (*http.Response, error) {
	log.Println("Comm::getUpdatesQuery")
	params := make(map[string]string)
	if offset != 0 {
		params["offset"] = strconv.Itoa(offset)
	}
	params["timeout"] = strconv.Itoa(1000)
	return get("getUpdates", params)
}

var client = &http.Client{
	Timeout: time.Second * 100,
}

func get(methodName string, params map[string]string) (*http.Response, error) {
	return request("GET", methodName, params, nil)
}

func post(methodName string, body interface{}) (*http.Response, error) {
	log.Println("Comm::post")
	b, _ := json.Marshal(body)
	return request("POST", methodName, nil, bytes.NewReader(b))
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

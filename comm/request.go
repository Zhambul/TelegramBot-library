package comm

type Updates struct {
	NextUpdateId int
	Messages     []*Message
	Callbacks    []*Callback
	Inlines      []*Inline
}

type MessageInfo struct {
	UpdateId int `json:"update_id"`
	Message  *Message `json:"message, omitempty"`
	Callback *Callback `json:"callback_query, omitempty"`
	Inline   *Inline `json:"inline_query, omitempty"`
}

type Message struct {
	UpdateId int
	Chat struct {
		Id int `json:"id"`
	} `json:"chat"`
	From           *From `json:"from"`
	Text           string `json:"text"`
	MessageId      int `json:"message_id"`
	ReplyToMessage *Message `json:"reply_to_message"`
}

type Callback struct {
	UpdateId     int
	Id           string `json:"id"`
	From         *From `json:"from"`
	CallbackData string `json:"data"`
	Message      *Message `json:"message"`
}

type Inline struct {
	UpdateId int `json:"update_id"`
	Id       string `json:"id"`
	Query    string `json:"query"`
	From     *From `json:"from"`
}

type From struct {
	Id        int `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

package comm

type webhook struct {
	url string
}
type Reply struct {
	ChatId      int `json:"chat_id"`
	Text        string `json:"text"`
	MessageId   int `json:"message_id,omitempty"`
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	ParseMode   string    `json:"parse_mode,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]*InlineKeyboardButton `json:"inline_keyboard,omitempty"`
}

type InlineKeyboardButton struct {
	Text              string `json:"text,omitempty"`
	Url               string `json:"url,omitempty"`
	CallbackData      string `json:"callback_data,omitempty"`
	SwitchInlineQuery string `json:"switch_inline_query,omitempty"`
}

type InputMessageContent struct {
	MessageText string `json:"message_text"`
	ParseMode   string `json:"parse_mode"`
}

type InlineQueryResult struct {
	Type                string `json:"type"`
	Id                  string `json:"id"`
	Title               string `json:"title"`
	InputMessageContent *InputMessageContent `json:"input_message_content"`
	Url                 string `json:"url"`
	HideUrl             bool `json:"hide_url"`
	Description         string `json:"description"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup"`
}

type InlineQueryAnswer struct {
	InlineQueryId string `json:"inline_query_id"`
	Results       []*InlineQueryResult `json:"results"`
	CacheTime     int `json:"cache_time"`
}

type DeleteMsg struct {
	ChatId    int `json:"chat_id"`
	MessageId int `json:"message_id"`
}

type MessageIdResp struct {
	Result struct {
		MessageId int `json:"message_id"`
	}    `json:"result"`
}

package bot

type Handler interface {
	Handle(*Context) *Response
}

type InlineHandler interface {
	Handle(*Context) *InlineAnswer
}
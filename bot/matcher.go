package bot

type Matcher interface {
	Match(*Context) bool
}

type SimpleMatcher struct {
	Text string
}

func (m *SimpleMatcher) Match(c *Context) bool {
	if m.Text == "" {
		panic("Not defined matcher pattern")
	}
	return c.Message.Text == m.Text
}
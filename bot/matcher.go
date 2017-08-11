package bot

type Matcher interface {
	Match(*Context) bool
}

type simpleMatcher struct {
	Text string
}

func (m *simpleMatcher) Match(c *Context) bool {
	if m.Text == "" {
		panic("Not defined matcher pattern")
	}
	return c.Message.Text == m.Text
}
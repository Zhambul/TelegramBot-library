package bot

import "log"

type Matcher interface {
	Match(*Context) bool
}

type simpleMatcher struct {
	Text string
}

func (m *simpleMatcher) Match(c *Context) bool {
	log.Println("SimpleMatcher::Match")
	if m.Text == "" {
		panic("Not defined matcher pattern")
	}
	return c.Message.Text == m.Text
}
package bot

import (
	"log"
	"fmt"
)

type contextLogger struct {
	Account *BotAccount
	Context *Context
}

func newContextLogger(account *BotAccount, context *Context) *contextLogger {
	log.Println("Bot::newContextLogger")
	return &contextLogger{Account: account, Context: context }
}

func (l *contextLogger) debug(msg string, v ...interface{}) {
	l.log(msg, "DEBUG", v...)
}

func (l *contextLogger) info(msg string, v ...interface{}) {
	l.log(msg, "INFO", v...)
}

func (l *contextLogger) err(msg string, v ...interface{}) {
	l.log(msg, "ERROR", v...)
}

func (l *contextLogger) log(msg, level string, v ...interface{}) {
	formated := fmt.Sprintf(msg, v...)
	log.Printf("%v: %v | acc: %+v, context %+v\n", level, formated,
		l.Account, l.Context)
}

package smtpd

import (
	"fmt"
	"net"
	"regexp"
	"time"
)

const (
	mod_COMMAND = 1
	mod_HEAD    = 2
	mod_BODY    = 4
)

type smtp_context_t struct {
	sock     net.Conn
	Address  string
	handlers SmtpServerConfigure
	conf     *ServerConfig
	Module   int
	Login    bool
	re       *regexp.Regexp
	Email    *Mail
	// 其他
	Msg  string
	User string
}

func initSmtpContext(sock net.Conn, config SmtpServerConfigure) *smtp_context_t {
	scxt := &smtp_context_t{
		sock:     sock,
		Address:  sock.RemoteAddr().String(),
		handlers: config,
		conf:     config.GetConfig(),
		Module:   mod_COMMAND,
		Login:    false,
		re:       regexp.MustCompile("<(.+)>"),
		Email:    &Mail{},
	}

	return scxt
}

// 发送
func (scxt *smtp_context_t) Send(content string) {
	fmt.Fprint(scxt.sock, content)
}

// 发送并关闭
func (scxt *smtp_context_t) End(content string) {
	fmt.Fprint(scxt.sock, content)
	scxt.sock.Close()
}

// 问候语
func (scxt *smtp_context_t) hola() {
	config := scxt.conf
	scxt.Send(fmt.Sprintf("220 %s %s Server (%s %s Server %s) ready %d\r\n",
		config.Domain, config.Type, config.Name, config.Type, config.Version, time.Now().Unix(),
	))
}

func (scxt *smtp_context_t) takeOff() {
	scxt.handlers.TakeOff(scxt.Email)
}

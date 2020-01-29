package smtpd

import (
    "fmt"
    "net"
    "regexp"
    "time"
)

const (
    MOD_COMMAND = 1
    MOD_HEAD = 2
    MOD_BODY = 4
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
    Msg      string
    User     string
}

func initSmtpContext(sock net.Conn, config SmtpServerConfigure) *smtp_context_t {
    this := &smtp_context_t{
        sock: sock,
        Address: sock.RemoteAddr().String(),
        handlers: config,
        conf: config.GetConfig(),
        Module: MOD_COMMAND,
        Login: false,
        re: regexp.MustCompile("<(.+)>"),
        Email: &Mail{},
    }

    return this
}

// 发送
func (this *smtp_context_t) Send(content string) {
    fmt.Fprint(this.sock, content)
}

// 发送并关闭
func (this *smtp_context_t) End(content string) {
    fmt.Fprint(this.sock, content)
    this.sock.Close()
}

// 问候语
func (this *smtp_context_t) hola() {
    config := this.conf
	this.Send(fmt.Sprintf("220 %s %s Server (%s %s Server %s) ready %d\r\n",
        config.Domain, config.Type, config.Name, config.Type, config.Version, time.Now().Unix(),
	))
}

func (this *smtp_context_t) takeOff() {
    this.handlers.TakeOff(this.Email)
}

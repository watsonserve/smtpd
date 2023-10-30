package smtpd

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"
)

// helo命令
func helo(ctx *smtp_context_t) {
	ctx.Module = mod_COMMAND
	addr := ctx.Address
	name := ctx.Msg[5:]

	ctx.Send(fmt.Sprintf("250 %s Hello %s (%s[%s])\r\n", ctx.conf.Domain, name, addr, addr))
}

// ehlo命令
func ehlo(ctx *smtp_context_t) {
	ctx.Module = mod_COMMAND
	addr := ctx.Address
	name := ctx.Msg[5:]
	msg := fmt.Sprintf(
		"250-%s Hello %s (%s[%s])\r\n%s\r\n%s\r\n%s\r\n%s\r\n",
		ctx.conf.Domain, name, addr, addr,
		"250-AUTH LOGIN PLAIN",
		"250-AUTH=LOGIN PLAIN",
		"250-PIPELINING",
		"250 ENHANCEDSTATUSCODES",
	)
	ctx.Send(msg)
}

// 授权
func auth(ctx *smtp_context_t) {
	content, err := base64.StdEncoding.DecodeString(ctx.Msg[11:])
	if nil != err {
		log.Printf("error: %s\n", err.Error())
		return
	}

	for i := 0; i < len(content); i++ {
		if 0 == content[i] {
			content[i] = '\n'
		}
	}
	Auth := ctx.handlers.Auth
	userPassword := strings.Split(string(content), "\n")
	userId := Auth(userPassword[0], userPassword[1])
	buf := "535 Authentication Failed\r\n"
	if "" != userId {
		ctx.User = userPassword[0]
		buf = "235 Authentication Successful\r\n"
		ctx.Login = true
		log.Println("auth by self")
	}
	ctx.Send(buf)
}

func quit(ctx *smtp_context_t) {
	ctx.End("221 2.0.0 " + ctx.conf.Domain + " Service closing transmission channel\r\n")
}

func xclient(ctx *smtp_context_t) {
	log.Println("auth by agency")
	ctx.Login = true
	ctx.hola()
}

func starttls(ctx *smtp_context_t) {
	ctx.Send("502 5.3.3 STARTTLS is not supported\r\n")
	log.Println("startTTS")
}

func help(ctx *smtp_context_t) {
	ctx.Send("502 5.3.3 HELP is not supported\r\n")
}

func noop(ctx *smtp_context_t) {
	ctx.Send("250 2.0.0 OK\r\n")
	log.Println("noop")
}

func rset(ctx *smtp_context_t) {
	ctx.Send("250 2.0.0 OK\r\n")
	log.Println("rset")
}

func mail(ctx *smtp_context_t) {
	ctx.Email.Sender = ctx.re.FindStringSubmatch(ctx.Msg)[1]
	clientDomain := strings.Split(ctx.Email.Sender, "@")[1]
	if (clientDomain == ctx.conf.Domain) != (!ctx.Login) { // 本域已登录 or 外域未登录
		ctx.Send("250 2.1.0 Sender <" + ctx.Email.Sender + "> OK\r\n")
		return
	}
	ctx.Send("530 5.7.1 Authentication Required\r\n")
}

func rcpt(ctx *smtp_context_t) {
	recver := ctx.re.FindStringSubmatch(ctx.Msg)[1]
	if strings.Split(recver, "@")[1] != ctx.conf.Domain && !ctx.Login { // 非登录用户 to 外域
		ctx.Send("530 5.7.1 Authentication Required\r\n")
		return
	}
	ctx.Email.Recver.PushBack(recver)
	ctx.Send("250 2.1.5 Recipient <" + recver + "> OK\r\n")
}

func data(ctx *smtp_context_t) {
	format := "from %s ([%s]) by %s over TLS secured channel with %s(%s)\r\n\t%d"
	ctx.Module = mod_HEAD
	config := ctx.conf
	ele := &KV{
		Name:  "Received",
		Value: fmt.Sprintf(format, config.Domain, config.Ip, config.Domain, config.Name, config.Version, time.Now().Unix()),
	}
	ctx.Email.Head = append(ctx.Email.Head, *ele)

	ctx.Send("354 Ok Send data ending with <CRLF>.<CRLF>\r\n")
}

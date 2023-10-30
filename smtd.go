// Simple Mail Transfer Daemon
package smtpd

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/watsonserve/maild"
)


// 接口集合
type SmtpServerConfigure interface {
	GetConfig() *maild.ServerConfig
	Auth(username string, password string) string
	TakeOff(email *maild.Mail)
}

type smtpd_t struct {
	config SmtpServerConfigure
	dict   map[string]func(*smtp_context_t)
}

func task(self *smtpd_t, conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	ctx := initSmtpContext(conn, self.config)
	ctx.hola()

	for scanner.Scan() {
		err := scanner.Err()
		if nil != err {
			log.Println("Reading standard input: ", err)
			break
		}
		// 行遍历
		msg := scanner.Text()
		ctx.Msg = msg
		switch ctx.Module {
		case mod_COMMAND:
			err = commandHash(self, ctx)
			if nil != err {
				log.Println("Unknow Cmd: ", err)
			}
		case mod_HEAD:
			dataHead(ctx)
		case mod_BODY:
			dataBody(ctx)
		}
	}
}

func commandHash(this *smtpd_t, ctx *smtp_context_t) error {
	var key string
	// 截取第一个单词
	_, err := fmt.Sscanf(ctx.Msg, "%s", &key)
	if nil != err {
		return err
	}
	// 查找处理方法
	method, exist := this.dict[key]
	if !exist {
		ctx.Send("method " + key + " not valid\r\n")
		return errors.New("method " + key + " not valid")
	}
	// 执行处理
	method(ctx)
	return nil
}

func dataHead(ctx *smtp_context_t) {
	if "" == ctx.Msg {
		ctx.Module = mod_BODY
	} else if ' ' == ctx.Msg[0] || '\t' == ctx.Msg[0] {
		ctx.Email.Head[len(ctx.Email.Head)-1].Value += "\r\n" + ctx.Msg
	} else {
		attr := strings.Split(ctx.Msg, ": ")
		ele := &maild.KV{
			Name:  attr[0],
			Value: attr[1],
		}
		ctx.Email.Head = append(ctx.Email.Head, *ele)
	}
}

func dataBody(ctx *smtp_context_t) {
	if "." == ctx.Msg {
		ctx.Module = mod_COMMAND
		ctx.Send("250 2.6.0 Message received\r\n")
		ctx.takeOff()
		return
	}
	ctx.Email.MailContent += ctx.Msg + "\r\n"
}

func Service(ln net.Listener, config SmtpServerConfigure) {
	if nil == config {
		return
	}
	that := &smtpd_t{
		config: config,
		dict: map[string]func(*smtp_context_t){
			"HELO":     helo,
			"EHLO":     ehlo,
			"AUTH":     auth,
			"QUIT":     quit,
			"XCLIENT":  xclient,
			"STARTTLS": starttls,
			"HELP":     help,
			"NOOP":     noop,
			"RSET":     rset,
			"MAIL":     mail,
			"RCPT":     rcpt,
			"DATA":     data,
		},
	}
	for {
		conn, err := ln.Accept()
		if nil != err {
			log.Println("a connect exception")
			continue
		}

		go task(that, conn)
	}
}

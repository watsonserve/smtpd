// Simple Mail Transfer Daemon
package smtpd

import (
	"bufio"
	"container/list"
    "crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type KV struct {
    Name string
    Value string
}

type Mail struct {
    Sender string
    Recver list.List
    Head []KV
    MailContent string
}

type ServerConfig struct {
	Domain  string
	Ip      string
	Name    string
	Type    string
	Version string
}

type SmtpServerConfigure interface {
	GetConfig() *ServerConfig
	Auth(username string, password string) string
	TakeOff(email *Mail)
}

type Smtpd struct {
	config SmtpServerConfigure
	dict   map[string]func(*smtp_context_t)
}

func TLSSocket(port string, crt string, key string) (net.Listener, error) {
    cert, err := tls.LoadX509KeyPair(crt, key)
    if nil != err {
        return nil, err
    }
    ln, err := tls.Listen("tcp", port, &tls.Config {
        Certificates: []tls.Certificate{cert},
        CipherSuites: []uint16 {
          tls.TLS_RSA_WITH_AES_256_CBC_SHA,
          tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
          tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
          tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
          tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
        },
        PreferServerCipherSuites: true,
    })
    if nil == err {
        defer ln.Close()
    }
    return ln, err
}

func Socket(port string) (net.Listener, error) {
    ln, err := net.Listen("tcp", port)
    if nil == err {
        defer ln.Close()
    }
    return ln, err
}

func New(config SmtpServerConfigure) *Smtpd {
	if nil == config {
		return nil
	}
	return &Smtpd {
		config: config,
		dict: map[string]func(*smtp_context_t) {
			"HELO": helo,
			"EHLO": ehlo,
			"AUTH": auth,
			"QUIT": quit,
			"XCLIENT": xclient,
			"STARTTLS": starttls,
			"HELP": help,
			"NOOP": noop,
			"RSET": rset,
			"MAIL": mail,
			"RCPT": rcpt,
			"DATA": data,
		},
	}
}

func (this *Smtpd) task(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	ctx := initSmtpContext(conn, this.config)
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
		case MOD_COMMAND:
			err = commandHash(this, ctx)
			if nil != err {
				log.Println("Unknow Cmd: ", err)
			}
		case MOD_HEAD:
			dataHead(ctx)
		case MOD_BODY:
			dataBody(ctx)
		}
	}
}

// 这里使用的是每个链接启动一个新的go程的模型，高并发的话，性能取决于go语言的协程能力
func (this *Smtpd) Listen(ln net.Listener) {
    for {
        conn, err := ln.Accept()
        if nil != err {
            log.Println("a connect exception")
            continue
        }
        defer conn.Close()
        go this.task(conn)
    }
}

func commandHash(this *Smtpd, ctx *smtp_context_t) error {
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
		ctx.Module = MOD_BODY
	} else if ' ' == ctx.Msg[0] || '\t' == ctx.Msg[0] {
		ctx.Email.Head[len(ctx.Email.Head)-1].Value += "\r\n" + ctx.Msg
	} else {
		attr := strings.Split(ctx.Msg, ": ")
		ele := &KV {
			Name:  attr[0],
			Value: attr[1],
		}
		ctx.Email.Head = append(ctx.Email.Head, *ele)
	}
}

func dataBody(ctx *smtp_context_t) {
	if "." == ctx.Msg {
		ctx.Module = MOD_COMMAND
		ctx.Send("250 2.6.0 Message received\r\n")
		ctx.takeOff()
		return
	}
	ctx.Email.MailContent += ctx.Msg + "\r\n"
}

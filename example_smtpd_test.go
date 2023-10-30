package smtpd_test

import (
	"fmt"
	"testing"

	"github.com/watsonserve/goutils"
	"github.com/watsonserve/smtpd"
	"github.com/watsonserve/maild"
)

type SmtpConf struct{}

func (sc *SmtpConf) GetConfig() *maild.ServerConfig {
	return &maild.ServerConfig{
		Domain:  "watsonserve.com",
		Ip:      "127.0.0.1",
		Type:    "SMTP",
		Name:    "WS_SMTPD",
		Version: "1.0",
	}
}

func (sc *SmtpConf) Auth(username string, password string) string {
	fmt.Printf("%s %s\n", username, password)
	return "null"
}

func (sc *SmtpConf) TakeOff(email *smtpd.Mail) {
	fmt.Println(email.Head, email.MailContent)
}

func TestExample(t *testing.T) {
	ln, err := goutils.Socket(":10025")
	if nil != err {
		t.Log(err)
	}

	t.Log("listen on port 10025")
	smtpd.Service(ln, &SmtpConf{})
}

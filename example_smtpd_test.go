package smtpd_test

import (
    "os"
    "io"
    "log"
    "fmt"
    "github.com/watsonserve/goutils"
    "github.com/watsonserve/smtpd"
)

type SmtpConf struct {}

func (this *SmtpConf) GetConfig() *smtpd.ServerConfig {
    return &smtpd.ServerConfig{
        Domain: "watsonserve.com",
        Ip: "127.0.0.1",
        Type: "SMTP",
        Name: "WS_SMTPD",
        Version: "1.0",
    }
}

func (this *SmtpConf) Auth(username string, password string) string {
    fmt.Printf("%s %s\n", username, password)
    return "null"
}

func (this *SmtpConf) TakeOff(email *smtpd.Mail) {
    fmt.Println(email.Head, email.MailContent)
}

func Example() {
    ln, err := goutils.Socket(":10025")
    if nil != err {
        log.Println(err)
    }

    log.Println("listen on port 10025")
    smtpd.Service(ln, &SmtpConf {})
}

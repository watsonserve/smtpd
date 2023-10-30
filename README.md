# smtpd

smtp server

使用者需实现`SmtpServerConfigure`接口

```
type SmtpServerConfigure interface {
    GetConfig() *lib.ServerConfig
    func (this *SmtpConf) Auth(username string, password string) string
    func TakeOff(email *smtpd.Mail)
}
```

准备妥当后只需构建一个 smtpd 实例并监听即可

```
smtpServer := smtpd.New(&SmtpConf {})

ln, err := smtpd.Socket(":10025")
if nil != err {
    log.Println(err)
}
smtpServer.Listen(ln)
```

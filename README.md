# smtpd
smtp server


# imapd
INTERNET MESSAGE ACCESS PROTOCOL - VERSION 4rev1 server

[RFC3501](http://www.faqs.org/rfcs/rfc3501.html)
```
typedef enum flags
{
    Seen = 1,
    Answered = 2,
    Flagged = 4,
    Deleted = 8,
    Draft = 16,
    Recent = 32
}
```
## 基础能力依赖
* 7层有状态代理：the IMAP client will keep a long TCP connect with server. So the server side need a proxy layer to handle socket connect from client and send real request with session_id to server
* 缓存状态机：客户端可订阅邮箱状态
* 语法分析器：IMAP协议的命令使用类似SQL的结构化查询语法
* 树形存储：IMAP协议要求可以建立目录，且目录下可以有子目录，所以需要类似LDAP或NFS的目录设计

client * n => proxy * 1
proxy * n  => server * 1
由于server上有事件发生时不知道对应的客户端在哪个proxy server上，所以应该由proxy去注册并正在客户端离开时清理资源

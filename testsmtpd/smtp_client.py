import socket
import ssl
import pprint
import time
import sys

def foo(self):
    if 3 == sys.version_info[0]:
        def sendStr(msg):
            return self.send(bytes(msg, encoding='utf-8'))
        return sendStr
    else:
        return self.send

def printB(buf):
    if 3 == sys.version_info[0]:
        print(str(buf, encoding='utf-8'))
    else:
        print(buf)

def tcp(server_name, port, use_ssl):
    sockfd = socket.socket(socket.AF_INET)
    if use_ssl:
        context = ssl.create_default_context()
        sockfd = context.wrap_socket(sockfd, server_hostname=server_name)
    sockfd.connect((server_name, port))
    return sockfd

if "__main__" == __name__:
    BUFSIZ = 4096
    # sockfd = tcp("imap-mail.outlook.com", 993, True)
    sockfd = tcp("127.0.0.1", 10025, False)
    sendStr = foo(sockfd)
    print("connected!")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)

    sendStr("EHLO WS_SMTP\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    sendStr("XCLIENT\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    sendStr("MAIL james<james@watsonserve.com>\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)

    sendStr("RCPT james<james@watsonserve.com>\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)

    sendStr("RCPT jameswatson<james@cn-bar.com>\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)

    sendStr("DATA\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    #sendStr("6 UID FETCH 126:127 (INTERNALDATE UID RFC822.SIZE FLAGS BODY.PEEK[HEADER.FIELDS (date subject from content-type to cc bcc message-id in-reply-to references)])\r\n")
    #buf = sockfd.recv(BUFSIZ)
    #printB(buf)
    sendStr("Title: test email\r\n")
    sendStr("\r\n\r\n")
    sendStr("abcdedghijklmn\r\n\r\n")
    sendStr("\r\n.\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    sendStr("QUIT\r\n")
    buf = sockfd.recv(BUFSIZ)
    printB(buf)
    time.sleep(10)
    #exit(0)

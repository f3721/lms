package email

import (
	"encoding/json"
	"errors"
	"go-admin/common/global"
	"go-admin/config"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/storage"
	"gopkg.in/gomail.v2"
)

const (
	smtpHost = "smtp.exmail.qq.com"
	smtpPort = 465
	from     = "ehsy-server@ehsy.com"
	password = "XIyu2018///"
)

const (
	InvalidEmail = "gomail: could not send email 1: 550 Mailbox not found or access denied"
)

type Email struct {
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
	AttachUrls []string `json:"attachUrls"`
}

// 同步发送邮件
func SendEmails(recipients []string, subject string, body string) error {

	// 单机模式不发邮件
	if config.ExtConfig.NoNetwork {
		return nil
	}

	// 连接到 SMTP 服务器
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// 执行发送
	for _, recipient := range recipients {
		if recipient == "" {
			continue
		}
		m := gomail.NewMessage()
		m.SetHeader("From", from)
		m.SetHeader("To", recipient)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)
		err := d.DialAndSend(m)
		// 单条，邮箱地址无效包装提醒
		if err != nil && len(recipients) == 1 && err.Error() == InvalidEmail {
			return errors.New("邮箱地址无效,请检查")
		}
	}

	return nil
}

// 异步发邮件
func AsyncSendEmails(recipients []string, subject string, body string) error {
	queueMap := map[string]any{
		"recipients": recipients,
		"subject":    subject,
		"body":       body,
	}

	// 推送到队列
	q := sdk.Runtime.GetMemoryQueue("")
	message, err := sdk.Runtime.GetStreamMessage("", global.AsyncEmail, queueMap)
	if err != nil {
		log.Printf("GetStreamMessage error, %s \n", err.Error())
	} else {
		err = q.Append(message)
		if err != nil {
			log.Printf("Append message error, %s \n", err.Error())
		}
	}

	return nil
}

// 异步处理邮件
func AsyncDealEmails(message storage.Messager) (err error) {
	// 解析数据
	values := message.GetValues()
	var rb []byte
	rb, err = json.Marshal(values)
	if err != nil {
		return err
	}
	var e Email
	err = json.Unmarshal(rb, &e)
	if err != nil {
		return err
	}

	// 执行发送
	err = SendEmails(e.Recipients, e.Subject, e.Body)
	if err != nil {
		return err
	}

	return nil
}

// 发送带附件的邮件
func SendEmailWithAttach(recipients []string, subject string, body string, attachUrls []string) error {

	// 单机模式不发邮件
	if config.ExtConfig.NoNetwork {
		return nil
	}

	// 连接到 SMTP 服务器
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// 执行发送
	for _, recipient := range recipients {
		if recipient == "" {
			continue
		}
		m := gomail.NewMessage()
		m.SetHeader("From", from)
		m.SetHeader("To", recipient)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)
		for _, attach := range attachUrls {
			// 获取文件内容
			resp, err := http.Get(attach)
			if err != nil {
				continue
			}
			defer resp.Body.Close()
			fileContent, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			// 获取文件名
			fileName, err := fileNameFromURL(attach)
			if err != nil {
				continue
			}

			// 添加附件
			m.Attach(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(fileContent)
				return err
			}))
		}
		err := d.DialAndSend(m)
		if err != nil {
			return err
		}
	}

	return nil
}

func fileNameFromURL(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	filename := filepath.Base(u.Path)
	return filename, nil
}

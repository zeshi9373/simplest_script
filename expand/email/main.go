package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"
)

// EmailConfig 邮件服务器配置
type EmailConfig struct {
	SMTPHost    string // SMTP 服务器地址
	SMTPPort    int    // 端口 (25, 465, 587)
	Username    string // 发件人邮箱，通常也是登录账号
	Password    string // 密码或授权码
	UseSSL      bool   // 是否使用 SSL (465端口通常为 true)
	StartTLS    bool   // 是否启用 STARTTLS (587端口通常为 true)
	InsecureTLS bool   // 是否跳过证书验证（测试环境）
}

// Message 邮件内容结构
type Message struct {
	From        string       // 发件人地址，可为空，默认使用 config.Username
	To          []string     // 收件人列表
	Cc          []string     // 抄送列表
	Bcc         []string     // 密送列表
	Subject     string       // 主题
	ContentType string       // "text/plain" 或 "text/html"，默认 "text/plain"
	Body        string       // 正文
	Attachments []Attachment // 附件列表
}

// Attachment 附件结构
type Attachment struct {
	Filename string // 附件文件名
	Data     []byte // 文件内容
	MimeType string // 可选，自动根据扩展名推断
}

// EmailSender 邮件发送器
type EmailSender struct {
	config *EmailConfig
	auth   smtp.Auth
}

// NewEmailSender 创建新的邮件发送器
func NewEmailSender(config *EmailConfig) *EmailSender {
	auth := smtp.PlainAuth("", config.Username, config.Password, config.SMTPHost)
	return &EmailSender{
		config: config,
		auth:   auth,
	}
}

// Send 发送邮件
func (s *EmailSender) Send(msg *Message) error {
	// 构建完整的邮件数据
	rawMsg, err := s.buildMessage(msg)
	if err != nil {
		return fmt.Errorf("构建邮件失败: %w", err)
	}

	// 确定服务器地址
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)

	// 根据配置选择发送方式
	if s.config.UseSSL {
		// 使用 SSL/TLS 直接连接 (465端口)
		return s.sendViaTLS(addr, msg, rawMsg)
	} else if s.config.StartTLS {
		// 使用 STARTTLS (587端口)
		return s.sendViaSTARTTLS(addr, msg, rawMsg)
	} else {
		// 普通不加密连接 (25端口)
		return smtp.SendMail(addr, s.auth, s.getFrom(msg), s.collectRecipients(msg), rawMsg)
	}
}

// sendViaTLS 使用 TLS 连接发送 (SMTP over SSL)
func (s *EmailSender) sendViaTLS(addr string, emailMsg *Message, msg []byte) error {
	tlsConfig := &tls.Config{
		ServerName:         s.config.SMTPHost,
		InsecureSkipVerify: s.config.InsecureTLS,
	}
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS 连接失败: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.SMTPHost)
	if err != nil {
		return fmt.Errorf("创建 SMTP 客户端失败: %w", err)
	}
	defer client.Quit()

	// 认证
	if err = client.Auth(s.auth); err != nil {
		return fmt.Errorf("认证失败: %w", err)
	}

	// 设置发件人
	from := s.getFrom(emailMsg)
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	// 收集所有收件人（To + Cc + Bcc）
	recipients := s.collectRecipients(emailMsg)
	for _, rcpt := range recipients {
		if err = client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("设置收件人 %s 失败: %w", rcpt, err)
		}
	}

	// 发送邮件数据
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取数据写入器失败: %w", err)
	}
	defer w.Close()

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("写入邮件数据失败: %w", err)
	}
	return nil
}

// sendViaSTARTTLS 使用 STARTTLS 升级连接
func (s *EmailSender) sendViaSTARTTLS(addr string, emailMsg *Message, msg []byte) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("连接 SMTP 失败: %w", err)
	}
	defer client.Quit()

	// 启动 STARTTLS
	tlsConfig := &tls.Config{
		ServerName:         s.config.SMTPHost,
		InsecureSkipVerify: s.config.InsecureTLS,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("STARTTLS 失败: %w", err)
	}

	// 认证
	if err = client.Auth(s.auth); err != nil {
		return fmt.Errorf("认证失败: %w", err)
	}

	// 设置发件人和收件人
	from := s.getFrom(emailMsg)
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}
	recipients := s.collectRecipients(emailMsg)
	for _, rcpt := range recipients {
		if err = client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("设置收件人 %s 失败: %w", rcpt, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取数据写入器失败: %w", err)
	}
	defer w.Close()
	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("写入邮件数据失败: %w", err)
	}
	return nil
}

// buildMessage 构建 MIME 邮件内容
func (s *EmailSender) buildMessage(msg *Message) ([]byte, error) {
	// 确定发件人
	from := s.getFrom(msg)

	// 边界字符串
	boundary := fmt.Sprintf("boundary_%d", time.Now().UnixNano())
	buf := &bytes.Buffer{}

	// 基础头部
	buf.WriteString(fmt.Sprintf("From: %s\r\n", from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ", ")))
	if len(msg.Cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(msg.Cc, ", ")))
	}
	// Bcc 不写入头部
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", encodeSubject(msg.Subject)))
	buf.WriteString("MIME-Version: 1.0\r\n")

	// 判断是否需要 multipart（有附件）
	if len(msg.Attachments) > 0 {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
		buf.WriteString("\r\n")

		// 正文部分
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		contentType := msg.ContentType
		if contentType == "" {
			contentType = "text/plain"
		}
		buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=UTF-8\r\n", contentType))
		buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		buf.WriteString(quotedprintableEncode(msg.Body))
		buf.WriteString("\r\n")

		// 附件部分
		for _, att := range msg.Attachments {
			buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			mimeType := att.MimeType
			if mimeType == "" {
				mimeType = mimeTypeByExtension(att.Filename)
			}
			buf.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", mimeType, att.Filename))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", att.Filename))
			buf.WriteString(base64Encode(att.Data))
			buf.WriteString("\r\n")
		}
		buf.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		// 无附件，简单正文
		contentType := msg.ContentType
		if contentType == "" {
			contentType = "text/plain"
		}
		buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=UTF-8\r\n", contentType))
		buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		buf.WriteString(quotedprintableEncode(msg.Body))
		buf.WriteString("\r\n")
	}

	return buf.Bytes(), nil
}

// getFrom 获取发件人地址
func (s *EmailSender) getFrom(msg *Message) string {
	if msg != nil && msg.From != "" {
		return msg.From
	}
	return s.config.Username
}

// collectRecipients 收集所有收件人（To + Cc + Bcc）
func (s *EmailSender) collectRecipients(msg *Message) []string {
	var recipients []string
	recipients = append(recipients, msg.To...)
	recipients = append(recipients, msg.Cc...)
	recipients = append(recipients, msg.Bcc...)
	return recipients
}

// encodeSubject 对中文主题进行 MIME 编码
func encodeSubject(subject string) string {
	// 如果包含非 ASCII 字符则编码
	for _, r := range subject {
		if r > 127 {
			return "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?="
		}
	}
	return subject
}

// quotedprintableEncode 将文本编码为 quoted-printable 格式
func quotedprintableEncode(s string) string {
	var buf bytes.Buffer
	w := quotedprintable.NewWriter(&buf)
	w.Write([]byte(s))
	w.Close()
	return buf.String()
}

// base64Encode 将字节数据编码为 base64 并换行
func base64Encode(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	// 每76字符换行
	var buf bytes.Buffer
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		buf.WriteString(encoded[i:end])
		buf.WriteString("\r\n")
	}
	return buf.String()
}

// mimeTypeByExtension 根据扩展名获取 MIME 类型（简单实现）
func mimeTypeByExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".html", ".htm":
		return "text/html"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}

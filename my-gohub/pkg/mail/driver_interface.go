package mail

type Driver interface {
	// Send 发送邮件
	Send(email Email, config map[string]string) bool
}

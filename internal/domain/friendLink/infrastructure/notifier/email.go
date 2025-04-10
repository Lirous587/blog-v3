package notifier

import (
	"blog/internal/domain/friendLink/model"
	"blog/pkg/config"
	"blog/pkg/email"
	"fmt"
	"github.com/pkg/errors"
)

type mailer struct {
	mailer email.Mailer
}

func NewMailer(m email.Mailer) Notifier {
	return &mailer{
		mailer: m,
	}
}

func (m *mailer) SendApprovalNotification(link *model.FriendLink) error {
	subject := "友链申请已通过通知"
	content := fmt.Sprintf(
		"<p>恭喜%s站长：</p>"+
			"<p>您提交的友链申请（站点：%s）已审核通过！</p>"+
			"<p>您的网站已被添加到我们的友链页面。</p>"+
			"<p>感谢您的支持！</p>",
		link.SiteName, link.Url)

	return errors.WithStack(m.mailer.SendHTML(link.Email, subject, content))
}

func (m *mailer) SendRejectionNotification(link *model.FriendLink) error {
	subject := "友链申请未通过通知"
	content := fmt.Sprintf(
		"<p>尊敬的%s站长：</p>"+
			"<p>很遗憾地通知您，您提交的友链申请（站点：%s）未能通过审核。</p>"+
			"<p>如有疑问，请回复此邮件与我们联系。</p>"+
			"<p>感谢您的理解！</p>",
		link.SiteName, link.Url)

	return errors.WithStack(m.mailer.SendHTML(link.Email, subject, content))
}

func (m *mailer) SendDeleteNotification(link *model.FriendLink, reason string) error {
	subject := "友链已被删除"
	content := fmt.Sprintf(
		"<p>尊敬的%s站长：</p>"+
			"<p>很遗憾地通知您，您的友链（站点：%s）因为以下原因已被删除。</p>"+
			"<br/> <b><p>%s</p></b> <br/>"+
			"<p>如有疑问，请回复此邮件与我们联系。</p>"+
			"<p>感谢您的理解！</p>",
		link.SiteName, link.Url, reason)

	return errors.WithStack(m.mailer.SendHTML(config.Cfg.Email.CC, subject, content))
}

func (m *mailer) SendPendingNotification(links []model.FriendLink) error {
	subject := "友链申请待审核"
	content := fmt.Sprintf(
		"<p>友链申请待审核：</p>"+
			"<p>当前共有%d条申请没有处理</p>",
		len(links))

	return errors.WithStack(m.mailer.SendHTML(config.Cfg.Email.CC, subject, content))
}

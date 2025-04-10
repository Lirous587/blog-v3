package notifier

import "blog/internal/domain/friendLink/model"

type Notifier interface {
	SendApprovalNotification(link *model.FriendLink) error
	SendRejectionNotification(link *model.FriendLink) error
	SendDeleteNotification(link *model.FriendLink, reason string) error
	SendPendingNotification(links []model.FriendLink) error
}

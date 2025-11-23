package common

const (
	ExchangeEmail       = "email.send"
	QueueNameAuthEmail  = "email.send.auth"
	RoutingKeyAuthEmail = "email.send.auth"

	ExchangeFile         = "file.action"
	QueueNameDeleteFile  = "file.action.delete"
	RoutingKeyDeleteFile = "file.action.delete"

	ExchangeNotification          = "notification.order"
	QueueNameServiceNotification  = "notification.order.service"
	RoutingKeyServiceNotification = "notification.order.service"

	RoleAdmin            = "admin"
	RoleAdminDisplayName = "Quản trị viên"
	RoleStaff            = "staff"
	RoleStaffDisplayName = "Nhân viên"
)

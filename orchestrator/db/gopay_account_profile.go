//go:build private_plugins

package db

type GoPayAccountProfile struct {
	GopayAccountID string `gorm:"primaryKey;column:gopay_account_id"`
	WAPhone        string `gorm:"column:wa_phone"`
	CreatedAt      int64  `gorm:"autoCreateTime"`
	UpdatedAt      int64  `gorm:"autoUpdateTime"`
}

func (GoPayAccountProfile) TableName() string {
	return "gopay_account_profiles"
}

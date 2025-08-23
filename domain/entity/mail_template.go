package entity

import (
	"queue-service/domain/common"
	"time"
)

type MailTemplate struct {
	tableName     struct{}      `pg:"mail_templates,alias:mt"`
	ID            string        `pg:"id,pk"`
	Name          string        `pg:"name"`
	Subject       string        `pg:"subject,unique"`
	Body          string        `pg:"body"`
	Keys          []string      `pg:"keys,array"`
	ProviderEmail string        `pg:"provider_email"`
	Status        common.Status `pg:"status"`
	CreatedBy     string        `pg:"created_by"`
	CreatedAt     time.Time     `pg:"created_at"`
	UpdatedAt     *time.Time    `pg:"updated_at"`
}

func (mt *MailTemplate) GetNameTable() any {
	return mt.tableName
}

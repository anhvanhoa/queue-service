package grpc_client

import "github.com/anhvanhoa/service-core/domain/oops"

var (
	ErrStatusHistoryClientNil = oops.New("Lịch sử email không tồn tại")
	ErrMailTemplateClientNil  = oops.New("Mẫu email không tồn tại")
	ErrMailProviderClientNil  = oops.New("Cấu hình gửi email không tồn tại")
	ErrMailHistoryClientNil   = oops.New("Lịch sử email không tồn tại")
)

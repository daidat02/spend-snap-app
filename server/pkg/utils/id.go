package utils

import "github.com/google/uuid"

// NewID sinh UUID v4 dạng string dùng làm khoá chính.
// DB cũng có DEFAULT gen_random_uuid()::text làm lớp bảo vệ phòng khi thiếu.
func NewID() string {
	return uuid.NewString()
}

package models

import "time"

type OperationLog struct {
	ID            int64
	NodeName      string
	NodeIP        string
	OperationType string
	Description   string
	Result        bool
	CreatedAt     time.Time
}

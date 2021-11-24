package models

import "time"

type UbisoftSession struct {
	Retries       uint8
	MaxRetries    uint8
	RetryTime     uint8
	SessionStart  time.Time
	SessionPeriod uint16
	SessionExpiry time.Time `json:"expiration"`
	SessionKey    string    `json:"sessionKey"`
	SpaceID       string    `json:"spaceId"`
	SessionTicket string    `json:"ticket"`
}

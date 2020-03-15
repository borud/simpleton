package model

import "time"

type Data struct {
	ID         int64     `db:"id" json:"id"`
	Timestamp  time.Time `db:"timestamp" json:"timestamp"`
	FromAddr   string    `db:"from_addr" json:"fromAddr"`
	PacketSize int       `db:"packet_size" json:"packetSize"`
	Payload    []byte    `db:"payload" json:"payload"`
}

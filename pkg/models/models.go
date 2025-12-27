package models

type Direction int

const (
	UP Direction = iota
	DOWN
	IDLE
)

type ExternalRequest struct {
	Direction Direction //Direction which person wants to go
	Floor     int64     //Floor on which the person is standing
}

type InternalRequest struct {
	Floor int64
}

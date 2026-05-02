package dto

// ActionType
// ENUM(kv_get, kv_set, kv_delete, kv_list)
//
//go:generate go tool go-enum --sql --marshal --names -f actions.go
type ActionType byte //nolint: recvcheck // required by JSON and SQL interfaces

type ActionKV struct {
	Type  ActionType `json:"type"            validate:"required"`
	Key   string     `json:"key"             validate:"required"`
	Value any        `json:"value,omitempty"`
}

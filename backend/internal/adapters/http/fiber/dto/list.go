package dto

import "palomnik/internal/ports"

type ListEnvelope[T any] struct {
	Data []T          `json:"data"`
	Meta ports.PageMeta `json:"meta"`
}

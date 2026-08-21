package dto

import "polomnik/internal/ports"

type ListEnvelope[T any] struct {
	Data []T          `json:"data"`
	Meta ports.PageMeta `json:"meta"`
}

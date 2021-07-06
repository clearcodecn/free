package model

import "time"

const (
	LinkTypeSSR   = "ssr"
	LinkTypeSS    = "ss"
	LinkTypeV2ray = "v2ray"
)

type Link struct {
	ID         string    `json:"id"`
	Host       string    `json:"host"`
	Port       string    `json:"port"`
	Raw        string    `json:"raw"`
	Type       string    `json:"type"`
	CreateTime time.Time `json:"createTime"`
}

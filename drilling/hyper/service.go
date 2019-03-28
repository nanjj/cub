package main

type Service struct {
	Address  string `json:"address"`
	Category string `json:"category"`
	Color    Color  `json:"color"`
}

//go:generate stringer -type=Color

type Color int

const (
	Red Color = iota
	Green
	Yellow
)

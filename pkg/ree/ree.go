package ree

import "time"

type ReePrices struct {
	Price    float64   `json:"value"`
	DateTime time.Time `json:"datetime"`
}
type ReeAttributes struct {
	Values []ReePrices `json:"values"`
}

type ReeIncluded struct {
	Type       string        `json:"type"`
	ID         string        `json:"id"`
	GroupID    string        `json:"groupId"`
	Attributes ReeAttributes `json:"attributes"`
}

type ReeError struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}
type ReeResponse struct {
	Included []ReeIncluded `json:"included"`
}

type ReeErrorResponse struct {
	Errors []ReeError `json:"errors"`
}

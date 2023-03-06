package main

import "testing"

func TestFillNewFileWith(t *testing.T) {
	fillNewFileWith(map[string]string{
		"message": "this is a test message...",
	})
}

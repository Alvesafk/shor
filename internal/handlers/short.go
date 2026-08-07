/*
SPDX-License-Identifier: GPL-3.0-only
Copyright (c) 2026 Alvesafk.
*/
package handlers

import (
	"crypto/rand"
	"log"
	"math/big"
	"strings"
)

const (
	defaultMaxShortURLSize = 5
	chars                  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateShortURL() (string, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered from panic:", err)
		}
	}()

	var sb strings.Builder
	for range defaultMaxShortURLSize {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars)-1)))
		if err != nil {
			return "", err
		}

		sb.WriteByte(chars[index.Int64()])
	}

	return sb.String(), nil
}

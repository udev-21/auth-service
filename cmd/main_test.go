package main_test

import (
	"testing"

	"golang.org/x/crypto/argon2"
)

var salt = []byte("salt")

func BenchmarkMain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		argon2.IDKey([]byte("$ecuRePa$$w0rd"), salt, 1, 15*1024, 2, 32)
	}
}

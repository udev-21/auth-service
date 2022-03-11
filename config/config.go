package config

import "time"

type Config struct {
	JWT            JWTConfig
	PasswordConfig PasswordConfig
}

type JWTConfig struct {
	SecretKey                  []byte
	AccessTokenExpireDuration  time.Duration
	RefreshTokenExpireDuration time.Duration
}

// look for owasp recommendations: https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
type ArgonEncodeConfig struct {
	Time      uint32 // should be between 1 and 3
	Memory    uint32 // should be between 16*1024
	Threads   uint8  // should be 1
	KeyLength uint32 // should be 64
}

type PasswordConfig struct {
	Salt  []byte
	Argon ArgonEncodeConfig
}

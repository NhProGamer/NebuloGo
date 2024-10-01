package salt

import (
	"NebuloGo/config"
	"encoding/hex"
	"golang.org/x/crypto/argon2"
)

func HashPhrase(phrase string) string {
	salt := []byte(config.Configuration.Argon.Salt)
	iterations := config.Configuration.Argon.Iterations
	memory := config.Configuration.Argon.Memory
	parallelism := config.Configuration.Argon.Parallelism
	hashlenght := config.Configuration.Argon.HashLenght
	return hex.EncodeToString(argon2.IDKey([]byte(phrase), salt, iterations, memory, parallelism, hashlenght))
}

func HashCompare(phrase string, hash string) bool {
	return HashPhrase(phrase) == hash
}

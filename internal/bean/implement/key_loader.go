package beanimplement

import (
	"fmt"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"os"
)

type KeyLoader struct {
	RSAKey []byte
	PGPKey []byte
}

func NewKeyLoader(path model.KeyPath) *KeyLoader {
	rsaKey, err := os.ReadFile(path.RSAKey)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	pgpKey, err := os.ReadFile(path.PGPKey)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &KeyLoader{
		RSAKey: rsaKey,
		PGPKey: pgpKey,
	}
}

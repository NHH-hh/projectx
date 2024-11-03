package crypto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyPairSignVerifySuccess(t *testing.T) {
	key := GeneratePrivateKey()
	pubKey := key.PublicKey()
	// address := pubKey.Address()
	msg := []byte("hello world")
	sign, err := key.Sign(msg)
	assert.Nil(t, err)
	assert.True(t, sign.Verify(pubKey, msg))
	fmt.Println(sign)
}

func TestKeyPairSignVerifyFail(t *testing.T) {
	priKey := GeneratePrivateKey()
	pubKey := priKey.PublicKey()
	// address := pubKey.Address()
	msg := []byte("hello world")
	sign, err := priKey.Sign(msg)
	assert.Nil(t, err)
	otherPriKey := GeneratePrivateKey()
	otherPubKey := otherPriKey.PublicKey()
	assert.False(t, sign.Verify(otherPubKey, msg))
	assert.False(t, sign.Verify(pubKey, []byte("xxx")))
}

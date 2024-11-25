package core

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"projectx/crypto"
	"testing"
)

func TestSignTransaction(t *testing.T) {
	data := []byte("hello world")
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: data,
	}
	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	data := []byte("hello world")
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: data,
	}
	assert.Nil(t, tx.Sign(privKey))
	assert.Nil(t, tx.Verify())
	otherPrivKey := crypto.GeneratePrivateKey()
	tx.From = otherPrivKey.PublicKey()
	assert.NotNil(t, tx.Verify())
}

func TestTxEncodeAndDecode(t *testing.T) {
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t, tx.Encode(NewGolTxEncoder(buf)))
	txDecode := Transaction{}
	assert.Nil(t, txDecode.Decode(NewGobTxDecoder(buf)))
	assert.Equal(t, tx, txDecode)
}

func randomTxWithSignature(t *testing.T) Transaction {
	tx := Transaction{
		Data: []byte("hello world"),
	}
	assert.Nil(t, tx.Sign(crypto.GeneratePrivateKey()))
	return tx
}

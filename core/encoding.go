package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/gob"
	"io"
	"projectx/crypto"
	"projectx/types"
)

type Encoder[T any] interface {
	Encode(T) error
}

type Decoder[T any] interface {
	Decode(T) error
}

type GobTxEncoder struct {
	w io.Writer
}

func NewGolTxEncoder(w io.Writer) *GobTxEncoder {
	gob.Register(elliptic.P256())
	return &GobTxEncoder{w: w}
}

func (e *GobTxEncoder) Encode(tx *Transaction) error {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(tx.From.Key)
	if err != nil {
		return err
	}
	x := &TxSerializable{
		Data:      tx.Data,
		FromBytes: pubKeyBytes,
		Signature: tx.Signature,
		hash:      tx.hash,
		firstSeen: tx.firstSeen,
	}
	return gob.NewEncoder(e.w).Encode(x)
}

type GobTxDecoder struct {
	r io.Reader
}

func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	gob.Register(elliptic.P256())
	return &GobTxDecoder{r: r}
}

func (d *GobTxDecoder) Decode(tx *Transaction) error {
	x := &TxSerializable{}
	err := gob.NewDecoder(d.r).Decode(x)
	if err != nil {
		return err
	}
	key, err := x509.ParsePKIXPublicKey(x.FromBytes)
	if err != nil {
		return err
	}
	from := crypto.PublicKey{
		Key: key.(*ecdsa.PublicKey),
	}
	tx.Data = x.Data
	tx.From = from
	tx.Signature = x.Signature
	tx.hash = x.hash
	tx.firstSeen = x.firstSeen
	return nil
}

type TxSerializable struct {
	Data      []byte
	FromBytes []byte
	Signature *crypto.Signature
	// cached version of the tx data hash
	hash types.Hash
	// firstSeen is timestamp of when this tx is first seen locally
	firstSeen int64
}

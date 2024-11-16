package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"projectx/crypto"
	"projectx/types"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
	Nonce         uint64
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)
	return buf.Bytes()
}

//func (h *Header) EncodeBinary(w io.Writer) error {
//	if err := binary.Write(w, binary.LittleEndian, &h.Version); err != nil {
//		return err
//	}
//	if err := binary.Write(w, binary.LittleEndian, &h.PrevBlock); err != nil {
//		return err
//	}
//	if err := binary.Write(w, binary.LittleEndian, &h.Timestamp); err != nil {
//		return err
//	}
//	if err := binary.Write(w, binary.LittleEndian, &h.Height); err != nil {
//		return err
//	}
//	return binary.Write(w, binary.LittleEndian, &h.Nonce)
//}
//
//func (h *Header) DecodeBinary(r io.Reader) error {
//	if err := binary.Read(r, binary.LittleEndian, &h.Version); err != nil {
//		return err
//	}
//	if err := binary.Read(r, binary.LittleEndian, &h.PrevBlock); err != nil {
//		return err
//	}
//	if err := binary.Read(r, binary.LittleEndian, &h.Timestamp); err != nil {
//		return err
//	}
//	if err := binary.Read(r, binary.LittleEndian, &h.Height); err != nil {
//		return err
//	}
//	return binary.Read(r, binary.LittleEndian, &h.Nonce)
//}

type Block struct {
	*Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature
	// cached version of the header hash
	hash types.Hash
}

func NewBlock(h *Header, txs []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txs,
	}
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, *tx)
}

func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sign, err := privKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}
	b.Validator = privKey.PublicKey()
	b.Signature = sign
	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return errors.New("block has no signature")
	}
	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return errors.New("block has invalid signature")
	}
	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}
	dataHash, err := CalculateDataHash(b.Transactions)
	if err != nil {
		return err
	}
	if dataHash != b.DataHash {
		return fmt.Errorf("block(%s) has invalid data hash", b.Hash(BlockHasher{}))
	}
	return nil
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}
	return b.hash
}

func CalculateDataHash(txx []Transaction) (hash types.Hash, err error) {
	buf := &bytes.Buffer{}
	for _, tx := range txx {
		if err = tx.Encode(NewGolTxEncoder(buf)); err != nil {
			return
		}
	}
	hash = sha256.Sum256(buf.Bytes())
	return
}

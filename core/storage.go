package core

type Storage interface {
	Put(*Block) error
	//Get(*Block) ([]byte, error)
}

type MemoryStorage struct{}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (ms *MemoryStorage) Put(b *Block) error {
	return nil
}

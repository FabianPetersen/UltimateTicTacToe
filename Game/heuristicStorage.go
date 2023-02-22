package Game

type Storage struct {
	heuristicStore map[[9]uint32]float64
}

func (storage *Storage) Count() int {
	return len(storage.heuristicStore)
}

func (storage *Storage) Get(hash GameHash) (float64, bool) {
	node, exists := storage.heuristicStore[*hash]
	return node, exists
}

func (storage *Storage) Set(hash GameHash, score float64) {
	storage.heuristicStore[*hash] = score
}

func (storage *Storage) Reset() {
	storage.heuristicStore = make(map[[9]uint32]float64, 50000)
}

func NewStorage() Storage {
	return Storage{
		heuristicStore: make(map[[9]uint32]float64, 50000),
	}
}

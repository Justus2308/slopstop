package main

import (
	"fmt"
	"io"
	"math/rand"
)

func main() {
	fmt.Printf("Hello, world!\n")
}

type Hasher struct {
	vocab       map[string]uint
	k           int
	hashVecs    []map[uint]uint
	isFinalized bool
}

func NewHasher(k int) (*Hasher, error) {
	if k < 1 || k > 128 {
		panic("k must be at least 1 and at most 128")
	}
	hasher := new(Hasher)
	hasher.vocab = make(map[string]uint, 128)
	hasher.k = k
	hasher.isFinalized = false
	return hasher, nil
}

func (hasher *Hasher) AddShingles(r io.Reader) (int, error) {
	if hasher.isFinalized {
		panic("cannot add shingles to finalized hasher")
	}
	buf := make([]byte, 1024)
	bytes_read, err := io.ReadFull(r, buf[:hasher.k])
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return 0, nil
		} else {
			return 0, err
		}
	}
	if bytes_read < hasher.k {
		return 0, nil
	}
	old_len := len(hasher.vocab)
	for {
		bytes_read, err = io.ReadFull(r, buf[hasher.k:])
		if err != nil {
			if err == io.EOF {
				s := buf[:hasher.k]
				hasher.vocab[string(s)] = 0
				hasher_added := (len(hasher.vocab) - old_len)
				return hasher_added, nil
			} else if err != io.ErrUnexpectedEOF {
				hasher_added := (len(hasher.vocab) - old_len)
				return hasher_added, err
			}
		}

		for i := 0; i < bytes_read; i += 1 {
			s := buf[i:][:hasher.k]
			hasher.vocab[string(s)] = 0
		}
		_ = copy(buf[:hasher.k], buf[bytes_read:][:hasher.k])
	}
}

func (hasher *Hasher) Finalize(vectorCount uint) {
	var idx uint = 0
	for s := range hasher.vocab {
		hasher.vocab[s] = idx
		idx++
	}
	hasher.hashVecs = make([]map[uint]uint, vectorCount)
	for i := range vectorCount {
		hasher.hashVecs[i] = hasher.makeHashVector()
	}
	hasher.isFinalized = true
}
func (hasher Hasher) makeHashVector() map[uint]uint {
	vector := make(map[uint]uint, len(hasher.vocab))
	for i := range len(hasher.vocab) {
		vector[uint(i+1)] = uint(i)
	}
	rand.Shuffle(len(vector), func(i, j int) {
		tmp, _ := vector[uint(i+1)]
		vector[uint(i+1)], _ = vector[uint(j+1)]
		vector[uint(j+1)] = tmp
	})
	return vector
}

func (hasher Hasher) Hash(str string) map[uint]struct{} {
	if !hasher.isFinalized {
		panic("cannot hash based on non-finalized hasher")
	}
	signature := make(map[uint]struct{}, len(hasher.vocab))
	vector := hasher.oneHot(str)
	for _, hvec := range hasher.hashVecs {
		for i := range len(hasher.vocab) {
			idx := hvec[uint(i+1)]
			signature_val := vector[idx]
			if signature_val {
				signature[idx] = struct{}{}
				break
			}
		}
	}
	return signature
}
func (hasher Hasher) oneHot(str string) []bool {
	// guaranteed to be zero-initialized (all false)
	vector := make([]bool, len(hasher.vocab))
	for i := range len(str) - hasher.k + 1 {
		s := str[i:][:hasher.k]
		idx, ok := hasher.vocab[s]
		if ok {
			vector[idx] = true
		}
	}
	// fmt.Printf("%v", vector)
	return vector
}

func Similiarity(a, b map[uint]struct{}) float64 {
	union := a
	sect := make(map[uint]struct{}, (len(a) / 2))
	for k := range b {
		_, ok := a[k]
		if ok {
			sect[k] = struct{}{}
		} else {
			union[k] = struct{}{}
		}
	}
	return float64(len(sect)) / float64(len(union))
}

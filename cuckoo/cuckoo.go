package cuckoo

import (
	"crypto/sha1"
	"math/rand"
)

const BucketSize = 4
const TableSize = 256

// kick次數超過maxNumKicks就視為insert失敗
const maxNumKicks = 256

// 0代表該slot沒有存data
type Bucket [BucketSize]byte

type HashTable struct {
	buckets [TableSize]Bucket
}

//hash完結果不可以為0
func getFingerPrint(data string) byte {
	h := sha1.New()
	h.Write([]byte(data))
	return h.Sum(nil)[0]
}

func getHashedFingerPrint(f byte) byte {
	fSlice := make([]byte, 1)
	fSlice[0] = f

	h := sha1.New()
	h.Write(fSlice)
	return h.Sum(nil)[0]
}

func getInsertLocations(data string, f byte) (byte, byte) {
	h := sha1.New()
	h.Write([]byte(data))

	// take the last byte of sha1 as i1
	i1 := h.Sum(nil)[0]

	i2 := i1 ^ getHashedFingerPrint(f)

	return i1, i2
}

func getAltInsertLocation(i byte, f byte) byte {
	return i ^ getHashedFingerPrint(f)
}

func (h *HashTable) Insert(data string) bool {
	f := getFingerPrint(data)
	i1, i2 := getInsertLocations(data, f)

	for i := 0; i < BucketSize; i++ {
		if h.buckets[i1][i] == 0 {
			h.buckets[i1][i] = f
			return true
		}
	}

	for i := 0; i < BucketSize; i++ {
		if h.buckets[i2][i] == 0 {
			h.buckets[i2][i] = f
			return true
		}
	}

	// both buckets[i1] and buckets[i2] are full
	return h.reInsert(f, i1)
}

func (h *HashTable) reInsert(f byte, i byte) bool {
	for k := 0; k < maxNumKicks; k++ {
		for j := 0; j < BucketSize; j++ {
			if h.buckets[i][j] == 0 {
				h.buckets[i][j] = f
				return true
			}
		}

		// kick a fingerprint out randomly
		index := rand.Intn(BucketSize)

		temp := f
		f = h.buckets[i][index]
		h.buckets[i][index] = temp

		// i xor hash(fingerPrint)
		hf := getHashedFingerPrint(f)
		i = getAltInsertLocation(i, hf)
	}

	return false
}

func (h *HashTable) Lookup(data string) bool {
	f := getFingerPrint(data)
	i1, i2 := getInsertLocations(data, f)

	for i := 0; i < BucketSize; i++ {
		if h.buckets[i1][i] == f || h.buckets[i2][i] == f {
			return true
		}
	}

	return false
}

func (h *HashTable) Delete(data string) bool {
	f := getFingerPrint(data)
	i1, i2 := getInsertLocations(data, f)

	for i := 0; i < BucketSize; i++ {
		if h.buckets[i1][i] == f {
			h.buckets[i1][i] = 0
			return true
		}
		if h.buckets[i2][i] == f {
			h.buckets[i2][i] = 0
			return true
		}
	}

	return false
}

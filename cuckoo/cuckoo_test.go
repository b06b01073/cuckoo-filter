package cuckoo

import (
	"crypto/sha1"
	"testing"
)

const testKey = "cat"

func TestGetFingerPrint(t *testing.T) {

	h := sha1.New()
	h.Write([]byte(testKey))
	if getFingerPrint(testKey) != h.Sum(nil)[0] {
		t.Error("getFingerPrint failed")
	}
}

func TestGetHashedFingerPrint(t *testing.T) {
	f := getFingerPrint(testKey)
	h := sha1.New()
	fSlice := make([]byte, 1)
	fSlice[0] = f
	h.Write(fSlice)

	if getHashedFingerPrint(f) != h.Sum(nil)[0] {
		t.Error("getHashedFingerPrint failed")
	}
}

func TestGetInsertLocations(t *testing.T) {
	f := getFingerPrint(testKey)
	hf := getHashedFingerPrint(f)
	i1, i2 := getInsertLocations(testKey, f)
	if i1^hf != i2 || i2^hf != i1 {
		t.Error("incorrect alternative location")
	}
}

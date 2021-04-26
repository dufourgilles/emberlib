package embertree_test

import (
	//"fmt"
	"testing"

	"github.com/dufourgilles/emberlib/asn1"
	. "github.com/dufourgilles/emberlib/embertree"
)

func TestDecodeTarget(t *testing.T) {
	buffer := []byte{163, 29, 48, 27, 160, 7, 110, 5, 160, 3, 2, 1, 1, 160, 7, 110, 5, 160, 3, 2, 1, 3, 160, 7, 110, 5, 160, 3, 2, 1, 5}
	encodedTargets := []int32{1, 3, 5}
	reader := asn1.NewASNReader(buffer)
	matrix, err := NewMatrix(1, OneToN, Linear)
	if err != nil {
		t.Error(err)
		return
	}
	err = matrix.DecodeTargets(reader)
	if err != nil {
		t.Error(err)
		return
	}
	targets, err := matrix.GetTargets()
	if err != nil {
		t.Error(err)
		return
	}
	for i, signal := range targets {
		target := signal.(*Target)
		if target.Number != encodedTargets[i] {
			t.Errorf("Target mismatch at position %d. Got %d instead of %d", i, target.Number, encodedTargets[i])
		}
	}
}

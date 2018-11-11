/*
 * Copyright (c) 2018 wellwell.work, LLC by Zoe
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *
 */

package pool

import (
	"math/rand"
	"testing"
	"time"
)

func TestIndex(t *testing.T) {
	testIndex(t, 0, 0)
	testIndex(t, 1, 0)

	testIndex(t, minSize-1, 0)
	testIndex(t, minSize, 0)
	testIndex(t, minSize+1, 1)

	testIndex(t, 2*minSize-1, 1)
	testIndex(t, 2*minSize, 1)
	testIndex(t, 2*minSize+1, 2)

	testIndex(t, maxSize-1, steps-1)
	testIndex(t, maxSize, steps-1)
	testIndex(t, maxSize+1, steps-1)
}

func testIndex(t *testing.T, n, expectedIdx int) {
	idx := index(n)
	if idx != expectedIdx {
		t.Fatalf("unexpected idx for n=%d: %d. Expecting %d", n, idx, expectedIdx)
	}
}

func TestPoolCalibrate(t *testing.T) {
	for i := 0; i < steps*calibrateCallsThreshold; i++ {
		n := 1004
		if i%15 == 0 {
			n = rand.Intn(15234)
		}
		testGetPut(t, n)
	}
}

func TestPoolVariousSizesSerial(t *testing.T) {
	testPoolVariousSizes(t)
}

func TestPoolVariousSizesConcurrent(t *testing.T) {
	concurrency := 5
	ch := make(chan struct{})
	for i := 0; i < concurrency; i++ {
		go func() {
			testPoolVariousSizes(t)
			ch <- struct{}{}
		}()
	}
	for i := 0; i < concurrency; i++ {
		select {
		case <-ch:
		case <-time.After(3 * time.Second):
			t.Fatalf("timeout")
		}
	}
}

func testPoolVariousSizes(t *testing.T) {
	for i := 0; i < steps+1; i++ {
		n := (1 << uint32(i))

		testGetPut(t, n)
		testGetPut(t, n+1)
		testGetPut(t, n-1)

		for j := 0; j < 10; j++ {
			testGetPut(t, j+n)
		}
	}
}

func testGetPut(t *testing.T, n int) {
	bb := ByteBufferGet()
	if len(bb.B) > 0 {
		t.Fatalf("non-empty byte buffer returned from acquire")
	}
	bb.B = allocNBytes(bb.B, n)
	ByteBufferPut(bb)
}

func allocNBytes(dst []byte, n int) []byte {
	diff := n - cap(dst)
	if diff <= 0 {
		return dst[:n]
	}
	return append(dst, make([]byte, diff)...)
}

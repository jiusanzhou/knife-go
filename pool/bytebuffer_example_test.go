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

package pool_test

import (
	"fmt"

	"github.com/jiusanzhou/knife-go/pool"
)

func ExampleByteBuffer() {
	bb := pool.ByteBufferGet()

	bb.WriteString("first line\n")
	bb.Write([]byte("second line\n"))
	bb.B = append(bb.B, "third line\n"...)

	fmt.Printf("bytebuffer contents=%q", bb.B)

	// It is safe to release byte buffer now, since it is
	// no longer used.
	pool.ByteBufferPut(bb)
}

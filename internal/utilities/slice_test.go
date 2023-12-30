// Copyright (c) 2022-2023 OpenSeaWaves.com/Rasbora
//
// This file is part of Rasbora Distributed Video Transcoding
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package utilities

import (
	"testing"
)

func TestInSlice_Int(t *testing.T) {
	var slice = []int{1, 2, 3}
	var value = 1

	if InSlice(value, slice) == false {
		t.Errorf("expected: %v", true)
	}
}

func TestInSlice_String(t *testing.T) {
	var slice = []string{"rasbora", "transcoder", "video"}
	var value = "video"

	if InSlice(value, slice) == false {
		t.Errorf("expected: %v", true)
	}
}

func TestInSlice_Float64(t *testing.T) {
	var slice = []float64{0.1, 0.2, 0.3}
	var value = 0.1

	if InSlice(value, slice) == false {
		t.Errorf("expected: %v", true)
	}
}

func BenchmarkInSlice(b *testing.B) {
	var slice = []int{1, 2, 3}
	for i := 0; i < b.N; i++ {
		InSlice(i, slice)
	}
}

func ExampleInSlice() {
	var slice = []int{1, 2, 3}
	var value = 1

	if InSlice(value, slice) {
		//Do something
	}
}

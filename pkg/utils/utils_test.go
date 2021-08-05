/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ByteArray2String(t *testing.T) {
	str := ByteArray2String([]byte("abc"))

	assert.Equal(t, "abc", str)
}

func Test_String2IntPointer(t *testing.T) {
	i, err := String2IntPointer("0")
	assert.Nil(t, err)

	want := int(0)
	assert.Equal(t, &want, i)
}

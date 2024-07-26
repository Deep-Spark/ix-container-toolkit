/*
# Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
# Copyright (c) 2021, NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/

package discover

import (
	"path/filepath"
	"sync"

	"gitee.com/deep-spark/ix-container-runtime/internal/lookup"
)

// mounts is a generic discoverer for Mounts. It is customized by specifying the
// required entities as a list and a Locator that is used to find the target mounts
// based on the entry in the list.
type mounts struct {
	None
	lookup   lookup.Locator
	root     string
	required []string
	sync.Mutex
	cache []Mount
}

var _ Discover = (*mounts)(nil)

// NewMounts creates a discoverer for the required mounts using the specified locator.
func NewMounts(lookup lookup.Locator, root string, required []string) Discover {
	return newMounts(lookup, root, required)
}

// newMounts creates a discoverer for the required mounts using the specified locator.
func newMounts(lookup lookup.Locator, root string, required []string) *mounts {
	return &mounts{
		lookup:   lookup,
		root:     filepath.Join("/", root),
		required: required,
	}
}

func (d *mounts) Mounts() ([]Mount, error) {
	return nil, nil
}

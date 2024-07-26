/**
# Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
# Copyright (c) NVIDIA CORPORATION.  All rights reserved.
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
**/

package lookup

import (
	"fmt"
	"os"
)

type executable struct {
	file
}

// NewExecutableLocator creates a locator to fine executable files in the path. A logger can also be specified.
func NewExecutableLocator(root string) Locator {
	paths := GetPaths(root)

	return newExecutableLocator(root, paths...)
}

func newExecutableLocator(root string, paths ...string) *executable {
	f := newFileLocator(
		WithRoot(root),
		WithSearchPaths(paths...),
		WithFilter(assertExecutable),
		WithCount(1),
	)

	l := executable{
		file: *f,
	}

	return &l
}

// assertExecutable checks whether the specified path is an execuable file.
func assertExecutable(filename string) error {
	err := assertFile(filename)
	if err != nil {
		return err
	}
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}

	if info.Mode()&0111 == 0 {
		return fmt.Errorf("specified file '%v' is not executable", filename)
	}

	return nil
}

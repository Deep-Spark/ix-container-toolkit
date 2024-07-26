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

package oci

import (
	"fmt"

	"github.com/opencontainers/runtime-spec/specs-go"
)

// SpecModifier defines an interface for modifying a (raw) OCI spec
type SpecModifier interface {
	// Modify is a method that accepts a pointer to an OCI Spec and returns an
	// error. The intention is that the function would modify the spec in-place.
	Modify(*specs.Spec) error
}

type Spec interface {
	Load() (*specs.Spec, error)
	Flush() error
	Modify(SpecModifier) error
	LookupEnv(string) (string, bool)
}

func NewSpec(args []string) (Spec, error) {
	bundleDir, err := GetBundleDir(args)
	if err != nil {
		return nil, fmt.Errorf("error getting bundle directory: %v", err)
	}

	ociSpecPath := GetSpecFilePath(bundleDir)

	ociSpec := NewFileSpec(ociSpecPath)

	return ociSpec, nil
}

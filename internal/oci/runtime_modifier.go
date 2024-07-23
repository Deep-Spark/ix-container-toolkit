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

	log "github.com/sirupsen/logrus"
)

type modifyingRuntimeWrapper struct {
	runtime  Runtime
	ociSpec  Spec
	modifier SpecModifier
}

func NewModifyingRuntimeWrapper(runtime Runtime, spec Spec, modifier SpecModifier) Runtime {
	if modifier == nil {
		return runtime
	}

	rt := modifyingRuntimeWrapper{
		runtime:  runtime,
		ociSpec:  spec,
		modifier: modifier,
	}
	return &rt
}

func (r *modifyingRuntimeWrapper) Exec(args []string) error {
	if HasCreateSubcommand(args) {
		err := r.modify()
		if err != nil {
			return fmt.Errorf("could not apply required modification to OCI specification: %v", err)
		}
		log.Printf("Applied required modification to OCI specification")
	} else {
		log.Printf("No modification of OCI specification required")
	}

	return r.runtime.Exec(args)
}

// modify loads, modifies, and flushes the OCI specification using the defined Modifier
func (r *modifyingRuntimeWrapper) modify() error {
	_, err := r.ociSpec.Load()
	if err != nil {
		return fmt.Errorf("error loading OCI specification for modification: %v", err)
	}

	err = r.ociSpec.Modify(r.modifier)
	if err != nil {
		return fmt.Errorf("error modifying OCI spec: %v", err)
	}

	err = r.ociSpec.Flush()
	if err != nil {
		return fmt.Errorf("error writing modified OCI specification: %v", err)
	}
	return nil
}

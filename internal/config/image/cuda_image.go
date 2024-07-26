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

package image

import (
	"strings"

	"gitee.com/deep-spark/ix-container-runtime/internal/config"
	"github.com/opencontainers/runtime-spec/specs-go"
)

// CUDA represents a CUDA image that can be used for GPU computing. This wraps
// a map of environment variable to values that can be used to perform lookups
// such as requirements.
type CUDA struct {
	env    map[string]string
	mounts []specs.Mount
	Cfg    *config.Config
}

// NewCUDAImageFromSpec creates a CUDA image from the input OCI runtime spec.
// The process environment is read (if present) to construc the CUDA Image.
func NewCUDAImageFromSpec(spec *specs.Spec, cfg *config.Config) (CUDA, error) {
	var env []string
	if spec != nil && spec.Process != nil {
		env = spec.Process.Env
	}

	return New(
		WithEnv(env),
		WithMounts(spec.Mounts),
		WithConfig(cfg),
	)
}

func (i CUDA) DevicesFromEnvvars(envVars ...string) VisibleDevices {
	// We concantenate all the devices from the specified env.
	var isSet bool
	var devices []string
	requested := make(map[string]bool)
	for _, envVar := range envVars {
		if devs, ok := i.env[envVar]; ok {
			isSet = true
			for _, d := range strings.Split(devs, ",") {
				trimmed := strings.TrimSpace(d)
				if len(trimmed) == 0 {
					continue
				}
				devices = append(devices, trimmed)
				requested[trimmed] = true
			}
		}
	}

	// Environment variable unset with legacy image: default to "all".
	if !isSet && len(devices) == 0 {
		return NewVisibleDevices("all")
	}

	// Environment variable unset or empty or "void": return nil
	if len(devices) == 0 || requested["void"] {
		return NewVisibleDevices("void")
	}

	return NewVisibleDevices(devices...)
}

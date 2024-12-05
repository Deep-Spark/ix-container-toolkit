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

package runtime

import (
	"os"

	log "github.com/sirupsen/logrus"

	"gitee.com/deep-spark/ix-container-runtime/internal/config"
	"gitee.com/deep-spark/ix-container-runtime/internal/config/image"
	"gitee.com/deep-spark/ix-container-runtime/internal/modifier"
	"gitee.com/deep-spark/ix-container-runtime/internal/oci"
)

var (
	defaultRuntime = []string{"docker-runc", "runc", "crun"}
)

func (r rt) Run(argv []string) (rerr error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	lowLevelRuntime, err := oci.NewLowLevelRuntime(defaultRuntime)

	if !oci.HasCreateSubcommand(argv) {
		return lowLevelRuntime.Exec(argv)
	} else {
		ociSpec, err := oci.NewSpec(argv)
		if err != nil {
			log.Printf("error: NewSpec\n")
			os.Exit(0)
		}

		rawSpec, err := ociSpec.Load()
		if err != nil {
			log.Printf("error: load spec\n")
			os.Exit(0)
		}

		image, err := image.NewCUDAImageFromSpec(rawSpec, cfg)
		if err != nil {
			log.Printf("new cuda image from spec\n")
			os.Exit(0)
		}

		gpuModifier := modifier.NewGraphicsModifier(image)
		sdkModifier := modifier.NewSdkModifier(image)
		mergeModifier := modifier.Merge(gpuModifier, sdkModifier)
		r := oci.NewModifyingRuntimeWrapper(lowLevelRuntime, ociSpec, mergeModifier)

		return r.Exec(argv)
	}
}

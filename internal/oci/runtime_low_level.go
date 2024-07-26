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

	"gitee.com/deep-spark/ix-container-runtime/internal/lookup"
	log "github.com/sirupsen/logrus"
)

// NewLowLevelRuntime creates a Runtime that wraps a low-level runtime executable.
// The executable specified is taken from the list of supplied candidates, with the first match
// present in the PATH being selected. A logger is also specified.
func NewLowLevelRuntime(candidates []string) (Runtime, error) {
	runtimePath, err := findRuntime(candidates)
	if err != nil {
		return nil, fmt.Errorf("error locating runtime: %v", err)
	}

	log.Infof("Using low-level runtime %v", runtimePath)
	return NewRuntimeForPath(runtimePath)
}

// findRuntime checks elements in a list of supplied candidates for a matching executable in the PATH.
// The absolute path to the first match is returned.
func findRuntime(candidates []string) (string, error) {
	if len(candidates) == 0 {
		return "", fmt.Errorf("at least one runtime candidate must be specified")
	}

	locator := lookup.NewExecutableLocator("/")
	for _, candidate := range candidates {
		log.Debugf("Looking for runtime binary '%v'", candidate)
		targets, err := locator.Locate(candidate)
		if err == nil && len(targets) > 0 {
			log.Debugf("Found runtime binary '%v'", targets)
			return targets[0], nil
		}
		log.Debugf("Runtime binary '%v' not found: %v (targets=%v)", candidate, err, targets)
	}

	return "", fmt.Errorf("no runtime binary found from candidate list: %v", candidates)
}

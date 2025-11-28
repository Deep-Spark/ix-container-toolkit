/**
# Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
# All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
**/

package modifier

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitee.com/deep-spark/ix-container-runtime/internal/config"
	"gitee.com/deep-spark/ix-container-runtime/internal/config/image"
	"gitee.com/deep-spark/ix-container-runtime/internal/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	log "github.com/sirupsen/logrus"
)

type sdkModifier struct {
	Change image.VisibleSdk
}

const (
	visibleSdkEnvvar   = "COREX_IMAGE"
	pathEnv            = "PATH"
	ldPathEnv          = "LD_LIBRARY_PATH"
	defaultDestination = "/usr/local/corex"

	pathAdded   = "/usr/local/corex/bin"
	ldPathAdded = "/usr/local/corex/lib64"
)

func (s sdkModifier) Modify(spec *specs.Spec) error {
	log.Debug("start sdk modification")

	if spec.Process == nil {
		spec.Process = &specs.Process{}
	}
	err := checkSourceDir(defaultDestination)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(defaultDestination)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", defaultDestination, err)
	}

	// Iterate through each entry in the directory.
	for _, entry := range entries {
		absPath := filepath.Join(defaultDestination, entry.Name())

		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: absPath,
			Source:      absPath,
			Type:        "bind",
			Options:     []string{"rbind", "ro", "nosuid", "nodev", "rprivate"},
		})

		log.Debugf("added mount: %s -> %s", absPath, absPath)
	}

	basePath := s.Change.Path()
	baseLdPath := s.Change.LdPath()
	newPathVal := joinPathLike(pathAdded, basePath)
	newLdPathVal := joinPathLike(ldPathAdded, baseLdPath)

	spec.Process.Env = setOrReplaceEnv(spec.Process.Env, pathEnv, newPathVal)
	spec.Process.Env = setOrReplaceEnv(spec.Process.Env, ldPathEnv, newLdPathVal)

	log.Debugf("---> pathval updated to :%v=%v", pathEnv, newPathVal)
	log.Debugf("---> ldpathval updated to :%v=%v", ldPathEnv, newLdPathVal)

	return nil
}

func NewSdkModifier(ig image.CUDA, cfg *config.Config) oci.SpecModifier {
	if !cfg.MountCorexSdk {
		log.Debug("MountCorexSdk is false, skipping Corex SDK modification")
		return nil
	}

	return sdkModifier{
		Change: ig.SdkFromEnvvars(visibleSdkEnvvar, pathEnv, ldPathEnv),
	}
}

func checkSourceDir(source string) error {
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", source)
	}
	return nil
}

func joinPathLike(add, base string) string {
	if add == "" {
		return base
	}
	if base == "" {
		return add
	}
	return add + ":" + base
}

func setOrReplaceEnv(env []string, key, value string) []string {
	prefix := key + "="
	out := env[:0]
	// remove all entries with the same key
	for _, e := range env {
		if !strings.HasPrefix(e, prefix) {
			out = append(out, e)
		}
	}
	// replace the value
	out = append(out, prefix+value)
	return out
}

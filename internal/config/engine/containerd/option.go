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

package containerd

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"

	"gitee.com/deep-spark/ix-container-runtime/internal/config/engine"
	log "github.com/sirupsen/logrus"
)

type builder struct {
	path                 string
	runtimeType          string
	useLegacyConfig      bool
	containerAnnotations []string
	configVersion        int
}

// Option defines a function that can be used to configure the config builder
type Option func(*builder)

// WithPath sets the path for the config builder
func WithPath(path string) Option {
	return func(b *builder) {
		b.path = path
	}
}

// WithRuntimeType sets the runtime type for the config builder
func WithRuntimeType(runtimeType string) Option {
	return func(b *builder) {
		b.runtimeType = runtimeType
	}
}

// WithUseLegacyConfig sets the useLegacyConfig flag for the config builder
func WithUseLegacyConfig(useLegacyConfig bool) Option {
	return func(b *builder) {
		b.useLegacyConfig = useLegacyConfig
	}
}

// WithContainerAnnotations sets the container annotations for the config builder
func WithContainerAnnotations(containerAnnotations ...string) Option {
	return func(b *builder) {
		b.containerAnnotations = containerAnnotations
	}
}

// WithConfigVersion sets the config version for the config builder
func WithConfigVersion(configVersion int) Option {
	return func(b *builder) {
		b.configVersion = configVersion
	}
}

func (b *builder) build() (engine.Interface, error) {
	if b.path == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	config, err := b.loadConfig(b.path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}
	config.RuntimeType = b.runtimeType
	config.UseDefaultRuntimeName = !b.useLegacyConfig
	config.ContainerAnnotations = b.containerAnnotations

	configVersion, err := b.parseVersion(config.Tree)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config version: %v", err)
	}

	log.Printf("Using config version: %v", configVersion)

	criRuntimePluginName, err := b.criRuntimePluginName(configVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get CRI runtime plugin name: %w", err)
	}

	log.Printf("Using CRI runtime plugin name: %v", criRuntimePluginName)

	config.CRIRuntimePluginName = criRuntimePluginName
	config.Version = configVersion
	switch configVersion {
	case 1:
		return (*ConfigV1)(config), nil
	default:
		return config, nil
	}

}

func (b *builder) criRuntimePluginName(configVersion int64) (string, error) {
	switch configVersion {
	case 1:
		return "cri", nil
	case 2:
		return "io.containerd.grpc.v1.cri", nil
	default:
		return "io.containerd.cri.v1.runtime", nil
	}
}

// loadConfig loads the containerd config from disk
func (b *builder) loadConfig(config string) (*Config, error) {
	info, err := os.Stat(config)
	if os.IsExist(err) && info.IsDir() {
		return nil, fmt.Errorf("config file is a directory")
	}

	if os.IsNotExist(err) {
		log.Infof("Config file does not exist; using empty config")
		config = "/dev/null"
	} else {
		log.Infof("Loading config from %v", config)
	}

	tomlConfig, err := toml.LoadFile(config)
	if err != nil {
		return nil, err
	}

	cfg := Config{
		Tree: tomlConfig,
	}
	return &cfg, nil
}

func (b *builder) parseVersion(c *toml.Tree) (int64, error) {
	if c == nil || len(c.Keys()) == 0 {
		// No config exists, or the config file is empty.
		if b.useLegacyConfig {
			// If a legacy config is explicitly requested, we default to a v1 config.
			return 1, nil
		}
		// Use the requested version.
		return int64(b.configVersion), nil
	}

	switch v := c.Get("version").(type) {
	case nil:
		return 1, nil
	case int64:
		return v, nil
	default:
		return -1, fmt.Errorf("unsupported type for version field: %v", v)
	}
}

func (b *Config) criRuntimePluginName(configVersion int64) (string, error) {
	switch configVersion {
	case 1:
		return "cri", nil
	case 2:
		return "io.containerd.grpc.v1.cri", nil
	default:
		return "io.containerd.cri.v1.runtime", nil
	}
}

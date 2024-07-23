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

package configure

import (
	"fmt"
	"log"

	"gitee.com/deep-spark/ix-container-runtime/internal/config/engine"
	"gitee.com/deep-spark/ix-container-runtime/internal/config/engine/containerd"
	"gitee.com/deep-spark/ix-container-runtime/internal/config/engine/docker"
	"github.com/urfave/cli/v2"
)

const (
	defaultRuntime                  = "docker"
	defaultContainerdConfigFilePath = "/etc/containerd/config.toml"
	defaultCrioConfigFilePath       = "/etc/crio/crio.conf"
	defaultDockerConfigFilePath     = "/etc/docker/daemon.json"
	defaultConfigFilePath           = "/etc/iluvatarcorex/ix-container-runtime/config.yaml"
	defaultConfigFileContent        = `librarypath: /usr/local/corex/lib64/libixml.so
sdksocketpath: /run/ix-sdk-manager/ix-sdk.sock`
)

type command struct {
}

// NewCommand constructs a configure command with the specified logger
func NewCommand() *cli.Command {
	c := command{}
	return c.build()
}

func (m command) build() *cli.Command {
	// Create a config struct to hold the parsed environment variables or command line flags
	config := newConfig()

	// Create the 'configure' command
	configure := cli.Command{
		Name:  "configure",
		Usage: "Add a runtime to the specified container engine",
		Action: func(c *cli.Context) error {
			return ConfigureRuntime(config)
		},
	}

	configure.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "runtime",
			Aliases:     []string{"r"},
			Usage:       "Specify the runtime to configure",
			Value:       defaultRuntime,
			Destination: &config.runtime,
		},
		&cli.BoolFlag{
			Name:        "ix-set-as-default",
			Aliases:     []string{"set-as-default"},
			Usage:       "set the Iluvatar runtime as the default runtime",
			Destination: &config.setAsDefault,
		},
		&cli.BoolFlag{
			Name:        "dry-run",
			Usage:       "update the runtime configuration as required but don't write changes to disk",
			Destination: &config.dryRun,
		},
		&cli.StringFlag{
			Name:        "ix-runtime-path",
			Aliases:     []string{"runtime-path"},
			Usage:       "specify the path to the ix runtime executable",
			Value:       engine.DefaultRuntimePath,
			Destination: &config.path,
		},
	}
	return &configure
}

type config struct {
	runtime        string
	configFilePath string
	setAsDefault   bool
	dryRun         bool
	path           string
}

func newConfig() config {
	return config{}
}
func (c *config) resolveConfigFilePath() string {
	if c.configFilePath != "" {
		return c.configFilePath
	}
	switch c.runtime {
	case "containerd":
		return defaultContainerdConfigFilePath
	case "crio":
		return defaultCrioConfigFilePath
	case "docker":
		return defaultDockerConfigFilePath
	}
	return ""
}

func (c *config) getOuputConfigPath() string {
	if c.dryRun {
		return ""
	}
	return c.resolveConfigFilePath()
}

func ConfigureRuntime(c config) error {
	configFilePath := c.resolveConfigFilePath()
	var cfg engine.Interface
	var err error
	switch c.runtime {
	case "containerd":
		cfg, err = containerd.New(
			containerd.WithPath(configFilePath),
		)
	case "docker":
		cfg, err = docker.New(
			docker.WithPath(configFilePath),
		)
	default:
		err = fmt.Errorf("unrecognized runtime '%v'", c.runtime)
	}
	if err != nil || cfg == nil {
		return fmt.Errorf("unable to load config for runtime %v: %v", c.runtime, err)
	}

	err = cfg.AddRuntime(
		engine.IluvatarRutimeName,
		c.path,
		c.setAsDefault,
	)
	if err != nil {
		return fmt.Errorf("unable to update config: %v", err)
	}

	outputPath := c.getOuputConfigPath()
	n, err := cfg.Save(outputPath)
	if err != nil {
		return fmt.Errorf("unable to flush config: %v", err)
	}

	if outputPath != "" {
		if n == 0 {
			log.Printf("Removed empty config from %v\n", outputPath)
		} else {
			log.Printf("Wrote updated config to %v\n", outputPath)
		}
		log.Printf("It is recommended that %v daemon be restarted.\n", c.runtime)
	}

	_, err = engine.Config(defaultConfigFilePath).Write([]byte(defaultConfigFileContent))
	if err != nil {
		return fmt.Errorf("unable to flush config: %v", defaultConfigFilePath)
	}

	return nil
}

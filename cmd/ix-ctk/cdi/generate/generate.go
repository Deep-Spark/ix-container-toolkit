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

package generate

import (
	"fmt"
	"log"
	"os"

	"gitee.com/deep-spark/ix-container-runtime/internal/config"
	"gitee.com/deep-spark/ix-container-runtime/pkg/ixcdi"
	"gitee.com/deep-spark/ix-container-runtime/pkg/ixcdi/spec"
	"gitee.com/deep-spark/ix-container-runtime/pkg/ixcdi/transform"
	"github.com/urfave/cli/v2"
)

const (
	allDeviceName = "all"
)

type command struct{}

type options struct {
	output               string
	deviceNameStrategies cli.StringSlice
	vendor               string
	class                string
}

// NewCommand constructs a generate-cdi command with the specified logger
func NewCommand() *cli.Command {
	c := command{}
	return c.build()
}

// build creates the CLI command
func (m command) build() *cli.Command {
	opts := options{}

	// Create the 'generate-cdi' command
	c := cli.Command{
		Name:  "generate",
		Usage: "Generate CDI specifications for use with CDI-enabled runtimes",
		Action: func(c *cli.Context) error {
			return m.run(c, &opts)
		},
	}

	c.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "output",
			Usage:       "Specify the file to output the generated CDI specification to. If this is '' the specification is output to STDOUT",
			Destination: &opts.output,
		},
		&cli.StringSliceFlag{
			Name:        "device-name-strategy",
			Usage:       "Specify the strategy for generating device names. If this is specified multiple times, the devices will be duplicated for each strategy. One of [index | uuid | type-index]",
			Value:       cli.NewStringSlice(ixcdi.DeviceNameStrategyIndex, ixcdi.DeviceNameStrategyUUID),
			Destination: &opts.deviceNameStrategies,
		},
		&cli.StringFlag{
			Name:        "vendor",
			Aliases:     []string{"cdi-vendor"},
			Usage:       "the vendor string to use for the generated CDI specification.",
			Value:       "iluvatar.com",
			Destination: &opts.vendor,
		},
		&cli.StringFlag{
			Name:        "class",
			Aliases:     []string{"cdi-class"},
			Usage:       "the class string to use for the generated CDI specification.",
			Value:       "gpu",
			Destination: &opts.class,
		},
	}

	return &c
}

func (m command) run(c *cli.Context, opts *options) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	spec, err := m.generateSpec(opts, cfg)
	if err != nil {
		return fmt.Errorf("failed to generate CDI spec: %v", err)
	}

	log.Printf("Generated CDI spec with version %v", spec.Raw().Version)

	if opts.output == "" {
		_, err := spec.WriteTo(os.Stdout)
		if err != nil {
			return fmt.Errorf("failed to write CDI spec to STDOUT: %v", err)
		}
		return nil
	}

	return spec.Save(opts.output)
}

func (m command) generateSpec(opts *options, cfg *config.Config) (spec.Interface, error) {
	var deviceNamers []ixcdi.DeviceNamer

	for _, strategy := range opts.deviceNameStrategies.Value() {
		deviceNamer, err := ixcdi.NewDeviceNamer(strategy)
		if err != nil {
			return nil, fmt.Errorf("failed to create device namer: %v", err)
		}
		deviceNamers = append(deviceNamers, deviceNamer)
	}

	cdilib, err := ixcdi.New(
		ixcdi.WithDeviceNamers(deviceNamers...),
		ixcdi.WithLibraryPath(cfg.LibraryPath),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create CDI library: %v", err)
	}

	deviceSpecs, err := cdilib.GetAllDeviceSpecs()
	if err != nil {
		return nil, fmt.Errorf("failed to create device CDI specs: %v", err)
	}

	return spec.New(
		spec.WithVendor(opts.vendor),
		spec.WithClass(opts.class),
		spec.WithDeviceSpecs(deviceSpecs),
		spec.WithPermissions(0644),
		spec.WithMergedDeviceOptions(
			transform.WithName(allDeviceName),
			transform.WithSkipIfExists(true),
		),
	)
}

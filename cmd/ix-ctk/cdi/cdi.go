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

package cdi

import (
	"gitee.com/deep-spark/ix-container-runtime/cmd/ix-ctk/cdi/generate"
	"github.com/urfave/cli/v2"
)

type command struct{}

// NewCommand constructs an info command with the specified logger
func NewCommand() *cli.Command {
	c := command{}
	return c.build()
}

// build
func (m command) build() *cli.Command {
	// Create the 'hook' command
	hook := cli.Command{
		Name:  "cdi",
		Usage: "Provide tools for interacting with Container Device Interface specifications",
	}

	hook.Subcommands = []*cli.Command{
		generate.NewCommand(),
	}

	return &hook
}

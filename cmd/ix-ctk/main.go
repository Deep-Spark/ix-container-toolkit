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

package main

import (
	"log"
	"os"

	"gitee.com/deep-spark/ix-container-runtime/cmd/ix-ctk/cdi"
	"gitee.com/deep-spark/ix-container-runtime/cmd/ix-ctk/runtime"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "ix-ctk",
		Usage: "A CLI tool to manage container runtimes and configurations",
		Commands: []*cli.Command{
			runtime.NewCommand(),
			cdi.NewCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

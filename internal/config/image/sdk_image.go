/*
*
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
*
*/
package image

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func (i CUDA) SdkFromEnvvars(sdkEnv, pathEnv, ldPathEnv string) VisibleSdk {
	log.Printf("enter SdkFromEnvvars !!!\n")
	setSdk := i.Cfg.DefaultSdk
	pathVal := ""
	ldPathVal := ""

	if sdk, ok := i.env[sdkEnv]; ok {
		setSdk = sdk
	}

	if path, ok := i.env[pathEnv]; ok {
		pathVal = strings.TrimSpace(path)
	}

	if ldPath, ok := i.env[ldPathEnv]; ok {
		ldPathVal = strings.Trim(ldPath, " \t\n\r")
	}

	log.Printf("setSdk name:%v  path:%v ldpath:%v\n", setSdk, pathVal, ldPathVal)
	return NewVisibleSdk(setSdk, pathVal, ldPathVal)
}

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

package ixcdi

import (
	"fmt"

	"gitee.com/deep-spark/go-ixml/pkg/ixml"
)

type toRequiredInfo struct {
	ixml.Device
}

type requiredInfo interface {
	getDevNodePath() (string, error)
}

func (d *toRequiredInfo) getDevNodePath() (string, error) {
	minor, ret := d.Device.GetMinorNumber()
	if ret != ixml.SUCCESS {
		return "", fmt.Errorf("error getting GPU device minor number: %d", ret)
	}
	path := fmt.Sprintf("/dev/iluvatar%d", minor)
	return path, nil
}
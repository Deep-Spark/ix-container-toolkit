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
	"tags.cncf.io/container-device-interface/specs-go"
)

// Interface defines the API for the ixcdi package
type Interface interface {
	// GetSpec() (spec.Interface, error)
	// GetCommonEdits() (*cdi.ContainerEdits, error)
	GetAllDeviceSpecs() ([]specs.Device, error)
	// GetGPUDeviceEdits(ixml.Device) (*cdi.ContainerEdits, error)
	// GetGPUDeviceSpecs(int, ixml.Device) ([]specs.Device, error)
	// GetDeviceSpecsByID(...string) ([]specs.Device, error)
}

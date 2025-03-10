/**
# Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
# Copyright (c) 2022, NVIDIA CORPORATION.  All rights reserved.
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

package discover

// charDevices is a discover for a list of character devices
type charDevices mounts

var _ Discover = (*charDevices)(nil)

// NewCharDeviceDiscoverer creates a discoverer which locates the specified set of device nodes.
func NewCharDeviceDiscoverer(devices []string) Discover {

	return (*charDevices)(newMounts(nil, "", devices))
}

// Mounts returns the discovered mounts for the charDevices.
// Since this explicitly specifies a device list, the mounts are nil.
func (d *charDevices) Mounts() ([]Mount, error) {
	return nil, nil
}

// Devices returns the discovered devices for the charDevices.
// Here the device nodes are first discovered as mounts and these are converted to devices.
func (d *charDevices) Devices() ([]Device, error) {
	var devices []Device
	for _, mount := range d.required {
		device := Device{
			HostPath: mount,
			Path:     mount,
		}
		devices = append(devices, device)
	}

	return devices, nil
}

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
	"errors"
	"fmt"
)

// DeviceNamer is an interface for getting device names
type DeviceNamer interface {
	GetDeviceName(int, string) (string, error)
}

// DeviceNamers represents a list of device namers
type DeviceNamers []DeviceNamer

// NewDeviceNamer creates a Device Namer based on the supplied strategy.
// This namer can be used to construct the names for GPU devices when generating the CDI spec.
func NewDeviceNamer(strategy string) (DeviceNamer, error) {
	switch strategy {
	case DeviceNameStrategyIndex:
		return deviceNameIndex{}, nil
	case DeviceNameStrategyTypeIndex:
		return deviceNameIndex{gpuPrefix: "gpu"}, nil
	case DeviceNameStrategyUUID:
		return deviceNameUUID{}, nil
	}

	return nil, fmt.Errorf("invalid device name strategy: %v", strategy)
}

type deviceNameIndex struct {
	gpuPrefix string
}

type deviceNameUUID struct{}

// Supported device naming strategies
const (
	// DeviceNameStrategyIndex generates devices names such as 0 or 1:0
	DeviceNameStrategyIndex = "index"
	// DeviceNameStrategyTypeIndex generates devices names such as gpu0 or mig1:0
	DeviceNameStrategyTypeIndex = "type-index"
	// DeviceNameStrategyUUID uses the device UUID as the name
	DeviceNameStrategyUUID = "uuid"
)

// GetDeviceName returns the name for the specified device based on the naming strategy
func (s deviceNameIndex) GetDeviceName(i int, _ string) (string, error) {
	return fmt.Sprintf("%s%d", s.gpuPrefix, i), nil
}

// GetDeviceName returns the name for the specified device based on the naming strategy
func (s deviceNameUUID) GetDeviceName(i int, uuid string) (string, error) {
	return uuid, nil
}

func (l DeviceNamers) GetDeviceNames(i int, d string) ([]string, error) {
	var names []string
	for _, namer := range l {
		name, err := namer.GetDeviceName(i, d)
		if err != nil {
			return nil, err
		}
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	if len(names) == 0 {
		return nil, errors.New("no names defined")
	}
	return names, nil
}

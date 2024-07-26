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
	"log"

	"gitee.com/deep-spark/go-ixml/pkg/ixml"
	"gitee.com/deep-spark/ix-container-runtime/pkg/ixcdi/discover"
	"gitee.com/deep-spark/ix-container-runtime/pkg/ixcdi/edits"
	"tags.cncf.io/container-device-interface/pkg/cdi"
	"tags.cncf.io/container-device-interface/specs-go"
)

type ixmllib ixcdilib

var _ Interface = (*ixmllib)(nil)

// GetAllDeviceSpecs returns the device specs for all available devices.
func (l *ixmllib) GetAllDeviceSpecs() ([]specs.Device, error) {
	var deviceSpecs []specs.Device
	var ret ixml.Return

	if l.libraryPath != "" {
		ret = ixml.AbsInit(l.libraryPath)
	} else {
		ret = ixml.Init()
	}

	if ret != ixml.SUCCESS {
		return nil, fmt.Errorf("failed to initialize ixml: %v", ret)
	}

	defer func() {
		if ret = ixml.Shutdown(); ret != ixml.SUCCESS {
			log.Printf("failed to shutdown ixml: %v", ret)
		}
	}()

	gpuDeviceSpecs, err := l.getGPUDeviceSpecs()
	if err != nil {
		return nil, err
	}
	deviceSpecs = append(deviceSpecs, gpuDeviceSpecs...)

	return deviceSpecs, nil
}

func (l *ixmllib) getGPUDeviceSpecs() ([]specs.Device, error) {
	var deviceSpecs []specs.Device
	var err error
	count, ret := ixml.DeviceGetCount()
	if ret != ixml.SUCCESS {
		log.Printf("failed to get count:%v\n", ret)
		return nil, fmt.Errorf("failed to get count:%v", ret)
	}

	log.Printf("Find GPU device count: %d\n", count)

	for i := uint(0); i < count; i++ {
		var device ixml.Device
		ret = ixml.DeviceGetHandleByIndex(i, &device)

		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get device at index %d: %v", i, ret)
		}
		specsForDevice, err := l.GetGPUDeviceSpecs(int(i), device)
		if err != nil {
			return nil, fmt.Errorf("failed to generate CDI edits for GPU devices: %v", err)
		}
		deviceSpecs = append(deviceSpecs, specsForDevice...)

	}
	return deviceSpecs, err
}

// GetGPUDeviceSpecs returns the CDI device specs for the full GPU represented by 'device'.
func (l *ixmllib) GetGPUDeviceSpecs(i int, d ixml.Device) ([]specs.Device, error) {
	edits, err := l.GetGPUDeviceEdits(d)
	if err != nil {
		return nil, fmt.Errorf("failed to get edits for device: %v", err)
	}

	var deviceSpecs []specs.Device
	uuid, ret := d.GetUUID()
	if ret != ixml.SUCCESS {
		return nil, fmt.Errorf("failed to get uuid for device: %v", ret)
	}
	names, err := l.deviceNamers.GetDeviceNames(i, uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get device name: %v", err)
	}
	for _, name := range names {
		spec := specs.Device{
			Name:           name,
			ContainerEdits: *edits.ContainerEdits,
		}
		deviceSpecs = append(deviceSpecs, spec)
	}

	return deviceSpecs, nil
}

// GetGPUDeviceEdits returns the CDI edits for the full GPU represented by 'device'.
func (l *ixmllib) GetGPUDeviceEdits(d ixml.Device) (*cdi.ContainerEdits, error) {
	device, err := l.newFullGPUDiscoverer(d)
	if err != nil {
		return nil, fmt.Errorf("failed to create device discoverer: %v", err)
	}

	editsForDevice, err := edits.FromDiscoverer(device)
	if err != nil {
		return nil, fmt.Errorf("failed to create container edits for device: %v", err)
	}

	return editsForDevice, nil
}

// newFullGPUDiscoverer creates a discoverer for the full GPU defined by the specified device.
func (l *ixmllib) newFullGPUDiscoverer(d ixml.Device) (discover.Discover, error) {
	ixmlDiscoverer, err := l.newIxmlDGPUDiscoverer(&toRequiredInfo{d})
	if err != nil {
		return nil, fmt.Errorf("failed to get devicenode: %v", err)
	}

	return ixmlDiscoverer, nil
}

func (l *ixmllib) newIxmlDGPUDiscoverer(d requiredInfo) (discover.Discover, error) {
	path, err := d.getDevNodePath()
	if err != nil {
		return nil, fmt.Errorf("error getting device node path: %w", err)
	}

	deviceNodePaths := []string{path}

	deviceNodes := discover.NewCharDeviceDiscoverer(
		deviceNodePaths,
	)
	return deviceNodes, nil
}

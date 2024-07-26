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

// Option is a function that configures the ixcdilib
type Option func(*ixcdilib)

// WithDeviceNamers sets the device namer for the library
func WithDeviceNamers(namers ...DeviceNamer) Option {
	return func(l *ixcdilib) {
		l.deviceNamers = namers
	}
}

// WithVendor sets the vendor for the library
func WithVendor(vendor string) Option {
	return func(o *ixcdilib) {
		o.vendor = vendor
	}
}

// WithClass sets the class for the library
func WithClass(class string) Option {
	return func(o *ixcdilib) {
		o.class = class
	}
}

// WithClass sets the class for the library
func WithLibraryPath(path string) Option {
	return func(o *ixcdilib) {
		o.libraryPath = path
	}
}

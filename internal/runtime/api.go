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

package runtime

type rt struct {
	modeOverride string
}

// Interface is the interface for the runtime library.
type Interface interface {
	Run([]string) error
}

// Option is a function that configures the runtime.
type Option func(*rt)

// New creates a runtime with the specified options.
func New(opts ...Option) Interface {
	r := rt{}
	for _, opt := range opts {
		opt(&r)
	}
	return &r
}

// WithModeOverride allows for overriding the mode specified in the config.
func WithModeOverride(mode string) Option {
	return func(r *rt) {
		r.modeOverride = mode
	}
}

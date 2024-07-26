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

package lookup

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// file can be used to locate file (or file-like elements) at a specified set of
// prefixes. The validity of a file is determined by a filter function.
type file struct {
	builder
	prefixes []string
}

// builder defines the builder for a file locator.
type builder struct {
	root        string
	searchPaths []string
	filter      func(string) error
	count       int
	isOptional  bool
}

type Option func(*builder)

func WithRoot(root string) Option {
	return func(f *builder) {
		f.root = root
	}
}

func WithSearchPaths(paths ...string) Option {
	return func(f *builder) {
		f.searchPaths = paths
	}
}

func WithCount(count int) Option {
	return func(f *builder) {
		f.count = count
	}
}

func WithFilter(assert func(string) error) Option {
	return func(f *builder) {
		f.filter = assert
	}
}

func getSearchPrefixes(root string, prefixes ...string) []string {
	seen := make(map[string]bool)
	var uniquePrefixes []string
	for _, p := range prefixes {
		if seen[p] {
			continue
		}
		seen[p] = true
		uniquePrefixes = append(uniquePrefixes, filepath.Join(root, p))
	}

	if len(uniquePrefixes) == 0 {
		uniquePrefixes = append(uniquePrefixes, root)
	}

	return uniquePrefixes
}

func (o builder) build() *file {
	f := file{
		builder: o,
		// Since the `Locate` implementations rely on the root already being specified we update
		// the prefixes to include the root.
		prefixes: getSearchPrefixes(o.root, o.searchPaths...),
	}
	return &f
}

func (p file) Locate(pattern string) ([]string, error) {
	var filenames []string

	log.Debugf("Locating %q in %v\n", pattern, p.prefixes)
visit:
	for _, prefix := range p.prefixes {
		pathPattern := filepath.Join(prefix, pattern)
		candidates, err := filepath.Glob(pathPattern)
		if err != nil {
			log.Debugf("Checking pattern '%v' failed: %v", pathPattern, err)
		}

		for _, candidate := range candidates {
			log.Debugf("Checking candidate '%v'\n", candidate)
			err := p.filter(candidate)
			if err != nil {
				log.Debugf("Candidate '%v' does not meet requirements: %v\n", candidate, err)
				continue
			}
			filenames = append(filenames, candidate)
			if p.count > 0 && len(filenames) == p.count {
				log.Debugf("Found %d candidates; ignoring further candidates\n", len(filenames))
				break visit
			}
		}
	}

	if !p.isOptional && len(filenames) == 0 {
		return nil, fmt.Errorf("pattern %v %w", pattern, ErrNotFound)
	}
	return filenames, nil
}

func newFileLocator(opts ...Option) *file {
	return newBuilder(opts...).build()
}

func newBuilder(opts ...Option) *builder {
	o := &builder{}
	for _, opt := range opts {
		opt(o)
	}

	if o.filter == nil {
		o.filter = assertFile
	}
	return o
}

// assertFile checks whether the specified path is a regular file
func assertFile(filename string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("error getting info for %v: %v", filename, err)
	}

	if info.IsDir() {
		return fmt.Errorf("specified path '%v' is a directory", filename)
	}

	return nil
}

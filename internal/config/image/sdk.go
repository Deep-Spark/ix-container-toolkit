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

type VisibleSdk interface {
	Name() string
	Path() string
	LdPath() string
}

type sdk struct {
	name   string
	path   string
	ldpath string
}

func (s sdk) Name() string {
	return s.name
}

func (s sdk) Path() string {
	return s.path
}

func (s sdk) LdPath() string {
	return s.ldpath
}

func NewVisibleSdk(sdkName, path, ldpath string) VisibleSdk {
	return sdk{
		name:   sdkName,
		path:   path,
		ldpath: ldpath,
	}
}

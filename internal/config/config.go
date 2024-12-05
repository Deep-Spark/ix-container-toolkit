/**
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
**/

package config

import (
	"fmt"
	"io"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"
)

const (
	cfgpath = "/etc/iluvatarcorex/ix-container-runtime/config.yaml"

	LogPath = "/var/log/iluvatarcorex/ix-container-toolkit/ix-container-runtime.log"

	LevelInfo    = "info"
	LevelDebug   = "debug"
	LevelTrace   = "trace"
	LevelWarning = "warning"
	LevelError   = "error"
	LevelFatal   = "fatal"
	LevelPanic   = "Panic"
)

type Config struct {
	Loglevel      string `json:"loglevel"             yaml:"loglevel,omitempty"`
	LogPath       string `json:"logpath"             yaml:"logpath,omitempty"`
	LibraryPath   string `json:"librarypath"             yaml:"librarypath,omitempty"`
	DefaultSdk    string `json:"defaultsdk" yaml:"defaultsdk"`
	SdkSocketPath string `json:"sdksocketpath" yaml:"sdksocketpath"`
}

func parseConfigFrom(reader io.Reader) (*Config, error) {
	var err error
	var configYaml []byte

	configYaml, err = io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read error: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(configYaml, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	return &cfg, nil
}

func (c *Config) update() error {
	var level log.Level
	if c.LogPath == "" {
		c.LogPath = LogPath
	}

	if c.Loglevel == "" {
		c.Loglevel = LevelInfo
	}

	switch c.Loglevel {
	case LevelInfo:
		level = log.InfoLevel
	case LevelDebug:
		level = log.DebugLevel
	case LevelTrace:
		level = log.TraceLevel
	case LevelWarning:
		level = log.WarnLevel
	case LevelError:
		level = log.ErrorLevel
	case LevelFatal:
		level = log.FatalLevel
	case LevelPanic:
		level = log.PanicLevel
	default:
		level = log.InfoLevel
	}
	log.SetLevel(level)

	reader, err := os.OpenFile(c.LogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("error opening config file: %v", err)
	}
	log.SetOutput(reader)

	// 启用日志行号
	log.SetReportCaller(true)

	// 自定义日志格式，包括文件名和行号
	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// 返回 "function" 和 "file:line" 的格式
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", f.File, f.Line)
		},
		FullTimestamp: true,
	})
	return nil
}

func LoadConfig() (*Config, error) {
	var cfg *Config
	reader, err := os.Open(cfgpath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = &Config{}
		} else {
			return nil, fmt.Errorf("error opening config file: %v", err)
		}
	} else {
		defer reader.Close()
		cfg, err = parseConfigFrom(reader)
		if err != nil {
			return nil, err
		}
	}
	err = cfg.update()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

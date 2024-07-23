#!/bin/bash

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


docker_tmp_name="docker_tmp.json"
containerd_tmp_name="containerd_tmp.toml"
docker_backup_name="docker_backup.json"
containerd_backup_name="containerd_backup.toml"
iluvatar_tmp_name="iluvatar_config_tmp.yaml"
iluvatar_backup_name="iluvatar_config_backup.yaml"

iluvatar_config_dir="/etc/iluvatarcorex/ix-container-runtime"
iluvatar_config_name="config.yaml"
iluvatar_config_path="${iluvatar_config_dir}/${iluvatar_config_name}"
out_name="ix-container-runtime"

iluvatar_log_dir="/var/log/iluvatarcorex/ix-container-toolkit/"

check_runtime_install() {
	isExist="$(which $1)"
	if [ -z "${isExist}" ];then
	  return 0
	else
	  return 1
	fi
}

configure_docker_env() {
	work_dir="$1"
	override="$2"
	isInstalled="$3"

	temp_config_path="${work_dir}/${docker_tmp_name}"
	backup_config_path="${work_dir}/${docker_backup_name}"
	docker_config_path="/etc/docker/daemon.json"
	docker_config_content=""
	check_runtime_install docker
	docker_installed="$?"

	if [ "${docker_installed}" -eq 1 ];then
		mkdir -p $(dirname ${docker_config_path})

		if [ -f "${docker_config_path}" ];then
			docker_config_content=$(cat ${docker_config_path})
			cp ${docker_config_path} ${backup_config_path}
		fi
		if [ "${isInstalled}" = true ]; then
			printf "%s\n" "${docker_config_content}" | ./yq_linux_amd64 '.runtimes += {"iluvatar": {"args": [], "path":"/usr/local/bin/ix-container-runtime"}}|."default-runtime"="iluvatar"' -o json > ${temp_config_path}
		else
			printf "%s\n" "${docker_config_content}" | ./yq_linux_amd64 'del(."runtimes"."iluvatar", ."default-runtime")' -o json > ${temp_config_path}
		fi

		if [ "${override}"  = true ];then
			cp ${temp_config_path} ${docker_config_path}
		fi
	fi
}

configure_containerd_env() {
	work_dir="$1"
	override="$2"
	isInstalled="$3"

	temp_config_path="${work_dir}/${containerd_tmp_name}"
	backup_config_path="${work_dir}/${containerd_backup_name}"
	containerd_config_path="/etc/containerd/config.toml"
	containerd_config_content=$(containerd config default)
	check_runtime_install containerd
	containerd_installed="$?"

	if [ "${containerd_installed}" -eq 1 ];then
		mkdir -p $(dirname ${containerd_config_path})

		if [ -f "${containerd_config_path}" ];then
			containerd_config_content=$(cat ${containerd_config_path})
			cp ${containerd_config_path} ${backup_config_path}
		fi

		if [ "${isInstalled}" = true ]; then
			printf "%s\n" "${containerd_config_content}"| ./yq_linux_amd64 --input-format=toml '.plugins."io.containerd.grpc.v1.cri".containerd.runtimes.iluvatar += {"runtime_type": "io.containerd.runc.v1", "options": {"BinaryName":"/usr/local/bin/ix-container-runtime"}}'|./dasel_linux_amd64 -r yaml -w toml > ${temp_config_path}
		else
			printf "%s\n" "${containerd_config_content}"| ./yq_linux_amd64 --input-format=toml 'del(.plugins."io.containerd.grpc.v1.cri".containerd.runtimes.iluvatar)' |./dasel_linux_amd64 -r yaml -w toml > ${temp_config_path}
		fi

		if [ "${override}"  = true ];then
			cp ${temp_config_path} ${containerd_config_path}
		fi
	fi
}

configure_iluvatar_env() {
	work_dir="$1"
	override="$2"
	isInstalled="$3"
	temp_config_path="${work_dir}/${iluvatar_tmp_name}"
	backup_config_path="${work_dir}/${iluvatar_backup_name}"
	iluvatar_config_content="librarypath: /usr/local/corex/lib64/libixml.so"

	mkdir -p ${iluvatar_config_dir}


	if [ -f "${iluvatar_config_path}" ];then
		iluvatar_config_content=$(cat ${iluvatar_config_path})
		cp ${iluvatar_config_path} ${backup_config_path}
	fi

	if [ "${isInstalled}" = true ]; then
		printf "%s\n" "${iluvatar_config_content}"| ./yq_linux_amd64 --input-format=yaml '.librarypath = "/usr/local/corex/lib64/libixml.so"|.sdksocketpath = "/run/ix-sdk-manager/ix-sdk.sock"' -o yaml > ${temp_config_path}

		if [ "${override}"  = true ];then
			cp ${temp_config_path} ${iluvatar_config_path}
		fi
	else
		if [ "${override}"  = true ];then
			rm -rf ${iluvatar_config_dir}
		fi
	fi
}

install() {
	local override=false
	local rm=false
	while [[ "$1" != "" ]];do
		case "$1" in
		--override)
			shift
			override=true
			;;
		--rm)
			shift
			override=true
			rm=true
			;;
		*)
			help
			exit 1
			;;
		esac
	done

	echo "Start to Install ix-Runtime"
	temp_dir=$(mktemp -d)
	if [ ! -d "${temp_dir}" ];then
		echo "Failed to create temp directory"
		exit 1
	else
		echo "Working temp directory is ${temp_dir}"
	fi

	mkdir -p ${iluvatar_log_dir}
	cp ix-container-runtime /usr/local/bin  
	cp ix-ctk /usr/local/bin
	chmod 755 /usr/local/bin/ix-ctk
	chmod 755 /usr/local/bin/ix-container-runtime
	configure_docker_env ${temp_dir} ${override} true
	configure_containerd_env ${temp_dir} ${override} true
	configure_iluvatar_env ${temp_dir} ${override} true

	if [ "${override}" = true ];then
		echo "----------------"
		echo "Install Finished"
		echo "Please restart your service to update your environment"
		if [ "${rm}" = true ]; then
			rm -rf ${temp_dir}
		else
			echo "Please delete the temporary directory ${temp_dir} once you have confirmed they are correct."
		fi
	else
		echo "********************************"
		echo "Generated Configuration Finished"
		echo "Please enter ${temp_dir} to configure the environment"
		echo "The file ${docker_tmp_name} used for configrue the docker config file"
		echo "The file ${containerd_tmp_name} used for configrue the containerd config file"
		echo "The file ${iluvatar_tmp_name} used for configrue the iluvatar config file, install "
		echo "mannual by put this file under ${iluvatar_config_dir} and renamed to ${iluvatar_config_name}"
	fi
}

uninstall() {
	local override=false
	local rm=false
	while [[ "$1" != "" ]];do
		case "$1" in
		--override)
			shift
			override=true
			;;
		--rm)
			shift
			override=true
			rm=true
			;;
		*)
			help
			exit 1
			;;
		esac
	done
	echo "Start to Uninstall ix-Runtime"
	temp_dir=$(mktemp -d)
	if [ ! -d "${temp_dir}" ];then
		echo "Failed to create temp directory"
		exit 1
	else
		echo "Working temp directory is ${temp_dir}"
	fi


	configure_docker_env ${temp_dir} ${override} false
	configure_containerd_env ${temp_dir} ${override} false
	configure_iluvatar_env ${temp_dir} ${override} false
	if [ "${override}" = true ];then
		rm /usr/local/bin/ix-container-runtime
		echo "----------------"
                echo "Uninstall Finished"
                echo "Please restart your service to update your environment"
                if [ "${rm}" = true ]; then
                        rm -rf ${temp_dir}
                else
                        echo "Please delete the temporary directory ${temp_dir} once you have confirmed they are correct."
                fi
	else
                echo "********************************"
                echo "Generated Configuration Finished"
                echo "Please enter ${temp_dir} to configure the environment"
                echo "The file ${docker_tmp_name} used for configrue the docker config file"
                echo "The file ${containerd_tmp_name} used for configrue the containerd config file"
	fi
}

help() {
	echo "Usage: ix-runtime-installer.run {install|uninstall} --override --rm"
	echo "--override    will override the configuration file(default is false)"
	echo "--rm          not only override the configuration file, but also remove the workspace(default is false)"
	exit 1
}

if [ -z "$1" ]; then
	help
fi

case "$1" in
	install)
		shift
		install "$@"
		;;
	uninstall)
		shift
		uninstall "$@"
		;;
	*)
		help
	;;
esac

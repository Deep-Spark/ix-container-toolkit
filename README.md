# IluvatarCorex Container Toolkit

- [Introducton](#introduction)  
- [Building the IX Container Toolkit](#building-the-ix-container-toolkit)  
- [Configuring](#configuring)
    - [Configuring Docker](#configuring-docker)
    - [Configuring Containerd](#configuring-containerd)
    - [Configuring Crio](#configuring-crio)
- [Running Samples](#running-samples)
    - [Running a Sample Workload with Docker](#running-a-sample-workload-with-docker)
    - [Running a Sample Workload with Containerd/Crio(for kubernetes 1.22+)](#running-a-sample-workload-with-containerd/crio(for-kubernetes-1.22+))
    - [Running a Sample Workload with Podman](#running-a-sample-workload-with-podman)
- [License](#license)

## Introduction

The Iluvatar Container Toolkit allows users to build and run GPU accelerated containers. The toolkit includes a container runtime(ix-container-runtime) and utilities to automatically configure containers to leverage Iluvatar GPUs.

## Building the IX Container Toolkit

```shell
make all
```
This will build the ix-container-runtime and ix-ctk binary, see logging for more details.

```shell
sudo make install
```
This will install the ix-container-runtime and ix-ctk binary, see logging for more details.

## Configuring

### Configuring Docker

1. Run the ix-ctk

```shell
sudo ix-ctk runtime configure --runtime docker --ix-set-as-default
```

Run the `runtime configure --runtime docker` command will automatically add the configuration of `ix-container-runtime` to `/etc/docker/daemon.json`, The default configuration is generated in the config.yaml `/etc/iluvatarcorex/ix-container-runtime/` file, and  the docker service needs to be restarted for the configuration to take effect.

2. Restart docker service

```shell
sudo systemctl daemon-reload
sudo systemctl restart docker
```

### Configuring Containerd (for Kubernetes)

1. Run the ix-ctk

```shell
sudo ix-ctk runtime configure --runtime containerd --ix-set-as-default
```

Run the `runtime configure --runtime containerd` command will automatically add the configuration of `ix-container-runtime` to `/etc/containerd/config.toml`, The default configuration is generated in the config.yaml `/etc/iluvatarcorex/ix-container-runtime/` file, and  the containerd service needs to be restarted for the configuration to take effect.

2. Restart containerd service

```shell
sudo systemctl daemon-reload
sudo systemctl restart containerd
```

### Configuring Crio (for Kubernetes)

1. Run the ix-ctk

```shell
sudo ix-ctk runtime configure --runtime crio --ix-set-as-default
```

Run the `runtime configure --runtime crio` command will automatically add the configuration of `ix-container-runtime` to `/etc/crio/crio.conf`, The default configuration is generated in the config.yaml `/etc/iluvatarcorex/ix-container-runtime/` file, and  the containerd service needs to be restarted for the configuration to take effect.

2. Restart crio service

```shell
sudo systemctl daemon-reload
sudo systemctl restart crio
```

## Running Samples

### Running a Sample Workload with Docker

After you install and configure the toolkit and install Iluvatar GPU Driver and SDK, you can verify your installation by running a sample workload.

```shell
sudo docker run -it --rm --runtime iluvatar -e IX_VISIBLE_DEVICES=0 corex:4.0.0 ixsmi
```

Your output should resemble the following output:

```shell
+-----------------------------------------------------------------------------+
|  IX-ML: 4.0.0       Driver Version: 4.1.0       CUDA Version: N/A           |
|-------------------------------+----------------------+----------------------|
| GPU  Name                     | Bus-Id               | Clock-SM  Clock-Mem  |
| Fan  Temp  Perf  Pwr:Usage/Cap|      Memory-Usage    | GPU-Util  Compute M. |
|===============================+======================+======================|
| 0    Iluvatar BI-V150S        | 00000000:8A:00.0     | 500MHz    1600MHz    |
| 0%   33C   P0    N/A / N/A    | 114MiB / 32768MiB    | 0%        Default    |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                       GPU Memory |
|  GPU        PID      Process name                                Usage(MiB) |
|=============================================================================|
|  No running processes found                                                 |
+-----------------------------------------------------------------------------+
```

### Running a Sample Workload with Containerd/Crio(for kubernetes 1.22+)


After you install and configure the toolkit and install Iluvatar GPU Driver and SDK, you can verify your installation by running a sample workload with kubernetes.

You can create a yaml file `corex-example.yaml` with the following content:
```shell
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  name: iluvatar
handler: iluvatar
---
apiVersion: v1
kind: Pod
metadata:
  name: corex-example
spec:
  runtimeClassName: iluvatar
  containers:
  - name: corex-example
    image: docker.io/library/corex:4.0.0
    command: ["/usr/local/corex/bin/ixsmi"]
    args: ["-l"]
    env:
    - name: IX_VISIBLE_DEVICES
      value: "0"
```

apply the yaml file to kubernetes cluster and see the logs:
```shell
kubectl apply -f corex-example.yaml
kubectl logs corex-example
```

Your output should resemble the following output:

```shell
+-----------------------------------------------------------------------------+
|  IX-ML: 4.0.0       Driver Version: 4.1.0       CUDA Version: N/A           |
|-------------------------------+----------------------+----------------------|
| GPU  Name                     | Bus-Id               | Clock-SM  Clock-Mem  |
| Fan  Temp  Perf  Pwr:Usage/Cap|      Memory-Usage    | GPU-Util  Compute M. |
|===============================+======================+======================|
| 0    Iluvatar BI-V150S        | 00000000:8A:00.0     | 500MHz    1600MHz    |
| 0%   33C   P0    N/A / N/A    | 114MiB / 32768MiB    | 0%        Default    |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                       GPU Memory |
|  GPU        PID      Process name                                Usage(MiB) |
|=============================================================================|
|  No running processes found                                                 |
+-----------------------------------------------------------------------------+
```

### Running a Sample Workload with Podman

After you install and configura the toolkit (including generating a CDI specification) and install Iluvatar GPU Driver and SDK, you can verify your installation by running a sample workload.

1. Generate a CDI specification

```shell
$ sudo ix-ctk cdi generate  --output=/etc/cdi/ix.yaml
```

- The CDI file path is usually `/etc/cdi` or `/var/run/cdi`.
- Use `sudo` to ensure that the `/etc/cdi/ix.yaml` file can be created.
- If the `--output` parameter is not used, the output will be output to `stdout` by default.

`Example output:`

```shell
2024/11/07 03:11:56 Find GPU device count: 2
2024/11/07 03:11:56 Generated CDI spec with version 0.5.0
```

2. Test CDI with Podman(v4.3.0+)

Podman can specify to use the GPU device defined in the CDI file with the `--device` parameter, the format is `iluvatar.com/gpu=device name`.  
Currently, the CDI file generated by ix-ctk contains three types of GPU device names: `index`, `UUID`, `all`.

```shell
$ podman run --rm --device iluvatar.com/gpu=0 localhost/corex:4.0.0 ixsmi
```

Your output should resemble the following output:

```shell
Timestamp    Thu Nov  7 03:23:57 2024
+-----------------------------------------------------------------------------+
|  IX-ML: 4.0.0       Driver Version: 4.2.0       CUDA Version: N/A           |
|-------------------------------+----------------------+----------------------|
| GPU  Name                     | Bus-Id               | Clock-SM  Clock-Mem  |
| Fan  Temp  Perf  Pwr:Usage/Cap|      Memory-Usage    | GPU-Util  Compute M. |
|===============================+======================+======================|
| 0    Iluvatar BI-V100         | 00000000:00:06.0     | 1500MHz   1200MHz    |
| 0%   35C   P0    56W / 250W   | 9MiB / 32768MiB      | 0%        Default    |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                       GPU Memory |
|  GPU        PID      Process name                                Usage(MiB) |
|=============================================================================|
|  No running processes found                                                 |
+-----------------------------------------------------------------------------+
```

## License

Copyright (c) 2024 Iluvatar CoreX. All rights reserved. This project has an Apache-2.0 license, as
found in the [LICENSE](LICENSE) file.

# SGT

## 背景

这个进程用于管理机器上面的其他agent，比如监控的agent、安全的agent，管理主要是：安装、升级、卸载、查看启动状态，不做其他事情。省去客户手工安装其他agent的工作。

## 安装

虚机创建的时候会自动安装此进程，如需对存量虚机安装，可以执行：

```
curl -s http://mirrors.intra.didiyun.com/didiyun_resource/sgd-v1.sh | sh
```

只能在滴滴云的虚机里运行这条指令，适用64位linux系统

## 资源占用

安装完成之后机器上会有sgd和sga两个进程，sgd内存占用小于10MB，承担管理其他agent的核心业务逻辑，sga内存占用小于4MB，是sgd进程的伴生进程，在sgd挂掉的时候负责将其拉起。cpu使用率小于1%

## 规范要求

sgd管理的其他agent需要提供control脚本，打到tar.gz包里，control脚本需要具备可执行权限，支持这些参数：pid | version | start | stop | uninstall | install，sgd就是利用业务agent的control脚本来做管理的
# 功能

支持健康检查的服务发现

# 说明

分布式服务提供者集成worker，用于向etcd注册自身节点，并定时发送心跳。

master作为服务发现者，向etcd查询服务。

可结合github.com/q191201771/load_blance选择合适的负载方式使用。

# TODO

* mixin consistent-hash or round-robin

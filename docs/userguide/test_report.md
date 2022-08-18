# SLIs/SLOs

可扩展性和性能是多集群联邦的重要特性，作为多集群联邦的用户，我们期望在以上两方面有服务质量的保证。在进行大规模性能测试之前，我们需要定义测量指标。在参考了Kubernetes社区的SLI(Service Level Indicator)/SLO(Service Level Objectives)和多集群的典型应用，Karmada社区定义了以下SLI/SLO来衡量多集群联邦的服务质量。

1. API Call Latency
   
对于联邦控制面内的资源，包括用户提交的需要下发到成员集群的资源模板、调度策略等：

| Status                      | SLI                                   | SLO |
| ------------------------- | ----------------------------------------------- | ----- |
| Offical    | 最近5min的单个Object Mutating API P99 时延                   | 除聚合API和CRD外，P99 <= 1s |
| Offical    | 最近5min的non-streaming read-only P99 API时延                 | 除聚合API和CRD外，Scope=resource, P99 <= 1s，Scope=namespace, P99 <= 5s，Scope=cluster, P99 <= 30s|

对于成员集群内的资源：

| Status                      | SLI                                   | SLO |
| ------------------------- | ----------------------------------------------- | ----- |
| Offical    | 最近5min的单个Object Mutating API P99 时延                   | 除聚合API和CRD外，P99 <= 1s |
| Offical    | 最近5min的non-streaming read-only P99 API时延                 | 除聚合API和CRD外，Scope=resource, P99 <= 1s，Scope=namespace, P99 <= 5s，Scope=cluster, P99 <= 30s|

2. Cluster Startup Latency

| Status                      | SLI                                   | SLO |
| ------------------------- | ----------------------------------------------- | ----- |
| Offical    | 集群从接入联邦控制面到状态能被控制面正确收集的时间，不考虑控制面与成员集群之间的网络波动                   | X |

3. Pod Startup Latency

| Status                      | SLI                                   | SLO |
| ------------------------- | ----------------------------------------------- | ----- |
| Offical    | 用户在联邦控制面提交资源模板和下发策略后到无状态Pod在成员集群的启动延时，不考虑控制面与成员集群之间的网络波动但考虑单集群的Pod启动延时                   | X |

4. Resource usage

| Status                      | SLI                                   | SLO |
| ------------------------- | ----------------------------------------------- | ----- |
| WIP    | 在接入一定数量的集群后集群联邦维持其正常工作所必需的资源使用量                   | X |


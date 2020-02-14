# slabinfo-exporter

### 背景
最近线上 k8s 集群遇到了一个业务接口延时规律性的变高，经过查看[阿里云给的文档](https://yq.aliyun.com/articles/697773)
得出结论是 `节点 dentry 链表过程，导致查询 slab 状态信息时，服务器致禁用中断时间过长`，上面文档中也是留下了两问题

* 是什么程序在反复地获取 slab 信息，产生类似 cat /proc/slabinfo 的效果?
* 这么多 dentry 生成的原因是什么?

### 怎么办？
上面两个问题我也不知道答案是什么，但是可以把 slabinfo 监控起来。找了一圈没发现现存监控轮子，只能自己造轮子了。。。

### 原理

```bash
# 把下面这条命令的输出结果转换成 metrics 数据
$ slabtop -o | grep -Ev "^$|OBJS ACTIVE|Minimum / Average / Maximum Object|Active / Total" |awk '{print $1","$2","$3","$4","$5","$6","$7","$8}'
```
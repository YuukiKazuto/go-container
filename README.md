# go-container
## 简介
本包参考Java Util包下的集合类型，是基于golang1.18+版本的泛型，
用迭代器模式实现的线性表数据结构。

本包封装了列表、栈、队列等线性表容器，各个容器都有顺序表和双向链表
两种对应的实现，且都是协程安全的。

本包目前还处于开发阶段，未通过测试，一些注释也尚未添加，暂时不建议使用。

如果想要提前试用，可以使用以下代码：
```shell
go get -u github.com/YuukiKazuto/go-container
```
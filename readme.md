# DHT (beta)

See project at [github](https://github.com/HsunGong/ppca-dht-test/wiki/) [acm wiki](https://acm.sjtu.edu.cn/wiki/PPCA_2019)

## Overview

- A DHT can be viewed as a dictionary service distributed over a network: it provides access to a common shared key->value data-store, distributed over participating nodes with great performance and scalability.
- From a user perspective, a DHT essentially provides a map interface, with two main operations: `put(key, value)` and `get(key)`. Get will retrieve values stored at a certain key while put (often called announce) will store a value on the network. Note that many values can be stored under the same key.
- There are many algorithms to implement dht, where **Kademlia** DHT algorithm requires to contact only O(log(N)) nodes for a get operation, N being the number of nodes in the network. This property makes DHTs very scalable as demonstrated.

More info in [Wiki:DHT](https://en.wikipedia.org/wiki/Distributed_hash_table)


## Target
- Use go-lang to implement a DHT with basic functions
- Use this DHT to implement an easy application

# Syllabus

## Schedule
- Learn about Golang and at least one DHT protocol
  > **Project 1(Not constrained)**: get close to golang: implement a client-server protocol using golang. With naive tests.
  > forward message and mutli-threads(2 with client and server)
- Implement a DHT protocol in Go
  > **Project 2**: implement DHT: using any one model is allowed. With naive and strong tests. E.g. chatroom, torrent, dht-filesystem.
- Build a higher level application on top of your DHT
  > **Project 3\***: implement an application. No test.
- Bonus: 
  > **Project 4\***implement DHT with another algorithm (as mentioned above in overview), or optimize your application.
## Reference

- Learn Go
> [A Tour of Go](https://tour.golang.org/)\
> [Go package docs](http://golang.org/pkg/)\
> [Books about Go](https://github.com/golang/go/wiki/Books)

- DHT models
> [Chord](https://en.wikipedia.org/wiki/Chord_(peer-to-peer))\
> [Pastry](https://en.wikipedia.org/wiki/Pastry_(DHT))\
> [Kademlia](https://en.wikipedia.org/wiki/Kademlia)

- Related project framework
> [Dixie](https://cit.dixie.edu/cs/3410/asst_chord.php)\
> [CMU](https://www.cs.cmu.edu/~dga/15-744/S07/lectures/16-dht.pdf)\
> [MIT](https://pdos.csail.mit.edu/papers/sit-phd-thesis.pdf)


# Test(beta)
  
## Range

- nodes in DHT network >= 50
- `<key, value>` in DHT network >= 1500
- Give some time for your network to resume stable.

## Contents

### Requirements

(1) Implement this interface in `interface.go`
```go 
type dhtNode interface {
    Get(k string) (string, bool)
	Put(k string, v string) bool
	Del(k string) bool
	Run(wg *sync.WaitGroup)
	Create()
	Join(addr string) bool
	Quit()
	Ping(addr string) bool
}
```

<!-- type dhtNode interface {
    Get(k string) (string, bool)
    Put(k string, v string) bool
    Del(k string) bool
    Run(wg *sync.WaitGroup)
    Create()
    Join(addr string) bool
    Quit()
    Ping(addr string) bool
} -->

(2) Overwrite in `userdef.go`
`NewNode()`

### Tests

- standard test
  > get, put, del, join, quit, ping
- standard test of additive functions
  > append info in `<k, v>`\
  > del_append info in `<k,v>`

- additive test
  > get, put, del, join, quit, ping\
  > force quit\
  > load(max nodes and max data)\
  > unstable put and get\
  > mixed put, get, join, quit

<!-- - 初始时由一个 DHT 节点组成网络
- 循环五次：
> A. 向网络中加入 15 个节点
> B. 选择 5 个节点退出网络
> 每次一个节点加入/一个节点退出后等待 3s
- 上述的 A, B 操作之后均进行：
> 等待 30s
向 DHT 网络中插入 300 个键值对
在网络中现有的键中随机选择 200 个进行查询操作并验证正确性
从网络中现有的键中随机选择 150 个删除 -->

## Score
- 10% - learn go of navie protocol
- 60% - DHT tests
- 30% - Applications with stable network


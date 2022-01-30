# 并发编程

## 当单例模式遇上多线程
作为最被滥用的设计模式之一，很多程序员都自己实现过至少一次单例模式（singleton pattern）。然而，在当今的多线程环境，要正确写对这个模式，其中的细节颇有一番讲究。   
以下是单例模式最基本的实现方式：
<figure>
  <pre><code> type singleton struct {}
    var instance *singleton
    func getSingleton() *singleton {
        if instance == nil {
            instance = new(singleton)
        }
        return instance
    }
  </code></pre>
  <figcaption>Fig.1 classic singleton</figcaption>
</figure>

```go
    type singleton struct {}
    var instance *singleton
    func getSingleton() *singleton {
        if instance == nil {
            instance = new(singleton)
        }
        return instance
    }
```
这个实现在单线程环境下基本没有问题，但如果在多线程的环境中使用便不再可靠。假设有两个线程A和B都需要获取singleton实例，考虑以下的执行顺序：
1. 线程A进入了`getSingleton`函数，且执行了instance是否为空的判断，结果为true；
2. 线程A被操作系统暂停运行了；
3. 线程B进入了`getSingleton`函数，并在完整执行了整个函数后返回；
4. 线程A被唤醒，从第11行开始执行。
可见，这样的执行顺序将重复创建singleton实例，导致该模式失效。我们可以添加一些代码，来观察这个实现在多线程环境中具体是如何失效的：
```go
    // 使用一个有可观测副作用的函数来构造singleton
    func makeSingleton() *singleton {
        fmt.Println("makeSingleton")
        return new(singleton)
    }
    func getSingleton() *singleton {
        if instance == nil {
            instance = makeSingleton()
        }
        return instance
    }

    // 开启多个线程并发获取singleton实例，观察输出
    func main() {
        n := 10
        wg := &sync.WaitGroup{}
        for i := 0; i < n; i++ {
            wg.Add(1)
            go func(i int) {
                defer wg.Done()
                fmt.Println(i, getSingleton())
            }(i)
        }
        wg.Wait()
    }
```
> 可以去这里亲自操作：https://go.dev/play/p/5AOMKGIgkZU  
如果运行这段代码，几乎每次都会得到不同的运行结果（偶尔能遇见正确结果）。以下是一种可能的运行结果：
```
makeSingleton
makeSingleton
3 &{}
makeSingleton
0 &{}
makeSingleton
6 &{}
makeSingleton
7 &{}
9 &{}
......
```
观察以上运行结果，可以发现：
1. `makeSingleton`被多次调用，构造了多个singleton对象。这说明有多个线程都读取到`instance`为`nil`；某个线程给instance赋值的结果，似乎并未没有同步给其他线程。
2. 线程的被执行的顺序毫无规律，即使代码中按标号递增的顺序去创建线程。

## 并发？
这个例子显示，思考、理解并发问题需要不同的mental model，先来看下并发的定义。很多资料用一句大白话来定义并发（concurrency）：并发指的是多个计算任务同时进行。这句话确实是对的，但对于初识并发的程序员来说，其实很不好理解，下面尝试qualify这句话。

这句定义中最显著的就是“同时”这个词，它并非指“同一时刻”，而是“同一时段”。这引入了两个概念，分别是时段和顺序。

时段（time period）意味着，计算任务的执行而非瞬间完成，而是需要一段时间的。以`instance == nil`为例，这是一条看似极其primitive的操作，但它可能要花多个CPU cycle才能完成：CPU先尝试从其私有缓存中查找`instance`的值，若未果，再先后尝试从公有缓存和内存中取（这些存储器都有不同的访问延迟，且访存优先级越低的存储器延迟越高），最后才是与`nil`去比较。
> memory hierarchy latency  
  
再来看如何定义顺序（order）。在有了时段的概念后，每个任务X的时段都可以由开始和结束这两个时刻（moment）来标识，记为`X = (tBeg, tEnd)`。那么两个任务的先后顺序可以这样定义：假设有任务A和B，如果`A.tEnd > B.tBeg`，则任务A先于B，记为`A < B`。

可见，并发意味着，以程序中的计算任务为元素形成一个偏序集，元素间的关系是刚才定义的先后顺序，不能比较的元素之间则为并发关系。取决于具体的环境，两个并发的计算任务之间可能是以下任意一种情况：
1. A < B；
2. B < A；
3. A和B执行时段有重叠。

那么，如果在`main`函数中创建多个线程，
```go
    func main() {
        n := 10
        wg := &sync.WaitGroup{}
        for i := 0; i < n; i++ {
            wg.Add(1)
            go func(i int) {
                defer wg.Done()
                fmt.Println(i, getSingleton())
            }(i)
        }
        wg.Wait()
    }   
```
设`main`所在的线程为`gMain`，`main`所创建的子线程设为`{gi | 0 <= i < 10}`，那么对于任意的`gi`都有`gMain < gi`；但由于gi起始时间不确定，gi之间无法比较，因此gi间是并发关系。

由于`gi`之间的顺序关系未定义，所以各`gi`间的这几个操作很可能是**相互交织在一起**的。假设把任意一个`gi`的操作简单地分为三步：
1. `getSingleton()`
2. `fmt.Println()`
3. `wg.Done()`

那么两个gi的其中一种交织执行情况可能是：  
|time\gi|g0     |g1     |
|----   |---    |---    |
|t0     |getSingleton() |   |
|t1     |               |getSingleton()|
|t2     |fmt.Println()  |   |
|t3     |               |fmt.Println()|
|t4     |               |wg.Done()|
|t5     |wg.Done()      |   |

回顾前文，之前以一个“可能的执行顺序”为例，来尝试证明经典单例实现在并发环境下的不可靠。在那个例子中，线程A和B也是并发关系，其操作是相互交织的。之所以可以这样证明，因为并发则顺序未定义，只要能找出一种让模式出错的执行顺序，就能证明其是错的。

## 如何
put a big lock on
sync


## 参考文献
[] Meyers, S., & Alexandrescu, A. (2004). C++ and the Perils of Double-Checked Locking.  
[] Wikipedia. Partially ordered set. https://en.wikipedia.org/wiki/Partially_ordered_set


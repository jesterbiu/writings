# Go语言碎碎念
由于工作需要，最近都在用 Go 语言进行开发，至今已半年有余。而我原本是半个 C++ 爱好者（虽然没有实战经验），这回却要学习Go。据说语言塑造人的思维，虽然编程语言非自然语言，但我也确实体验到了与 C++ 大相径庭的编程思维。

## Go的背景
Go语言初创团队三剑客分别是：Ken Thompson, Rob Pike, 和Robert Griesemer。据 Ken Thompsen 所言，当初发明 Go 是因为：
> When the three of us got started, it was pure research. The three of us got together and decided that we hated C++.

10年后，Go 语言增长的势头非常好——但 C++ 也是。Go 发布两年后，C++ 正式推出了 C++11 标准，并在此后规律地每三年出一个新版本，开发者数量再次回升。一定不能小看这部分工作，C++标准从来没有被某个大公司把持过（我指的 Java），C++ 也没有公司为其做商业推广（我指的还是 Java），C++ 标委会全部用爱发电，这部分工作不容忽视。

而 Rob Pike 作为知名嘴炮，多次在公开场合“冷嘲热讽” C++ 的复杂，强调 Go 自有另一套注重简约的语言哲学。那他们是什么样的关系，为什么听着是竞争对手，但实际应用领域看着不搭边？这段话被我掏出来咀嚼了很多次，看了 Rob Pike 的一些表述和其他的一些资料，我隐约察觉到这句话有其语境，那就是解决 Google 的问题。

Google 的海量规模，会遇到前所未有的问题，也会遇到已知问题被放大到难以忍受的程度。对于前所未有的问题，Google的那三篇最有名的分布式论文就是非常好的例子：Google File System, MapReduce 和 BigTable。有兴趣的同学可以自行去查阅。

举例几点已知的问题：
1. C++的编译耗时，这是当初Go想解决的一大痛点。C++为了兼容C以提高易用度（他俩都诞生在Bell Lab），包括其编译模型。Google内部的C++代码量级本身相当可观，加上Google内部代码仓库的组织方式，再加上用70年代的C编译模型来编译，属于凌迟。
2. 再比如语言本身的复杂度，C++必读书目中Scott Meyers的两本Effective C++加起来将近百条款项，大几百页，而你才刚刚入门。Go对此表示“大道至简”，语言一共只有不到30个关键字，以及programming language学界的paper一律不看。
3. 编写并发程序。C++直到C++11才有内存模型以及线程的概念，而Go团队在初期就认为语言要内置对并发的支持。

说了这么些，意思是Go的提出有其特定背景，曰“知其然，知其所以然”。这点非常重要，如果你是一个C++开发者，也请不要照搬 Google 的 C++ Style Guide，那是他们内部的需要！“我们讨厌C++”更多是一句吐槽，Go无意也无法替代某某语言，但它提供一个新颖独特的选项。所以，事后两个语言都活得好好的，C++开发者继续在性能敏感（~~历史遗留~~）领域圈地自萌，而Go开发者里来了很多之前写Python却觉得太慢的朋友。

[] Interview with Ken Thompson, http://web.archive.org/web/20130105013259/https://www.drdobbs.com/open-source/interview-with-ken-thompson/229502480  
[] "Lang NEXT 2014 Panel Systems Programming in 2014 and Beyond", https://www.youtube.com/watch?v=ZQR32nTVF_4  
[] Expressiveness of Go, https://talks.golang.org/2010/ExpressivenessOfGo-2010.pdf  
[] “为什么Go语言如此不受待见？”, https://zhihu.com/question/27867348/answer/114125733  

## 基于CSP的并发模式
在中文互联网上，介绍Go几乎总是离不开“Go协程”这个概念。我之前专门写有一篇文章来argue这个问题：我不认为这是个好的说法。是的，goroutine 实现中的一些机制很像协程，但协程实在是太原始而底层的一个概念了，goroutine 无论是设计还是使用都更像线程，但比操作系统的线程更轻量—— Rob Pike 本人也称其为线程[Concurrency Is Not Parallelism]，我想设计者的看法也足够权威了。

把 Go 的并发特性排在第一位来聊，我认为是相当程度上是众望所归的。据说并发编程有两大套路，一是 shared memory，二是 message passing。前者是最普遍的模式，它直接源自 Dijkstra 于60年代发布的一系列关于并发编程的文章，这些文章奠定了计算机科学中的并发概念；后者一般追溯到Hoare于1976年发布的论文Communicating Sequential Processes（以下简称 CSP）。

### Communicating Sequential Processes
CSP 是一篇很有野心的论文（原文如此），它尝试定义一系列编程原语（primitive）来解决并行进程的同步问题。这里要qualify所谓的“并行进程”：
1. 原文确实写的“parallel”，但按当代眼光来看，并行更多是硬件的能力而非程序的组织结构，这篇论文也确实是提出了处理并发问题的模型。我认为在原文的语境中，我们可以暂时不去细究并发和并行之别。
2. 原文中的“process”并非对应到某个具体OS的process上，比如Linux process，文中是一个更抽象的概念。我认为理解为执行流（thread of execution）就不错。

下面介绍CSP 的 Introduction 部分提出的要点。

首先，CSP 引入了一个基于Dijkstra的`parbegin`的原语，可以开启一组新的并行进程，并等待它们执行完毕。类似Go语言中的`go`关键字与`sync.WaitGroup`操作的合体。CSP中对进程有特殊的限制，它们不能通过更新全局变量来通讯或同步。

论文里称，当时赋值操作已经被研究得相当清楚，但影响外部环境的IO操作却并不明晰。CSP 提出将进程间通信的IO操作作为编程语言内置的原语。这个对应到Go的`channel`操作：`<-`和`->`。

进程间通信操作以两个进程名为参数，分别是输出和输入进程；且通信操作是无缓冲的，等待通信的进程会阻塞，这对应到Go默认的无缓冲`channel`。

此外：输入操作可以作为在guarded command的条件，输入条件为真即意味着该通信操作当前可以无阻塞地立即进行；循环体也可以以输入操作为条件，循环会运行至所有消息源进程都已终结。

显然，Go的并发模式很大程度上借鉴了这篇论文的成果，并使得这篇论文的成果在近40年后再次为人所熟知。事实上，Rob Pike等人在Bell Lab做过Plan 9操作系统，当时该操作系统就使用了CSP模式的并发设计；如今他们在Google又有机会掏出来用了。另一方面，erlang也是一个使用CSP做并发的编程语言。

### Message Passing
这一套的好处，鼓励更高层次的处理并发问题，少见raw sync primitive，并非去除。在实现和OS层还是一个shared memory。

## 多态
Go的多态实现也是非常有意思的一点。

[] goroutine, 协程, COE - Hungbiu的文章 - 知乎, https://zhuanlan.zhihu.com/p/404452442  
[] Concurrency Is Not Parallelism - Rob Pike, https://www.youtube.com/watch?v=qmg1CF3gZQ0  
[] Communicating Sequential Processes - C.A.R. Hoare, http://www.cs.ox.ac.uk/people/bill.roscoe/publications/4.pdf  
[] Guarded Command - E.W. Dijkstra, https://dl.acm.org/doi/pdf/10.5555/1074100.1074433#:~:text=The%20term%20guarded%20command%2C%20as,execution%20is%20controlled%20by%20B.  


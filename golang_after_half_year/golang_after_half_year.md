# Go语言碎碎念
由于工作需要，最近都在用Go语言进行开发，至今已半年有余。而我原本是半个C++爱好者（虽然没有实战经验），这回却要学习Go。据说语言塑造人的思维，虽然编程语言非自然语言，但我也确实体验到了与C++大相径庭的编程思维。

## Go的背景
Go语言初创团队三剑客分别是：Ken Thompson, Rob Pike, 和Robert Griesemer。据 Ken Thompsen 所言，当初发明 Go 是因为：
> When the three of us got started, it was pure research. The three of us got together and decided that we hated C++.

10年后，Go语言增长的势头非常好——但C++也是。Go发布两年后，C++正式推出了C++11标准，并在此后规律地每三年出一个新版本，开发者数量再次回升。一定不能小看这部分工作，C++标准从来没有被某个大公司把持过（我指的Java），C++也没有公司为其做商业推广（我指的还是Java），C++标委会全部用爱发电，这部分工作不容忽视。

而 Rob Pike 作为知名嘴炮，多次在公开场合冷嘲热讽C++的复杂，强调Go自有另一套注重简约的语言哲学。那他们是什么样的关系，为什么听着是竞争对手，但实际应用领域看着不搭边？这段话被我掏出来咀嚼了很多次，看了 Rob Pike 的一些表述和其他的一些资料，我隐约察觉到这句话有其语境，那就是解决Google的问题。

Google 的海量规模，会遇到前所未有的问题，也会遇到已知问题被放大到难以忍受的程度。对于前所未有的问题，Google的那三篇最有名的分布式论文就是非常好的例子：Google File System, MapReduce和BigTable。有兴趣的同学可以自行去查阅。

举例几点已知的问题：
1. C++的编译耗时，这是当初Go想解决的一大痛点。C++为了兼容C以提高易用度（他俩都诞生在Bell Lab），包括其编译模型。Google内部的C++代码量级本身相当可观，加上Google内部代码仓库的组织方式，再加上用70年代的C编译模型来编译，属于凌迟。
2. 再比如语言本身的复杂度，C++必读书目中Scott Meyers的两本Effective C++加起来将近百条款项，大几百页，而你才刚刚入门。Go对此表示“大道至简”，语言一共只有不到30个关键字，以及programming language学界的paper一律不看。
3. 编写并发程序。C++直到C++11才有内存模型（memory model）以及线程的概念，而Go团队在初期就认为语言要内置对并发的支持。

说了这么些，意思是Go的提出有其特定背景，曰“知其然，知其所以然”。这点非常重要，如果你是一个C++开发者，也请不要照搬 Google 的 C++ Style Guide，那是他们内部的需要！“我们讨厌C++”更多是一句吐槽，Go无意也无法替代某某语言，但它提供一个新颖独特的选项。所以，事后两个语言都活得好好的，C++开发者继续在性能敏感（~~历史遗留~~）领域圈地自萌，而Go开发者里来了很多之前写Python却觉得太慢的朋友。

[] Interview with Ken Thompson, http://web.archive.org/web/20130105013259/https://www.drdobbs.com/open-source/interview-with-ken-thompson/229502480  
[] "Lang NEXT 2014 Panel Systems Programming in 2014 and Beyond", https://www.youtube.com/watch?v=ZQR32nTVF_4  
[] Expressiveness of Go, https://talks.golang.org/2010/ExpressivenessOfGo-2010.pdf  
[] “为什么Go语言如此不受待见？”, https://zhihu.com/question/27867348/answer/114125733  

## 基于CSP的并发模式
在中文互联网上，介绍Go几乎总是离不开“Go协程”这个概念。我之前专门写有一篇文章来argue这个问题：我不认为这是个好的说法。是的，goroutine 实现中的一些机制很像协程，但协程实在是太原始而底层的一个概念了，goroutine 无论是设计还是使用都更像线程，但比操作系统的线程更轻量—— Rob Pike 本人也称其为线程[Concurrency Is Not Parallelism]，我想设计者的看法也足够权威了。

把 Go 的并发特性排在第一位来聊，我认为是相当程度上是众望所归的。据传，并发编程的模式可以对应到进程间通信（inter-process communication）的两大套路：一是共享内存（shared memory），二是消息传递（message passing）。

前者是最普遍的模式，它更多地源自Dijkstra于60年代发布的一系列关于并发编程的文章，89年的[SRC]对这些情况做了总结，这些文章奠定了计算机科学中的并发概念；后者则可能追溯到Hoare于1976年发布的论文Communicating Sequential Processes（以下简称 CSP）。当然，接下来我们会看到这两个套路更多是不同的问题建模方法，而在实现上非互斥的关系。

### 并发的世界
当下，一台X86-64架构的服务器可能拥有百倍于20年前的可用内存，但内存容量增长的速度依旧远远落后于人类软件需求的增长速度，CPU的单核心性能进步也遇到瓶颈[Sutter]。这是提高程序并发性能的根本需求：程序员不能再依赖新的硬件来解决性能问题。并发编程因而不再停留于学术论文或实验室中，而成为工业界的热门话题。

C++一直是Google内部研发的主力语言，但在C11/C++11标准前，这两个编程语言是没有“线程”概念的，也没有内存模型。要进行并发编程，只能依赖第三方库或编译器内置函数（compiler intrinsic）。而这个第三方库，通常就是POSIX提供的线程及相关同步设施。然而，很多时候，直接使用系统接口来编写并发代码是困难的。

直接使用OS线程的问题是，其粒度较大，而且属于OS实现，应用程序对其的控制力有限；既然OS线程的粒度相对固定，那就很难能直接映射到应用程序多样的、更细粒度的并发工作负载。为此，程序员们可谓想尽了办法。比如，最常见的模式是线程池，Java标准库对此有很好的支持；C++社区还有在努力引入协程（通常比`goroutine`更底层），比如boost fiber，腾讯的libco，最新的标准C++20提供语言级的协程支持。

而系统提供的同步原语（synchronization primitive）大部分直接来自经典计算机论文，比如mutex、semaphore来自Dijkstra，monitor（condition variable）来自Hoare。这些同步原语的设计背景多是操作系统实现，而对于编写应用程序显得过于底层，从而阻碍程序员表达并发的程序结构或模式。

什么是并发的程序结构或模式？考虑以下最有名或常用的几个模式。
1. Map-Reduce：  
![Map-Reduce](https://res.cloudinary.com/talend/image/upload/f_auto/q_auto/v1633234207/resources/seo-articles/seo-what-is-mapreduce_gj9ehi.jpg)  
2. Fork-Join：  
![Fork-Join](https://www.researchgate.net/profile/David-Beckingsale-2/publication/305806654/figure/fig2/AS:669527134175236@1536639116271/The-fork-join-model-used-for-thread-level-parallelism-in-OpenMP.png)
3. Pipeline:
![Pipeline](https://www.researchgate.net/publication/351818084/figure/fig2/AS:1026998491160582@1621866929847/Example-of-pipeline-parallel-training.png)

对于有并发直觉的程序员，使用这样的模式去对并发问题建模，是非常符合直觉的；但要用同步原语来实现，却是非常反直觉的。**实现这些模式本身，仍然是实现程序的控制流而非业务；再要用同步原语来实现，代码只会更加肮脏**
。出于以上现实的需要，Go语言团队认为编程语言本身需要提供表达并发模式的基础设施。他们的选择是Hoare的那篇CSP，玩消息传递流派。实际上，他们（或者说，Rob Pike？）对于CSP的偏好或者说选择，要追溯至80年代了，这方面的历史可以参考Russ Cox的博客[]。

并发可能是Go语言最漂亮优雅的特性。我认为其中很重要的一个原因，就是其并发设计来自于一篇自顶向下、面向高层次抽象的论文，而非面向实现的。此处谈到的“面向实现”，指的是从硬件属性出发，向高级语言及其编程模式构建。例如C语言就是一门面向实现的语言，它设计初衷就想做离汇编一步之遥的高级语言，从而获得极强的移植性。而从高层次抽象向实现拟合，可以举例的是SQL，它一举将数据定义和查询的语言，与关系数据库的执行计划、存储引擎的实现解耦。

### Communicating Sequential Processes
CSP 是一篇很有野心的论文（原文如此），它尝试通过定义进程间的通信原语（primitive）来解决进程间的同步问题。这里要qualify所谓的“并行进程”：
1. 原文确实写的“parallel”，但按当代眼光来看，并行更侧重强调利用硬件的执行能力，而并发强调程序的组织结构；这篇论文也更多是提出处理并发问题抽象模型。我认为在原文的语境中，“parallel”暂时先做并发之意理解。
2. 原文中的“process”并非映射到某个具体OS的进程实现上，如Linux进程或Windows进程，文中是一个更抽象的概念。我认为理解为执行流（thread of execution）就不错。

下面介绍CSP 的 Introduction 部分提出的要点。

首先，CSP 引入了一个原语，可以开启一组新的并行进程，并等待它们执行完毕。这个原语基于Dijkstra的`parbegin`。类似Go语言中的`go`关键字与`sync.WaitGroup`操作的合体。CSP中对进程有特殊的限制，它们不能通过更新全局变量来通讯或同步。

论文里称，当时赋值操作已经被研究得相当清楚，但影响外部环境的IO操作却并不明晰。CSP 提出将进程间通信的IO操作作为编程语言内置的原语。这个对应到Go的`channel`操作：`<-`和`->`。

进程间通信操作以两个进程名为参数，分别是输出和输入进程；且通信操作是无缓冲的，等待通信的进程会阻塞，这对应到Go默认的无缓冲`channel`。

此外：输入操作可以作为在guarded command的条件，输入条件为真即意味着该通信操作当前可以无阻塞地立即进行；循环体也可以以输入操作为条件，循环会运行至所有消息源进程都已终结。

显然，Go的并发模式很大程度上借鉴了这篇论文的成果，并使得这篇论文的成果在近40年后再次为人所熟知。事实上，Rob Pike等人在Bell Lab做过Plan 9操作系统，当时该操作系统就使用了CSP模式的并发设计；如今他们在Google又有机会掏出来用了。另一方面，erlang也是一个使用CSP做并发的编程语言。

### 从CSP到Go并发模式
indirection, indirection, indirection! 见GMP。


[] goroutine, 协程, COE - Hungbiu的文章 - 知乎, https://zhuanlan.zhihu.com/p/404452442  
[] An Introduction to Programming with Threads - Andrew D. Birrell, https://www.hpl.hp.com/techreports/Compaq-DEC/SRC-RR-35.pdf  
[] Concurrency Is Not Parallelism - Rob Pike, https://www.youtube.com/watch?v=qmg1CF3gZQ0  
[] Communicating Sequential Processes - C.A.R. Hoare, http://www.cs.ox.ac.uk/people/bill.roscoe/publications/4.pdf  
[] Guarded Command - E.W. Dijkstra, https://dl.acm.org/doi/pdf/10.5555/1074100.1074433#:~:text=The%20term%20guarded%20command%2C%20as,execution%20is%20controlled%20by%20B.  
[] Go Concurrency Patterns, https://go.dev/blog/pipelines 
[] Advanced Go Concurrency Patterns, https://go.dev/blog/io2013-talk-concurrency  
[] Solution of a Problem in Concurrent Programming Control - E.W. Dijkstra, http://rust-class.org/static/classes/class19/dijkstra.pdf  
[] EWD35, http://www.cs.utexas.edu/users/EWD/ewd00xx/EWD35.PDF  
[] Bell Labs and CSP Threads - Russ Cox, https://swtch.com/~rsc/thread/  
[] Structured Parallelism programming Patterns: Patterns for Efficient Computation
## 多态
Go的多态实现也是非常有意思的一点。
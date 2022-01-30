# Go语言碎碎念：用Go开发半年之后
由于工作需要，最近都在用Go语言进行开发，至今已半年有余。而我原本是半个C++爱好者（虽然没有实战经验），这回却要学习Go。据说语言塑造人的思维，虽然编程语言非自然语言，但我也确实体验到了与C++大相径庭的编程思维。

## Go的背景
Go语言初创团队三剑客分别是：Ken Thompson, Rob Pike, 和Robert Griesemer。据Ken Thompsen所言，当初发明Go是因为：
> When the three of us got started, it was pure research. The three of us got together and decided that we hated C++.

10年后，Go语言增长的势头非常好——但C++也是。Go发布两年后，C++正式推出了C++11标准，并在此后规律地每三年出一个新版本，开发者数量再次回升。一定不能小看这部分工作，C++标准从来没有被某个大公司把持过（我指的Java），C++也没有公司为其做商业推广（我指的还是Java），C++标委会全部用爱发电，这部分工作不容忽视。

而Rob Pike作为知名嘴炮，多次在公开场合“冷嘲热讽”C++的复杂，强调Go自有另一套注重简约的语言哲学。那他们是什么样的关系，为什么听着是竞争对手，但实际应用领域看着不搭边？这段话被我掏出来咀嚼了很多次，看了Rob Pike的一些表述和其他的一些资料，我隐约察觉到这句话有其语境，那就是解决Google的问题。

Google的海量规模，会遇到前所未有的问题，也会遇到已知问题被放大到难以忍受的程度。对于前所未有的问题，Google的那三篇最有名的分布式论文就是非常好的例子：Google File System, MapReduce和BigTable。有兴趣的同学可以自行去查阅。

举例几点已知的问题：
1. C++的编译耗时，这是当初Go想解决的一大痛点。C++为了兼容C以提高易用度（他俩都诞生在Bell Lab），包括其编译模型。Google内部的C++代码量级本身相当可观，加上Google内部代码仓库的组织方式，再加上用70年代的C编译模型来编译，属于凌迟。
2. 再比如语言本身的复杂度，C++必读书目中Scott Meyers的两本Effective C++加起来将近百条款项，大几百页，而你才刚刚入门。Go对此表示“大道至简”，语言一共只有不到30个关键字，以及programming language学界的paper一律不看。
3. 编写并发程序。C++直到C++11才有内存模型以及线程的概念，而Go团队在初期就认为语言要内置对并发的支持。

说了这么些，意思是Go的提出有其特定背景，曰“知其然，知其所以然”。这点非常重要，如果你是一个C++开发者，也请不要照搬Google的C++ Style Guide，那是他们内部的需要！“我们讨厌C++”更多是一句吐槽，Go无意也无法替代某某语言，但它提供一个新颖独特的选项。所以，事后两个语言都活得好好的，C++开发者继续在性能敏感（~~历史遗留~~）领域圈地自萌，而Go开发者里来了很多之前写Python却觉得太慢的朋友。

[] Interview with Ken Thompson, http://web.archive.org/web/20130105013259/https://www.drdobbs.com/open-source/interview-with-ken-thompson/229502480
[] "Lang NEXT 2014 Panel Systems Programming in 2014 and Beyond", https://www.youtube.com/watch?v=ZQR32nTVF_4
[] Expressiveness of Go, https://talks.golang.org/2010/ExpressivenessOfGo-2010.pdf
[] “为什么Go语言如此不受待见？”, https://zhihu.com/question/27867348/answer/114125733

## 基于CSP的并发模式
在中文互联网上，介绍Go几乎总是离不开“Go协程”这个概念。我之前和朋友或同事聊过，也专门写有一篇文章来argue这个问题：我不认为这是个好的说法。是的，goroutine实现中的一些机制很像协程，但协程实在是太原始而底层的一个概念了，goroutine无论是设计还是使用都更像线程，但比操作系统的线程更轻量——Rob Pike本人也称其为thread，我想设计者的看法也足够权威了。

把这个特性排在第一位来聊，我认为是相当程度上是众望所归的。据说并发编程有两大套路，一是shared memory，二是message passing。前者是最普遍的模式，它直接源自Dijkstra于60年代发布的一系列关于并发编程的文章，奠定了计算机科学中的并发概念；后者一般追溯到Hoare于1976年发布的论文Communicating Sequential Processes。

他们在plan 9实践过。

[] goroutine, 协程, COE - Hungbiu的文章 - 知乎, https://zhuanlan.zhihu.com/p/404452442
[] Concurrency Is Not Parallelism, https://www.youtube.com/watch?v=qmg1CF3gZQ0
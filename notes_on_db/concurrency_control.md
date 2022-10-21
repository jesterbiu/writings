# 关系数据库的并发控制

关系数据库对用户提供的工作单位是 transaction，concurrency control 机制要处理多个 transaction 同时读写同一组对象的问题，此处对象可能是逻辑对象，比如一条 record、一个 table、index，也可能是物理对象，比如硬盘上存放的 data page 或 log file。

## Serializability: Essential Property

例子：

> Assume at first A and B each have $1000. What are the possible outcomes of running T1 and T2 ? 

```
T1: BEGIN; A=A-100; B=B+100; COMMIT 
T2: BEGIN; A=A*1.06; B=B*1.06; COMMIT
```

> Many! But A+B should be: → $2000*1.06=$2120 There is no guarantee that T1 will execute before T2 or vice-versa, if both are submitted together. But the net effect must be equivalent to these two transactions running serially in some order

DB 实际会并发执行事务，意味着不同事务中的操作是交叉进行的（interleaving）。同样的一组事务可能的执行序列有很多，然而，为了保证正确性，我们只能允许其中的一部分存在。

如何定义执行序的正确性？如果一个交叉执行序的最终结果，等于某一个串行执行序的结果，那么就认为是正确的——这被称为 *serializability*，而”等于“的定义是：

> Equivalent with respect to what? *Conflict equivalence*, the relative order of execution of the conflicting operations belonging to committed transactions in two schedules are the same.

 又，之所以只需要满足”某一个“串行序列的结果，是因为每个串行序列中的操作是全序的，会明确定义每个事务的先后顺序；而并发是偏序关系，并发的事务间没有明确先后发生顺序，这一点给到DB 一定的灵活空间以调度执行顺序，最大化操作输出。

> concurrency control, uwaterloo: https://cs.uwaterloo.ca/~tozsu/courses/cs448/notes/9.ConcurrencyControl-ho.pdf
>
> concurrency control, cmu: https://15445.courses.cs.cmu.edu/fall2021/slides/15-concurrencycontrol.pdf

conflict operation（DB 语境）的定义如下很简单（在并发编程语境里，data race 的定义也类似）：

- They are by different transactions
- They operate on the same object and at least one of them is a write.

在 DB 并发控制的语境里，不正确地处理 conflict operation 导致的结果被称为 *anomaly*。



two-phase locking

hierarchical lock

multi-version concurrency control
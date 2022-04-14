# 论文导读：How to make ad-hoc polymorphism less ad-hoc 

## significance

这篇论文提出了 type class 的概念来更好地支持 ad-hoc polymorphism 的应用。Type class 可以视为一种函数重载的机制，解决了 OOP、bounded quantification 的一些问题。

## Polymorphism 简介

原文介绍了 parametric polymorphism 和 ad-hoc polymorphism 的定义。

## Motivation

介绍当前 ad-hoc（函数重载）的局限。以 equality 运算为例：

1. 每个类型都写一份重载，按照 concrete type 来解析函数（类似 C++/Java 的函数重载）。`eq(Int, Int)`, `eq(Bool, Bool)`...
2. `(==) :: a -> a -> Bool`，让相等运算有类似 parametric polymorphism，但纯粹的 parametric 会导致部分类型的相等语义逻辑不正确。比如比较两个 set 是否相等，是指两个集合互为彼此的子集，而不考虑 set 中元素是否同样的顺序存放。
3.  `(==) :: a(==) -> a(==) -> Bool`，只有实现了 equality 这个接口的类型才可以比较。

对于 OOP：借鉴虚表的实现方式。对于3.，我们这样实现：`eq(vtable, a, a)`，其中 vtable 指向变量a的类型对应的相等函数。

## Subclass

要求 `Num` 的类型必须能具备相等性：

```haskell
class Eq a => Num a where 
	(+)    :: a -> a -> a
	(*)    :: a -> a -> a
	negate :: a -> a
```

`a` 属于 `Eq`，是 `a` 属于 `Num` 的必要条件。

## Go Generic

dict 实现

## C++ Concept

```c++
template <typename T>
concept eq = requires(T a, T b) {
    { a == b } -> std::convertible_to<bool>;
};

template <typename T>
concept ord = requires(T a, T b) {
    { a < b } -> std::convertible_to<bool>;
};

template <typename T>
concept cmp = eq<T> && ord<T>;

struct person {
    int identity;
    std::string name;
};

bool operator==(const person& lhs, const person& rhs) {
    return lhs.identity == rhs.identity;
}

bool operator<(const person& lhs, const person& rhs) {
    return lhs.identity < rhs.identity;
}
```




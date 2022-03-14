# Go 语言碎碎念：多态

说到多态，你会想到什么？是 Java 的类继承与虚函数表，还是 C++ 的模板与重载决议，抑或是 Go 的 interface？

> In programming language theory and type theory, polymorphism is the provision of a single interface to entities of different types or the use of a single symbol to represent multiple different types. [1]

编程语言中多态性（polymorphism）的不严谨的定义是，用同一个 interface 来使用多个不同的类型。简单地说，一个函数可以运作与不同类型的参数上。此处，“函数”不能仅仅是`func`，而是一个更抽象的概念；而“类型”指的是变量的类型，即“type”（而非“kind”）的概念。这样的定义显得含糊，本文接下来以 Go 语言的多态为主线，详细地聊一下多态。

## 无处不在的多态
如今 C++/Java 流派的面向对象编程有最广泛的使用，因此可能有不少程序员对于多态或其他计算机科学词汇的理解都建立在该流派的语境之下。说到多态，脑子里可能出现的是基于类继承实现的动态派发。然而，可以想得再简单而基础一些：各种运算符通常也是多态的！比如，相等运算符`==`。

虽然 Go 不支持自定义运算符重载，但它的`==`实现却是多态的，可以类比于语言内置的重载。  
- 对于标量原始类型，如`int`或`bool`，`==`的语义是比较标量的值是否相等，其实现一般就是一条CPU的cmp指令。  
- 对于`string`，`==`比较两个字符串是否为相同的 rune 序列；又由于 Go 使用 UTF-8 编码，`==`的实现很可能就是一条`memcmp`（可以去这个网站查看对应的汇编代码：https://godbolt.org）。  
- 对于`chan`, 可以比较两个`chan`变量是否引用同一个 channel。  
- 对于`map`或 slice，`==`仅可以把`map`或 slice 变量和`nil`值比较，检查是否进行过初始化。  

Go 中`==`所具备的多态性，即函数或操作符重载（讨论多态定义时，操作符可认为是前文所述的抽象的函数），被归类为 ad-hoc polymorphism，“Ad-hoc”的含义是特设的。这种多态指的是，根据不同类型的参数，去选择多态函数调用的具体实现；而每个函数实现是根据参数类型定制的，不同的实现之间不需要存在任何联系。

比如上述`int`和`string`所对应的`==`实现，是分别根据`int`和`string`定制的；编译器会根据`==`左右两个参数的类型去选择正确的`==`实现，且不同的实现根本互不相干（cmp指令和memcmp）。

对于函数重载版本间不同的参数列表，此处也解释为不同的类型。比如一个参数列表`(T1, T2, T3)`可以唯一地对应到一个类型`tuple<T1, T2, T3>`。

类比同样不支持自定义操作符重载的Java，其值类型和引用类型的`==`操作符分别具备多态性吗？

## Parametric Polymorphism
Go 中有另一种更显而易见的多态，那就是`map`和 slice 类型的操作。以 slice 为例，对于每个类型`T`，`[]T`都是不同的类型;但是`[]T`作为`len()`、`[]`、`append()`等函数的参数时，这些函数处理具体的值的逻辑都是一样的，而与具体的`T`无关。这一类函数逻辑不依具体类型而定的多态，被称为 parametric polymorphism。

比如，用 Go generics 的语法来描述 slice 与`len()`函数，它们可能长下面这样：
```go
// see reflect.SliceHeader
type Slice[T any] struct {
    data *T
	len  int
	cap  int
}

func (s *Slice[T]) Len() int {
    return s.len
}
```

不论`T`的实际类型是什么，每个由`slice[T]`实例化得到类型都有一个`len`字段，这是`Len()`方法能够运作的基础。由此引出，假设一个 parametric polymorphism 函数所能作用的参数类型的集合为`G`，那么`G`中的`T`都需要具备某种结构或性质。在很多具备泛型特性的编程语言中，类型约束（type constraint）就对应上述抽象的`G`，它作用于泛型函数的类型参数。

举一个简单的例子，`map`的键类型必须满足 Comparable 的约束，Comparable 类型能通过`==`和`!=`操作符比较相等性，否则实例化的`map`无法判断键值是否相等。根据 Go spec [2]：“... comparable types are boolean, numeric, string, pointer, channel, and interface types”。

当前的 Go generics 主要就是通过复用已有的 interface 特性来表达类型限制；而上述的 Comparable 应该只能内置实现。

## Java Generics
Java 大量使用了继承的概念，因此 Java generics 中的类型约束主要通过`extends`和`super`关键字来分别限制类型参数在继承链中的上界和下界。比如编写一个泛型函数`copy`来拷贝存放在`List`中的`T`类型对象，其函数签名如下：
```java
static <T> void copy(List<? extends T> src, List<? super T> dest)
```
`extends`关键字限定了`src`元素类型的上界，必须能从中读取到`T`；而`super`限定了`dest`元素类型的下界，必须能够存放一个类型`T`。背后的机制是，继承意味着可替换性——可以给父类参数提供子类对象，语法上表现为upcast，如果我转存T对象，那么必须从能upcast为T类型对象的list中读取；然后写入能够T能够upcast的类型的List中去。

这个例子行不行？

covariance：ArrayList<Number> 有 ArrayList<Integer> 父类子类关系吗？
contravariance：Function<Number, String> 是 Function<Integer, String> 的子类

## interface
Java or C++ inheritance is a mix
structural typing is a mix of inclusion and ad-hoc polymorphism.

[1] Wikipedia - Polymorphism (computer science). https://en.wikipedia.org/wiki/Polymorphism_(computer_science)  
[2] Go spec - Comparison operators. https://go.dev/ref/spec#Comparison_operators  
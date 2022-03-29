# Go 语言碎碎念：多态

说到多态，你会想到什么？是 Java 的类继承与虚函数表，还是 C++ 的模板与重载决议，抑或是 Go 的 interface？

> In programming language theory and type theory, polymorphism is the provision of a single interface to entities of different types or the use of a single symbol to represent multiple different types. [1]

编程语言中多态性（polymorphism）的不严谨的定义是，用同一个 interface 来使用多个不同的类型。简单地说，一个函数可以运作与不同类型的参数上。此处，“函数”不能仅仅是 `func`，而是一个更抽象的概念；而“类型”指的是变量的类型，即“type”（而非“kind”）的概念。这样的定义显得含糊，本文接下来以 Go 语言的多态为主线，详细地聊一下多态。

## 无处不在的多态
如今 C++/Java 流派的面向对象编程有最广泛的使用，因此可能有不少程序员对于多态或其他计算机科学词汇的理解都建立在该流派的语境之下。说到多态，脑子里可能出现的是基于类继承实现的动态派发。然而，可以想得再简单而基础一些：各种运算符通常也是多态的！比如，相等运算符`==`。

虽然 Go 不支持自定义运算符重载，但它的 `==` 实现却是多态的，可以类比于语言内置的重载。  
- 对于标量原始类型，如 `int` 或 `bool`，`==` 的语义是比较标量的值是否相等，其实现一般就是一条 CPU 的 cmp 指令。  
- 对于 `string`，`==` 比较两个字符串是否为相同的 rune 序列；又由于 Go 使用 UTF-8 编码，`==` 的实现很可能就是一条 `memcmp`（可以去这个网站查看对应的汇编代码：https://godbolt.org）。  
- 对于 `chan`，由于其具备引用语义，比较的是两个 `chan` 变量是否引用同一个 channel 实例。  
- 对于 `map` 或 slice，`==`仅可以把 `map` 或 slice 变量和 `nil` 值比较，检查是否进行过初始化。  


## Ad-hoc Polymorphism
Go 中`==`所具备的多态性，类似函数或操作符重载的效果（讨论多态定义时，操作符可认为是前文所述的抽象的函数），这种多态被归类为 ad-hoc polymorphism。“ad-hoc”的意为“特设的”，这种多态指的是，根据不同类型的参数，去选择一个多态函数的调用的具体实现；而每个函数实现是根据参数类型定制的，不同的实现之间不需要存在任何联系。

比如上述 `int` 和 `string` 所对应的 `==` 实现，是分别根据 `int` 和 `string` 定制的；编译器会根据 `==` 左右两个参数的类型去选择正确的 `==` 实现，且不同的实现根本互不相干（cmp指令和memcmp）。

对于函数重载版本间不同的参数列表，此处也解释为不同的类型。比如一个参数列表 `(T1, T2, T3)` 可以唯一地对应到一个类型 `tuple<T1, T2, T3>`。

类比同样不支持自定义操作符重载的 Java，其值类型和引用类型的 `==` 操作符分别具备多态性吗？

## Parametric Polymorphism
Go 中有另一种更显而易见的多态，那就是 `map` 和 slice 类型的操作。以 slice 为例，对于每个类型 `T`，`[]T` 都是不同的类型；但是 `[]T` 作为 `len()`、`[]`、`append()` 等函数的参数时，这些函数处理具体的值的逻辑都是一样的，而与具体的 `T` 无关。这一类函数逻辑不依具体类型而定的多态，被称为 parametric polymorphism，使用这种多态进行编程有个更为人熟悉的名字：泛型编程（generic programming）。

比如，用 Go 的泛型语法来描述 slice 与 `len()` 函数，它们可能长下面这样：
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

不论`T`的实际类型是什么，每个由 `slice[T]` 实例化得到类型都有一个 `len` 字段，这是 `Len()` 方法能够运作的基础。由此引出，假设一个 parametric polymorphism 函数所能作用的参数类型的集合为 `G`，那么 `G` 中的 `T` 都需要具备某种结构或性质。在很多具备泛型特性的编程语言中，类型约束（type constraint）就对应上述抽象的`G`，它作用于泛型函数的类型参数。

举一个简单的例子，`map` 的键类型必须满足 Comparable 的约束，Comparable 类型能通过`==`和`!=`操作符比较相等性，否则实例化的`map`无法判断键值是否相等。根据 Go spec [2]：“... comparable types are boolean, numeric, string, pointer, channel, and interface types”。

当前，Go 泛型主要通过复用已有的 interface 特性来表达类型限制；而上述的 Comparable 目前只能由编译器内置实现。

## Subtype Polymorphism
Java 大量使用了继承的概念，Java 中基于继承（此处包括`extends` superclass 和`implements` interface两类情况）和虚函数实现的多态可能是很多程序员最熟悉的一种多态，这种多态被称为 subtype polymorphism（又称 inclusion polymorphism）。广义的继承（inheritance in computer science），目的是代码复用，而非规定如何处理接口或类型关系。但在 Java 的语境中，subclass 不仅继承了 superclass 的代码，也继承了 superclass 的接口，使得 subclass 成为了 superclass 的 subtype（子类型，为了避免与“子类”混淆，本文此处还是用英文）。

Subtyping 的定义为一种类型间的关系：设有两个类型 S 和 T，S 是 T 的 subtype，则 S 类型的元素可以被“安全地”用于期待 T 的上下文中（如何满足“可替换性”还与具体的编程语言或编程准则相关）。因此，这个关系又被称作可替换性（substitutability）。显然，subtyping 与代码复用的重点不同，前者从 supertype 获取的是接口及其语义，后者获取的是代码实现（implementation），没有默认实现的 Java interface 相比 class inheritance 更接近 subtyping 的定义。

此处不再赘述基于继承和虚函数实现的多态，而是来讨论另一个问题：Java 于 SE 5 引入了泛型特性，那么 subtyping 和泛型相遇时会是什么情况呢？

Java 泛型中通过`extends`和`super`关键字表达类型约束，这两个关键字用于分别限制类型参数在其继承体系中的上界和下界。比如泛型的`Arrays.sort`其函数签名如下：

```java
public static <T> void sort(T[] a,
            int fromIndex,
            int toIndex,
            Comparator<? super T> c)
```
可以看到`Comparator<? super T>`的类型约束是`Comparator`类型参数下界为`T`，该Comparator必须能够消费（consume）两个类型`T`的变量进行比较；说得再肤浅一些，对于类型为`T`的变量`t1`和`t2`，表达式`c.compare(t1, t2)`是合法的。

这样设计的原因是，上文提到泛型函数依赖类型间具备共有结构和性质；而 Java 中正是由继承关系声明类型之间的 subtype 关系。又由于 Java 缺乏 structural typing 的能力（详细讨论见下文），所以 Java 泛型选择使用继承体系来约束类型参数的共同性质（subtype 拥有与 supertype 相兼容的接口），是一个自然而然的选择。

然而，上述`sort`例子的“肤浅”的说法暴露出一个问题，那就是 subtyping 在例中的体现仅仅停留在语法层面。面向对象编程中常提到的里氏替换原则（Liskov substitution principle）是语义上的，实际严谨的叫法是 behavioral subtyping，对于多态的行为作出了要求 [3]。正如 Oracle 的文档中对`Comparator<T>.compare()`的语义有详细的描述 [5]，但代码是否正确地实现了该语义却完全取决于程序员，编译器无从判断。

Alan Kay 曾在一场 talk 中谈到，他认为面向对象编程最脆弱的地方之一，是依赖程序员能按照接口的语义去编写代码，仅仅满足接口的输入、输出类型不能满足他对“面向对象编程”的定义 [4]。

*可选内容*：subtyping 和 parameteric polymorphism 还会遇到另一个问题：covariance 与 contravariance。Variance 指的是：

> how subtyping between more complex types relates to subtyping between their components? 蹩脚的翻译：复合类型之间的 subtyping 关系与构成它们的元素的 subtyping 有什么联系？

- covariance：`Integer extends Number`，那么`ArrayList<Number>` 与 `ArrayList<Integer>`具备超类子类关系吗？
- contravariance：`Function<Number, String>` 是 `Function<Integer, String>` 的子类？如果有个地方期待后者，我可以塞一个前者进去正常运行？

## Structural Subtyping

Go 和 Java 一样也有 `interface`，同样地用于指定一个类型的行为，且以一组方法签名来定义。比如`context.Context`的定义为以下一组方法：

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}
```

有了上文的铺垫，可以推断出 Go `interface` 的使用也是一种 subtype polymorphism。但它跟 Java `interface` 不一样之处是，Go `interface` 不需要显式地声明实现与 interface 的关系，而是由编译器进行检查。`context` 的实现分别使用了 `emptyCtx`， `cancelCtx`，`timerCtx` 和 `valueCtx` 等一系列结构体，但这些结构体的声明和定义都不需要提及 `context` interface。

Go 的方式被称为 structural typing，它基于类型的实际结构和定义，来决定类型之间的关系 [6]；而基于声明来决定的则被称为 nominal typing。Structural typing 又可进一步地分为静态和动态的——后者俗称 duck typing。在 Go 的情况中，假设有`type S`和`interface I`，Go 编译器通过检查`S`是否实现了`I`所定义的方法集，来判断`S`是否为`I`的 subtype。

Go 没有 Java 式的继承，并且把 nominal typing 改为 structural，能够减少代码间的耦合性。

与此同时，Go interface 也是 ad-hoc 的。

```go
type Iterator[T] interface {
  Next() bool
  Get()  T
}
```



举例 C++ template 也有类似功效

## 参考链接

[1] Wikipedia - Polymorphism (computer science). https://en.wikipedia.org/wiki/Polymorphism_(computer_science) <br/>[2] Go spec - Comparison operators. https://go.dev/ref/spec#Comparison_operators <br/>[3] Wikipedia - Liskov substitution principle. https://en.wikipedia.org/wiki/Liskov_substitution_principle <br/>[4] Seminar with Alan Kay on Object Oriented Programming (VPRI 0246). https://youtu.be/QjJaFG63Hlo<br/>[5] Comparator (Java Platform SE 8). https://docs.oracle.com/javase/8/docs/api/java/util/Comparator.html<br/>[6] Wikipedia - Structural type system. https://en.wikipedia.org/wiki/Structural_type_system<br/>


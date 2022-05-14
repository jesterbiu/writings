# 多态

- Why do we need polymorphism? using examples from the type-class paper. 
  - Polymorphism means that a value may have multiple types. Why need multiple types? Why need types at all?
- How do we achieve that? introduce various kinds of polymorphism in some clear distinct ways: parametric (C++ template - ish), ad-hoc (operator overloading).
- How do they evolve with each other? The expressiveness  is limited using only one kind of polymorphism at a time, thus almost all PLs employ more various 



## 无处不在的多态

说到多态，你会想到什么？是 Java 的类继承与虚函数表，还是 C++ 的模板与重载决议，抑或是 Go 的 interface？

> In programming language theory and type theory, polymorphism is the provision of a single interface to entities of different types or the use of a single symbol to represent multiple different types. [1]

在编程语言中，多态性（polymorphism）指的是用同一个接口提供多个不同类型的使用。

如今， C++/Java 流派的面向对象编程有最广泛的使用，因此可能有不少程序员对于“多态”（或其他计算机科学词汇的理解）都建立在该流派的语境之下。说到多态，脑子里可能出现的是基于类继承实现的动态派发。然而，不妨想得再简单而基础一些：把各种运算符视为函数的话，这些运算符通常也是多态的！比如，相等运算符 `==`：

```c++
auto b1 = 42 == (40 + 2);
auto b2 = vector<int>{1,2,3} == vector<int>{4,5,6};
auto b3 = string{"Edward Yang"} == string{"David Cronenberg"};
```

上面三条 C++ 语句分别比较了一对 `int`、`vector<int>` 和 `string` 变量的值是否相等。虽然每条语句中进行比较的类型并不相同（不同类型的使用），但它们都用着同样的接口：`==`。


## Ad-hoc Polymorphism
上述代码中 `==` 所具备的多态性是基于函数重载实现的，而函数重载是 *ad-hoc polymorphism*（特设多态）的主要形式——这种多态定义为一个函数可以以不同的类型的参数进行调用，并表现出不同的行为 [8]。

以函数重载来说，重载能把一个函数名关联到多个独立的实现，让编译器根据调用时的上下文选择合适的实现进行调用；而每个版本的实现是根据参数类型定制的，不同的版本间不需要存在其他的联系。因此，从另一角度来说，函数重载是一种编程语言内置的派发机制  [Wiki: ad-hoc poly]；与基于泛型的多态相比，它几乎完全独立于类型系统存在。

函数重载的优势是几乎不会给代码引入额外的耦合性（比起为了一个函数要继承一个类的情况，比如 Java 的 `Runnable`）；但它缺乏扩展性：它只能以单个函数为单位实现多态，没法表达一个多态函数集合构成的接口，也没法针对一组具有相同接口的类型共用同一份实现。

## Parametric Polymorphism
在实现通用的数据结构或算法时，我们通常希望仅实现每个数据结构或算法一次，但这份实现可以用于无限多种类型。比如编写一款动态数组，我们希望动态数组的代码能用于任意的元素类型，因为对于动态数组各种操作，无论是追加元素还是访问元素，其逻辑都与元素的类型无关。当然，这里还有个隐性的需求是，动态数组还应保留元素的类型信息。

鉴于以上的要求，人们（具体点，ML 语言的贡献者们）引入了类型变量来替代函数的参数类型；待到函数被调用时，才根据实参类型替换类型变量。这种做法引入了一种新的多态：*parametric polymorphism*（参数多态）。一个参数多态的函数能定义在一组（而非仅仅一个）类型上，且这一组类型通常共同具备某种接口。由于同一个函数定义不再拘束于单一类型，这种多态又得名 *generic programming*（泛型编程）。

以刚刚添加了泛型特性的 Go 语言为例，我们可以利用泛型实现迭代器，以统一的接口来访问容器：

```go
type iter[T any] interface {
	Next() bool
	Elem() T
}

type sliceIter[T any] struct {
	slice []T
	index int
}

func (iter *sliceIter[T]) Next() bool {
	iter.index++
	return iter.index < len(iter.slice)
}

func (iter *sliceIter[T]) Elem() T {
	return iter.slice[iter.index]
}

func fromSlice[T any](s []T) iter[T] {
	return &sliceIter[T]{s, -1}
}
```



// copy-pasted

第二种实现方式，其签名是 `(==) :: a -> a -> Bool`，其中 `a` 是一个类型变量。这是一种泛型的实现，即 parametric polymorphism。这种方式为每个类型都定义了相等比较，通常只需要按类型的类别来编写常数数量的版本，即可为任意的类型实现 `==`，而不需挨个单独实现。换个角度说，这种方式看似最省心，多个类型都能用同一份代码来判断其相等性，但它的缺点是可能导致 `==` 的语义难以理解或不符合期待。

接下来，以 Go 标准库的 `reflect.DeepEqual` 为例讨论这种实现。它虽然非泛型函数，但却具备类似泛型版本 `==` 的语法和语义——它接收两个任意类型的参数，并比较其“是否相等”。这个函数先以 Go 的类型类别的方式分别定义了基础类型的相等性，再基于这些定义来建立任意类型变量的相等定义 [2]。其中，类型的“类别”分别定义了数组、结构体、指针、函数、`interface` 等。

即使 Go 有着相对简单的类型系统，由于 `reflect.DeepEqual` 要处理任意类型的相等比较，其最终拥有比较复杂的语义。以下例子衍生自官方文档，不见得每个开发者**每时每刻都能流畅背诵**这么复杂的规则：

- Go 的函数变量只能与 `nil` 值进行比较，否则是编译错误；但 `reflect.DeepEqual` 会毫无怨言地帮你进行比较：

  ```go
  foo := func() error { return nil }
      
  // false
  ok1 := reflect.DeepEqual(foo, foo)
      
  // invalid operation: foo == foo (func can only be compared to nil)
  ok2 := (foo == foo)
  ```

- Go 所谓的空 slice 比较也很 tricky：指向长度为 0 的数组的指针和 nil 指针并不相等，因此：

  ```go
  // 再次敲响警钟，
  // 不要把pointer value用作object identity
  s1, s2 := []int{}, []int(nil)
  
  // true
  ok1 := len(s1) == len(s2)
  
  // false
  ok2 := reflect.DeepEqual(s1, s2)
  ```

显然，这一类泛型实现的语义并非总是合乎预期，而且伸缩性较差，我们仍然需要为一些类型特设地实现相等比较。

## Subtype Polymorphism
Java 大量使用了继承的概念，Java 中基于继承（此处包括`extends` superclass 和`implements` interface两类情况）和虚函数实现的多态可能是很多程序员最熟悉的一种多态，这种多态被称为 subtype polymorphism（又称 inclusion polymorphism）。

需要注意的是，广义的继承（inheritance in computer science），目的是代码复用，而非规定如何处理接口或类型关系；但在 Java 中，subclass 不仅继承了 superclass 的代码，也继承了 superclass 的接口，声明了 subclass 为 superclass 的 subtype（子类型，为 了避免与“子类”混淆，本文此处还是用英文）。 GOF 和 [10] 也讨论了这个区别，并将 subtype 称为接口继承（interface inheritance）[7]。

Subtype 的定义为一种类型间的关系：设有两个类型 S 和 T，S 是 T 的 subtype，则 S 类型的元素可以被“安全地”用于期待 T 的上下文中（如何满足“可替换性”还与具体的编程语言或编程准则相关）。这个关系意味着类型间的可替换性（substitutability）。显然，subtype 与代码复用的重点不同：前者从 supertype 获取的是接口及其语义，后者获取的是代码实现（implementation），没有默认实现的 Java interface 相比 class inheritance 更接近 subtype 的定义。

此处不再赘述基于继承和虚函数实现的多态，而是来讨论另一个问题：Java 于 SE 5 引入了泛型特性，那么 subtype 和泛型相遇时会是什么情况呢？

*可选内容*：subtype 还会遇到另一个问题：covariance 与 contravariance。Variance 指的是：

> how subtyping between more complex types relates to subtyping between their components? 蹩脚的翻译：复合类型之间的 subtyping 关系与构成它们的元素的 subtyping 有什么联系？

- covariance：`Integer extends Number`，那么`ArrayList<Number>` 与 `ArrayList<Integer>`具备超类子类关系吗？
- contravariance：`Function<Number, String>` 是 `Function<Integer, String>` 的子类？如果有个地方期待后者，我可以塞一个前者进去正常运行？

## Subtype is Cross-Cutting

目前提到的三种多态彼此不是互斥的，而是是会相互作用的，其中又以 subtype 体现得最为明显。接下来我们来讨论 subtype 与其他多态的交叉作用。前文提到，在 parametric polymorphism 或者泛型中，多态的运作依赖于类型的共同结构或接口；而 subtype 关系正好就定义了类型间的共同接口。

比如，Java 泛型中通过 `extends `和 `super` 关键字表达类型约束，这两个关键字用于分别限制类型参数在其继承体系中的上界和下界。这种使用 subtype hierarchy 上下界限来做约束的方法不是一种巧合 todo，而是编程语言学界自1980年代便开始研究的问题，它的学名是 bounded quantification [8]。

举例来说，比如，泛型的`Arrays.sort`其函数签名如下：

```java
public static <T> void sort(T[] a,
            int fromIndex,
            int toIndex,
            Comparator<? super T> c)
```
这个函数声明的类型形参为 `T`，可以看到 `Comparator<? super T>` 约束其类型参数的下界为 `T`，意味着该 `c` 必须能够消费（consume）类型 `T` 的变量进行比较。而参数 `T`  再次出现在 `Comparator` 的参数列表中，被称为todo？

然而，上述 `sort` 例子的引出出一个问题，那就是 subtype 在 type theory 中体现止步于编程语言的语法层面，而面向对象编程中常提到的里氏替换原则（Liskov substitution principle）却是语义上的。后续基于里氏原则提出的严谨概念是 behavioral subtyping，对于多态的行为作出了要求 [3]。在 `Comparator` 一例中，虽然 Oracle 的文档中对`Comparator<T>.compare()`的语义有详细的描述 [5]，但代码是否正确地实现了该语义却完全取决于程序员，编译器无从判断。

事实上，Alan Kay 曾在一场 talk 中谈到，他认为面向对象编程最脆弱的地方之一，是依赖程序员能按照接口的语义去编写代码，仅仅满足接口的输入、输出类型不能满足他对“面向对象编程”的定义 [4]。 

// copy-pasted: bounded quantification

函数式编程并非与命令式或面向对象等范式水火不容，他们总是有千丝万缕的联系。前文通过举例说明了 type class 在函数式语言 haskell 中作为类型约束机制的应用，本节将以 Java 泛型为例，来讨论 OOP 中基于接口继承的类型约束。

试想如下情景：我们在 OOP 中引入泛型编程时，如何约束类型参数？尝试为 Java 编写一个泛型的查找函数：

```Java
<T> int find(ArrayList<T> list, Predicate<T> pred)
```

抽象地说，泛型函数作用于拥有同样接口的类型。在本例中，`find` 要求 `pred.test` 方法必须拥有 `T -> boolean` 的签名。然而，如果 `T` 继承了某个 `class S` （或实现了某个  `interface S`），且这个断言也仅需用到 `S` 的接口，那么上述函数签名会阻碍使用 `Predicate<S>` ，而需要额外实现一个 `Predicate<T>`。

因此，假设有 `ArrayList<T>` 和 `Predicate<U>`，`find` 并不要求 `T == U`，而只需要 `T` 是 `U` 的 subtype 即可。先来看 Java 语境下 subtype 关系的定义：对于类型 `T`， 如果它继承 `class U` 或实现 `interface U`，则称 `T` 是 `U` 的 subtype，记为 `T <: U`。可见，subtype 总是拥有 supertype 的接口，对应着 type class 的成员类型总是拥有 type class 的接口。

既然，Java 中的继承体系就是一个 subtype 体系，那么 Java 或其他 OOP 语言便可以利用 subtype 关系来做类型约束。这个泛型与 subtype 两种多态发生交互的场景，最终引出了 *bounded quantification* 的定义：它基于 `<:` 的关系来约束类型参数。这个术语的严谨定义和相关讨论可以参考 *Types And Programming Language, Chapter 26* [3]，本文处理概念还是以感性认识为主。

在 Java 类型约束中：

- upper bound（`? <: T`）写作 `? extends T`；
- lower bound （`T <: ?`）写作 `? super T`。

因此， `find` 的签名应改为：

```java
<T> int find(ArrayList<T> list, Predicate<? super T> pred)
```



## Go Interface

Go 和 Java 一样也有 `interface` 用于说明一个类型的行为，并以一组方法签名来定义。比如 `context.Context` 的定义为以下一组方法：

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}
```

有了上文的铺垫，可以识别出 Go interface 的使用也是一种 subtype polymorphism。但它跟 Java interface 不一样之处是，Go interface 不需要显式地声明实现与 interface 的关系，而是由编译器进行检查。目前 `context` 的实现分别使用了 `emptyCtx`， `cancelCtx`，`timerCtx` 和 `valueCtx` 等一系列结构体，但这些结构体的声明中都不需要提及 `context` interface。

这种方式被称为 structural typing，它基于类型的实际结构和定义，来判断类型之间的关系 [6]；而基于声明来判断的则被称为 nominal typing。Structural typing 又可进一步地分为静态和动态的——后者俗称 duck typing。在 Go 的情况中，假设有 `type S` 和 `interface I`，Go 编译器通过检查 `S` 是否实现了 `I` 所定义的方法集，来判断 `S `是否为 `I `的 subtype。

在 subtype 的实现上，Go 比 C++ 或 Java 都更清晰：它明确地区分了代码复用和 subtype polymorphism 两种情况，其中前者在 Go 中主要通过组合的方式达到。继承是一种侵入式的特性，subclass 可以访问并依赖 superclass 的实现，这种耦合可能导致难以在不破坏 subclass 的前提下修改 superclass 的实现，这个问题被称为“The fragile base class problem” [10]。

在 C++ 社区，其实一直有使用非侵入式的技巧来实现 subtype 多态的主张，而非遵循传统地显式地使用继承特性，比如 Sean Parent 的一系列重磅级演讲（其中比较详细的一场见 [11]）。但其实早在90年代末就有 paper 总结过这种模式 [12]；我也写过一系列文章讨论过其实现 [13]。

说到 cross-cutting，Go interface 是 ad-hoc 的。

```go
type Iterator[T] interface {
  Next() bool
  Get()  T
}
```

## Type Class

第三种方式，签名变成了 `(==) :: a(==) -> a(==) -> Bool`。这种方式中，只有实现了 `==` 函数的类型才可以比较。此处的 `(==)` 不再仅仅表示某个函数，而是引入了一个超越类型本身的概念，这个概念说明了类型的属性或接口——`(==)` 意味着类型具备相等比较的接口。

原文就是在这样的背景下，提出了 *type class*。它是一种基于重载的多态机制，调和了完全依赖重载的方式一和完全泛型的方式二。论文中以 `Num` 为例来介绍 type class 的概念。`class Num` 的声明如下：

```haskell
class Eq a where
	(==)   :: a -> a -> Bool
class Eq a => Num a where
	(+)    :: a -> a -> a
	(*)    :: a -> a -> a
	negate :: a -> a
```

这段声明意味着：

1. 类型 `a` 属于 `class Eq`（这个 type class实现了相等比较），是 `a` 属于 `Num` 的必要条件。这一点展示了 type class 的组合能力，熟悉 OOP 继承的读者可以想象 type class 间基于组合关系形成的 hierarchy。
2. 在满足 1. 的基础上，如果类型 `a` 具备 `+`, `*` 和 `negate` 函数，则 `a` 属于 `class Num`。

我们通过 `instance` 关键字和定义 type class 要求的函数，将一个类型实现为某个 type class 的成员：

```haskell
instance Num Int where
	(+)    = addInt
	(*)    = mulInt
	negate = negInt
```

这段代码需要注意的另外一点是，type class 更多是一种 ad-hoc polymorphism。上面这段代码在为 `Int` 实现 `Num` 的接口时，虽然提到了后者，但这段代码是与类型 `Int` 的定义相分离的，这一点与 Java 或 C++ 为首的“OOP 语言”不同。这种 OOP 并不显式地区分代码和接口继承，因此实现接口意味着侵入式地修改类型的声明与定义；而为一个类型实现 type class 时，是结构化地为该类型定义了一组重载函数。

要使用 type class 时，可以声明一个函数仅于属于某 type class 的类型之上有定义。例如，我们声明多态的平方函数仅在属于 `Num` 的类型上有定义：

```haskell
-- 类型参数a只能是属于Num的类型
square :: Num a => a -> a -- 签名
square x =  x * x         -- 定义
```

小结一下，type class 可以视作一种类型的约束机制（原文的用词是 *bounded quantifier*）。Type class 的定义描述了类型的接口；而一个类型属于某个 type class，等价于这个类型具备该 type class 要求的接口。看到这里，熟悉 OOP 方法论的读者可能会说，interface 也是一组方法的集合，难道 type class 就是发明了一次 interface（如果不是又一次）？

其实，从编程语言的角度上来看，这两个概念还是泾渭分明的：type class 不是类型，但具体的 interface 通常是一个类型。OOP 里的 interface 是 subtype polymorphism 的一种应用：如果类型 S 是类型 T 的 subtype 关系（记为 S <: T），意味着可以在期待 T 对象之处，用 S 对象安全地替换 T 使用（此处关联里氏替换原则）。然而 type class 不是类型，那么某个类型与其所属的 type class 之间就没有所谓可替换性一说了。

## C++ Concept

原文中提到，用翻译来引入 type class 的问题之一，就是多态函数将会依赖动态派发的。然而，如果在编译时能通过类型实参来确定多态函数实际需要使用的函数，便能够针对性地生成多态函数的特化版本，以代码体积上升的代价换取运行时的性能。这方面最为人熟知的例子之一，就是 C++ template。

C++ template 类似于语言内置的宏或代码生成特性，对于一个 template，编译器会为其每一组不同的实参生成一份代码。而 C++ 20 标准引入的用于约束或限制 template 的特性 concepts，其语法和功能都能听见 type class 的回响。

比如，定义一个 concept，要求类型参数为支持相等比较的类型：

```c++
template <typename T>
concept eq = requires(T a, T b) {
    { a == b } -> std::convertible_to<bool>;
};
```

这个 `concept eq` 是不是很像上文中的 type class ` Eq` 呢：

```haskell
class Eq a where 
    (==) :: a -> a -> Bool
```

进一步地，再定义一个 concept 要求偏序比较，并将同时满足相等和偏序比较的类型称为可比较 `comparable`，可见 concept 也和 type class 或 Go interface 有易用的组合性：

```C++
template <typename T>
concept less = requires(T a, T b) {
    { a < b } -> std::convertible_to<bool>;
};

template <typename T>
concept comparable = eq<T> && less<T>;
```

再看 concept 的实现，由于 concept 是在 template 层面的，与传统意义的类继承几乎平行，我们也可以像实现 type class 那样以 ad-hoc polymorphism 的方式来实现concept——其实就是写重载函数。

```C++
struct person {
	std::string id;
    std::string name;
};

bool operator==(const person& lhs, const person& rhs) {
    return lhs.id == rhs.id;
}

bool operator<(const person& lhs, const person& rhs) {
    return lhs.name < rhs.name;
}

static_assert(comparable<person>); // true
```

可以看到，concept 的不同之处在于，template 之于类型有着类似 *structural typing* 语法（俗称 duck typing），所以实现某个 concept 不需要提及它；但反过来说，由于 concept 作用的 template 层面与其他语言特性所在的层次有距离感，其使用体验便不足 type class 的那般丝滑。事实上，C++ 长期被 template 把语言生态割裂得两极分化，是不利于社区发展的。

## 参考链接

[1] Wikipedia - Polymorphism (computer science). https://en.wikipedia.org/wiki/Polymorphism_(computer_science) <br/>[2] Go spec - Comparison operators. https://go.dev/ref/spec#Comparison_operators <br/>[3] Wikipedia - Liskov substitution principle. https://en.wikipedia.org/wiki/Liskov_substitution_principle <br/>[4] Seminar with Alan Kay on Object Oriented Programming (VPRI 0246). https://youtu.be/QjJaFG63Hlo<br/>[5] Comparator (Java Platform SE 8). https://docs.oracle.com/javase/8/docs/api/java/util/Comparator.html<br/>[6] Wikipedia - Structural type system. https://en.wikipedia.org/wiki/Structural_type_system<br/>[7] "Gang of Four" - Design Patterns: Elements of Reusable Object-Oriented Software.<br/>[8] Benjamin Pierce - Types And Programming Languages.<br/>[9] Mikhajlov, Leonid; Sekerinski, Emil - A Study of The Fragile Base Class Problem. http://www.cas.mcmaster.ca/~emil/Publications_files/MikhajlovSekerinski98FragileBaseClassProblem.pdf<br/>[10] Allen Holub - Why extends is evil. https://www.infoworld.com/article/2073649/why-extends-is-evil.html<br/>[11] Sean Parent - Better Code: Runtime Polymorphism. https://youtu.be/QGcVXgEVMJg<br/>[12] Chris Cleeland, Douglas C. Schmidt - External Polymorphism. https://www.dre.vanderbilt.edu/~schmidt/PDF/C++-EP.pdf<br/>[13] 深入浅出C++类型擦除（1） - 知乎. https://zhuanlan.zhihu.com/p/351291649<br/>[14]


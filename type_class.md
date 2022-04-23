# How to make ad-hoc polymorphism less ad hoc

## Disclaimer

丑话说在前：本 CRUD boy 没有 PL 背景，以下内容纯属阅读原文和相关材料后的自行发挥。发出来是希望引发读者对于编程语言设计问题的思考，而非作为研究的参考材料之用。

## Intro

Philip Wadler, Stephen Blott 于 1988 年发表了论文 *How to make ad-hoc polymorphism less ad hoc* [1]，这篇论文提出了 *type class* 的概念。Type class 是基于 ad-hoc polymorphism 的一种多态机制，解决了面向对象编程（OOP）、bounded quantification 的一些问题。

论文还描述了一种 type class 实现方式：对于没有 type class 但具备 Hindley-Milner 类型系统的编程语言，可以在编译时通过一系列算法把使用 type class 的代码翻译为没有 type class 但功能相等的代码。

进入讨论前，先来介绍下两种重要的多态类型 [4]：

- *parametric polymorphism*（参数多态），函数定义于一组类型上，对于其中的类型，函数的行为都一样——即泛型函数。
- *ad-hoc polymorphism*（特设多态），函数的定义是为每个类型特设的，例子是函数重载。

 然而，主要体现为函数重载的 ad-hoc poymorphism 存在一定局限。论文以编写类型的相等运算（equality comparison）为例，详细地讨论了这个问题。以下的例子 假设编程语言支持以定义重载函数的形式来定义操作符重载，即，以 `==`  作为相等运算的函数名。

## The Ad-hoc Way

第一种方式是，为每个需要相等比较的类型都写一份定义，按照 concrete type 来解析函数（类似 C++/Java 的函数重载）。这种方式下，假如为 `Int, Bool, MyRecord` 三种类型实现 `==`，需要编写以下三个函数：

```haskell
(==) :: Int      -> Int      -> Bool
(==) :: Bool     -> Bool     -> Bool
(==) :: MyRecord -> MyRecord -> Bool
```

对于本文的例子，可以这样简单地解读上述函数签名：最后一个类型是返回值的类型，其他的类型名是参数的类型。

这种方式非常与 C++ 的操作符重载非常相似，但实现起来可能需要有复杂的函数解析机制。以 C++ 标准库中的常用容器 `std::vector` 为例，我们能很自然地写出这样的代码：

```c++
std::vector<int> expect = { /*...*/ };
std::vector<int> actual = get_result();
return expect == actual;
```

然而， `vector` 和它重载的运算符 `==` 都是在 `std` 命名空间中：

```c++
namespace std {
    template<class T, class Alloc> 
    class vector;
    
    template<class T, class Alloc>
	bool operator==(const std::vector<T,Alloc>& lhs,
                    const std::vector<T,Alloc>& rhs);
}
```

那么，为什么声明一个 vector 变量需要前缀 `std::`，但调用 `==` 比较两个 `vector<int>` 变量时却不需要该前缀呢？这是因为 C++ 在函数调用解析中引入了 *argument-dependent lookup*，这个机制使得编译器除了在当前的命名空间查找合适的，还会去函数实参所属的命名空间去查找合适的重载函数。

## The Parametric Way

第二种方式实现相等比较是 `(==) :: a -> a -> Bool`，其中 `a` 是一个类型变量。这种是一种泛型的实现，通常只需要按类型的类别来编写常数数量的版本，即可为任意的类型实现 `==`。反过来说，任意多个类型都可能用了同一份代码来判断其相等性，这可能导致 `==` 拥有理解或不符期待的语义。

接下来，以 Go 标准库的 `reflect.DeepEqual` 为例讨论这种实现。它虽然非泛型函数，但却具备类似泛型 `==` 的语义——它接收两个任意类型的参数并比较其“是否相等”。这个函数先以 Go 的类型类别的方式分别定义了基础类型的相等性（分别定义了 array, struct, func, interface, map, slice 和 pointer，类型类别定义参见 `reflect.Kind`），再基于这些定义来建立任意类型变量的相等定义 [2]。

显然，这样的实现并不能覆盖所有的 case，我们仍然需要为一些类型特设地实现相等比较，而且泛型实现的语义并非总是合乎预期。

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

## Type Class

第三种方式，其签名变成了 `(==) :: a(==) -> a(==) -> Bool`：只有实现了 `==` 函数的类型才可以比较。如果此处冒号左边的 `(==)` 仍然表示某个函数，这句话显然是多余的。此处的 `(==)` 表示的应当是一个超越类型本身的概念，对应着类型的属性或接口。

原文就是在这样的背景下，提出了 *type class*，它是一种有限的特设多态，调和了完全特设化（方式一）和完全参数化（方式二）。论文中以 `Num` 为例来介绍 type class 的概念。`class Num` 的声明如下：

```haskell
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

## Tranlastion

原文实现 type class 的方式是，基于现有的语法和类型系统，把 type class 代码翻译成现有而不需要在编译器上做太多改造的（看起来就是代码生成）。这种翻译方式借鉴了 OOP 中 subtype polymorphism 的实现，即通过函数表进行派发（subtype 的定义在后文，此处感性认识为虚函数机制即可）。基本的算法是，为每个 type class 都生成一个新的类型用作其函数表（原文是 method dictionary），并让每个使用 type class 的函数都添加相应 type class 的函数表作为额外参数。

```haskell
-- 泛型的Num函数表，NumD类似一个type constructor
-- 实例化出类型NumDict，其字段分别Num所声明的三个函数
data NumD a = NumDict (a -> a -> a) (a -> a -> a) (a -> a)

-- Num的三个函数就对应NumDict的三个getter函数
add (NumDict a m n) = a -- 字段：加法函数
mul (NumDict a m n) = m -- 字段：乘法函数
neg (NumDict a m n) = n -- 字段：取反函数

-- Int的对应的函数表定义
numDInt :: NumD Int
numDInt :: NumDict addInt mulInt negInt

-- square 调用翻译
-- 签名：添加参数NumD a
square' :: NumD a -> a -> a 
-- 定义：用getter mul从表中获取类型a的乘法函数，再执行乘法
-- 用C语法写的话，类似：numDa.mul(x,x)
square' numDa x = mul numDa x x 
```

Go 泛型当前的实现，就很像上述的设计：语法上，用 `interface` 来声明类型参数的约束；实现上，为每个泛型函数都添加额外的一个字典参数（dictionary），其中包含根据类型实参匹配的函数和类型实参相关的信息 [5]。

## Bounded Quantification

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

## Reference

[1] Philip Wadler, Stephen Blott - How to make ad-hoc polymorphism less ad hoc. https://people.csail.mit.edu/dnj/teaching/6898/papers/wadler88.pdf<br/>[2] Luca Cardelli, Peter Wegner - On Understanding Types, Data Abstraction, and Polymorphism. https://classes.cs.uoregon.edu/14S/cis607pl/Papers/onunderstanding.a4.pdf<br/>[3] Pierce Benjamin - Types and Programming Languages.<br/>[4] https://pkg.go.dev/reflect#DeepEqual<br/>[5] Go 1.18 Implementation of Generics via Dictionaries and Gcshape Stenciling. https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries-go1.18.md

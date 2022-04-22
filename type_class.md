# 阅读笔记：How to make ad-hoc polymorphism less ad-hoc

## Disclaimer

丑话说在前：本 CRUD boy 没有 PL 背景，以下内容纯属阅读原文后的自行发挥。发出来是希望引发读者对于编程语言设计问题的思考，而非作为研究的参考材料之用。

## Intro

这篇论文提出了 *type class* 的概念，一种应用 ad-hoc polymorphism 的方式。Type class 可以视为一种函数重载的机制，解决了面向对象编程（OOP）、bounded quantification 的一些问题。

首先来介绍下两种重要的多态类型：

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

对于本文的例子，可以这样简单地解读上述函数签名：最后一个类型是返回值类型，其他的类型名是参数类型。

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

第二种方式实现相等比较是 `(==) :: a -> a -> Bool`，其中 `a` 是一个类型变量。这种是一种泛型的实现，要通常只需要按类型的类别来编写常数数量的版本，即可为任意的类型实现 `==`。反过来说，任意多个类型都可能用了同一份代码来判断其相等性，这可能导致 `==` 拥有理解或不符期待的语义。

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

原文就是在这样的背景下提出了 *type class*，它作为完全特设化（方式一）和完全参数化（方式二）的调和——一种有限的特设多态。论文中以 `Num` 为例来介绍 type class 的概念。`class Num` 的声明如下：

```haskell
class Eq a => Num a where
	(+)    :: a -> a -> a
	(*)    :: a -> a -> a
	negate :: a -> a
```

这段声明意味着：

1. 类型 `a` 属于 `class Eq`（这个 type class实现了相等比较），是 `a` 属于 `Num` 的必要条件。这一点展示了 type class 的组合能力，熟悉 OOP 继承的读者可以想象 type class 间基于组合关系形成的 hierarchy。
2. 在满足 1. 的基础上，如果类型 `a` 具备 `+`, `*` 和 `negate` 函数，则 `a` 属于 `class Num`。

，我们通过 `instance` 关键字和定义 type class 要求的函数，将一个类型实现为某个 type class 的成员：

```haskell
instance Num Int where
	(+)    = addInt
	(*)    = mulInt
	negate = negInt
```

在使用时，可以通过 type class 声明一个函数仅于一组类型之上有定义，这一组类型属于这个 type class（或同时属于这一组 type class）。例如，我们声明多态的加法、乘法和取相反值的函数仅在属于 `Num` 的类型上有定义：

```haskell
add :: Num a => a -> a -> a
mul :: Num a => a -> a -> a
neg :: Num a => a -> a
```

看到这里，熟悉 OOP 方法论的读者可能会说，interface 也是一组方法的集合，难道 type class 就是发明了一次 interface（如果不是又一次）？

其实，从编程语言设计的角度上来看，这两个概念还是泾渭分明的：type class 不是类型，但具体的 interface 通常是一个类型。OOP 里的 interface 是 subtype polymorphism 的一种应用：如果类型 S 是类型 T 的 subtype 关系（记为 S <: T），意味着可以在期待 T 对象之处，用 S 对象安全地替换 T 使用（此处关联里氏替换原则）。然而 type class 不是类型，那么某个类型与其所属的 type class 之间就没有所谓可替换性一说了。

## Tranlastion

原文还花了相当的篇幅来讲解“翻译技巧”：对于没有 type class 的概念的编程语言，可以在编译时通过一系列算法把使用了 type class 的代码翻译为没有 type class 但功能相等的代码。

这种翻译的方式借鉴了 OOP 中 subtype polymorphism 的实现方式：通过函数表进行派发。每个 type class 都对应一个新的类型作为其函数表（原文：method dictionary）。 那么，

## Go Generic

dict 实现

## Bounded Quantification

严谨定义和相关讨论参见 *Types And Programming Language, Chapter 26* [1]，这里仅给出一个感性认识。试想如下情景：我们在 OOP 中引入泛型编程时，如何约束类型参数？以 Java 为例：

```Java
<T> int find(ArrayList<T> list, Predicate<T> pred)
```

泛型函数作用于拥有同样接口的类型，这个 `find` 要求 `pred.test` 方法必须拥有 `T -> boolean` 的签名。然而，如果 `T` 拥有`class S` 或  `interface S` 的接口，且这个断言也仅需用到 `S` 的接口，那么上述函数签名会阻碍使用 `Predicate<S>` 而需要实现一个 `Predicate<T>`。

因此，设 `pred` 的类型是 `Predicate<U>`，`find`  并不要求 `T == U`，而只需要 `T <: U` 即可（对于 `T`， 如果它继承 `class U` 或实现 `interface U`，记为 `T <: U`）。显然，我们可以利用 Java 中的继承体系来表达这种约束关系——这引出了 bounded quantification，它基于 `<:` 的关系来约束类型参数。最终， `find` 的签名改为：

```java
<T> int find(ArrayList<T> list, Predicate<? super T> pred)
```

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

## Reference

[1] Pierce Benjamin - Types and Programming Languages.<br/>[2] https://pkg.go.dev/reflect#DeepEqual<br/>[3]

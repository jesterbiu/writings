# 多态

- How do we achieve that? introduce various kinds of polymorphism in some clear distinct ways: parametric (C++ template - ish), ad-hoc (operator overloading).
- How do they evolve with each other? The expressiveness  is limited using only one kind of polymorphism at a time, thus almost all PLs employ more various 
- TBD: fix ref 



## 无处不在的多态

> what is type and why do we need it?

起初，编程语言中并没有所谓“类型”的概念，程序员看到的数据都是存储器中的比特串（bit string）。此时，唯一接近类型的概念，就是“字”（word）——机器寄存器大小的定长比特串；比特串只是程序或数据的**二进制表示**（representation），其具体的含义则只能从对比特串的**解释**（interpretation）中获得。

然而，程序员会将串解读为字符、数字、指针或者指令等不同种类的数据，这些数据都有各自的用途或行为。这种分类的开始，也标志着类型系统的演化的开始。

一个“类型”的概念，主要包括对比特串的解释及其操作的限制。比如，如果一个比特串 `1100001` 被记为整型变量，那么应该按整型的编码去解释它的值，使用 two's complement 可得这个串对应的十进制整型为 97。如果按照其他类型解释它，则可能导致程序产生非预期的行为：按字符类型解释，使用 ASCII 会得到值为英文字符  `a`；按指针类型解释，可以得到地址为 `0x61`——用这个值去寻址可能会导致内存错误。

如今，大部分编程语言都有类型系统，类型系统很大程度上保证了程序是“类型安全”的——程序对于一个值的解释是前后一致的，不会在值上执行其类型不允许的操作（比如上述把整型直接当指针使用的作法）。

> Why do we need polymorphism, that a value may have multiple types? using examples from the type-class paper. 

然而，一些编程语言的类型系统只允许每个值有唯一的类型（如 Pascal, C），无论这个值是函数或是函数的参数类型。这种语言被称为**单态的**（monomorphic）。但在实际编程中，经常遇到需要为不同的类型实现同一种操作的情况，比如变量的相等性比较。

如果每个函数只允许拥有一种类型，每个需要比较相等的类型就需要各编写一个不重名的函数，代码相对冗长。比如，为整型、字符串和列表类型编写相等比较，可能就有如下的三个不同的函数声明：

```
// 单态的相等比较
bool equal_int(int a, int b)
bool equal_str(string a, string b)
bool equal_lst(list a, list b)
```


## Ad-hoc Polymorphism
对于这个问题，一种解决方案是**函数重载**（function overloading）。函数重载机制能把一个函数名关联到多个类型不同的实现，由编译器根据调用上下文选择合适的版本进行调用：

```
// 基于函数重载的相等比较
bool equal(int a, int b)
bool equal(string a, string b)
bool equal(list a, list b)
```

大部分编程语言以重载操作符的形式，在语言里内置了字节、整型及浮点数等原始数据类型的相等比较，并以库的形式提供字符串、列表等数据结构的相等比较。比如，在 C++ 中为一个类型实现相等比较的惯例作法为其重载 `==` 操作符（操作符重载本质也是函数重载）；更有甚者，在类型系统设计上从简的 Go 把字符串重载的 `==` 也内置到语言里去了。

因为上述的 `equal` 在重载后可以被解释为不同的函数类型，所以它是多态的（polymorphic）——函数重载在几乎不影响类型系统的前提下，为编程语言引入了多态性。

在编程语言理论里，函数重载是**特设多态**（ad-hoc polymorphism）的主要形式。特设多态定义是，同一个值能与多个单独定义、异构的类型相关联 [8]。这个定义很拗口，但联系到函数重载的使用场景，其实特设多态的意图是相对清晰的：对于一个多态的值或接口，其关联的每个类型都需要单独、专门地去实现（因此而得名“特设”）。

显然，当前大部分编程语言内置的运算符都属于特设多态。在相等函数一例中，相等函数名 `equal` 是一个多态的值；而比较整型、字符串和列表的函数类型则是三个单独定义的类型；其实现也是异构的，比如整型的比较可能只需要一条机器指令 `cmp`，但字符串的比较则可能基于专门设计的算法。

从实现角度来说，函数重载是一种编程语言内置的派发机制  [Wiki: ad-hoc poly]。编译器内部通常重命名不同的重载版本，在获取重载函数的实参类型后，把被重载函数名派发到合适的重载版本上去，而每个重载版本本身是单态的函数。

## Parametric Polymorphism
通用的数据结构或算法通常与其操作的元素类型无关，无论元素类型如何，它们都会对元素执行同样的操作——比如数组和排序算法，同一份定义应当可用于几乎任意的元素类型。

比如数组就通常被实现为一个多态的类型，它可能是整型数组、字符串数组或数组的数组等。然而，与每个类型单独实现的特设多态正好相反，这个场景需要一种对于每个类型都使用同一份定义的多态。

至此，便引出了**参数多态**（parametric polymorphism）的概念 [ML paper?]。这种多态通过把具体类型参数化，使得编程语言能够表达与具体类型无关的程序；待使用时，再根据上下文提供的类型实参，对参数多态的程序进行实例化（instantiation）得到最终的具体类型。

```
// 使用尖括号添加类型参数
type array<T>
void sort(array<T>)

// 使用时，实例化参数多态的类型
array<int> arr
sort(arr)
```

在实现方面，为编程语言引入参数多态通常需要不小的工作量。

首先，参数多态需要向类型系统中添加类型实例化的算法：对于给定的参数多态值和上下文，尝试实例化出合法的类型，而合法类型的定义通常与语言的规范和实现有深刻的关系。以 Java 为例，由于其泛型（generic，参数多态的别名）实现基于“类型擦除” [source?]，所以泛型类型的类型实参不得是原生类型，形如 `ArrayList<int>` 的就属于非法类型。

另一方面，参数多态的内存管理也是个困难的问题。由于参数多态对不同类型一视同仁，那么编译器或运行时如何确定程序的内存布局呢？

- 对于带垃圾回收的语言，通常倾向于用指针（引用）来多态类型的内存布局，并通过添加对象头等结构来保留类型信息。这种方式使得实例化产生的代码趋同，比如 Java 的泛型容器对不同元素类型都使用同一份字节码。
- 对于不带垃圾回收的语言，更倾向于给每个不同类型实参实例化一个版本，从而确定内存布局。比如，C++ 基于 template 来实现参数多态，编译器会对每一组不同的模板实参实例化一个具体的类型或函数，从而避免使用指针来覆盖不同的类型——因为指针不可地避免涉及对象的生命周期管理乃至堆分配，这些问题引入的复杂度和开销是难以估量的。

// 这部分是否放到后面去？

此外，需要考虑参数多态与语言中其他特性的交互：为什么C++ template或Java generic允许递归的类型定义（curiously reccuring template pattern)，为什么 Go generic不允许parameterized method？

## Subtyping

面向对象编程是基于“对象”概念的一种编程范式，这种对象捆绑了数据和操作这些数据的函数（称为方法），程序则由对象之间通过方法进行的交互构成。这种作法屏蔽了数据的实现细节，从而提高程序的封装性，以对抗软件日益增长的复杂度。

常见的面向对象语言基于“类”（class-based）的概念，类是对象的定义，包括对象的字段和方法，同时也是类型定义。类的导出方法（或公共方法）构成的集合又叫称为“接口”（interface），因为类的导出方法决定了如何使用该类的对象。

面向对象编程鼓励针对对象的接口进行编程，而不是对象的具体类型。例如，假设有一个函数 `process_file` 需要通过读取字符串来处理一个文件，且 `process_file` 不关注数据的源头，它只需要读取到字符串即可。那么，可以如此编写该函数的签名，及其入参的类接口：

```
void process_file(string_stream r) {
	for r.has_more() {
		process(r.read())
	}
}
class string_stream { 
	string read()
	bool has_more()
}
```

由于字符串可能来自磁盘文件，也可能来自网络传输，此处又分别用不同的类来处理这两种数据源：

```
class file { 
	string read() // 从磁盘文件读取字符串
	bool has_more()
	void close()
}
class http_response { 
	string read() // 从http回复读取字符串
	bool has_more()
	http_request get_request()
}
```

显然， `string_stream` 的接口是 `file` 或 `http_response` 接口的子集。那么，在程序语法上，可以认为把后二者的对象当作 `string_stream` 对象使用安全的。这便引出了面向对象编程中的“多态”：

>  需要区分的是，类继承（inheritance）和子类型多态是不同的概念。根据编程语言的具体情况，类继承可能同时继承了基类的接口和实现，但只有继承接口是为了实现子类型多态，而继承实现是为了代码复用（code reuse）。比如，C++ 或 Java 的类继承就是同时继承实现和接口的；相较之下，Java 的 interface 更单纯地实现了子类型多态。

子类型的定义中强调一种类型间的可替代关系：假设 S 是 T 的子类型，那么 S 类型的元素可以被安全地用于期待 T 的上下文中。这种安全的可替换性通常被称为里氏替换原则（Liskov substitution principle），由 Barbara Liskov 于 80 年代的一场演讲中而为人所知。

子类型多态中的“安全替换”有两层概念：一是在代码的语法上，子类型确实拥有父类型的接口；二是在其语义上，子类型的确实现了接口所声称的行为。举例来说，Java 的 `Comparator` 接口的文档中，对该接口 `compare()` 方法的语义有详细的描述，要求了分别在什么情况下返回正数和负数、方法必须具备传递性等 [5]。

然而，在大部分情况下，编译器或其他检查工具无从证明代码是否正确地实现了上述语义（完全的证明等价于 the halting problem），而只能校验程序语法的正确性。因此，接口语义一般以文档方式记载，是一种“君子协议”。Alan Kay （面向对象设计先驱之一）曾在一场 talk 中谈到面向对象编程最脆弱的地方之一，就是依赖程序员能按照接口的语义去编写代码；仅仅在语法上满足接口的输入、输出类型不能满足他对“面向对象编程”的定义 [4]。

*可选内容*：subtype 还会遇到另一个问题：covariance 与 contravariance。Variance 指的是：

> how subtyping between more complex types relates to subtyping between their components? 蹩脚的翻译：复合类型之间的 subtyping 关系与构成它们的元素的 subtyping 有什么联系？

- covariance：`Integer extends Number`，那么`ArrayList<Number>` 与 `ArrayList<Integer>`具备超类子类关系吗？
- contravariance：`Function<Number, String>` 是 `Function<Integer, String>` 的子类？如果有个地方期待后者，我可以塞一个前者进去正常运行？

## Subtyping is Cross-Cutting

目前，我们已经归纳出了三种重要的基本多态类型。然而，从它们的定义可以看出，它们彼此不是互斥的，而是可以相互作用的。其中，又数子类型多态与其他类型的交互最为人所熟知。

在参数化多态（泛型）中，代码能够不依赖具体类型编写具备一个前提，那就是其作用的类型都具备共同结构或接口。比如一个泛型的排序算法，它其实假设了被排序的元素是存在某种偏序关系的。因此，在保证类型安全的前提下，需要限制元素的类型为可以比较顺序的类型。

而在面向对象编程中，对象可以进行比较的性质通常以方法的形式表达，且子类型多态的定义正好就基于了类型间的共同接口（子类型会拥有父类型的接口）。因此，可以利用子类型关系来表达参数类型的要求或限制——它的正式名称是 *bounded quantification* [8]。

Java 的泛型便利用了子类型关系来进行类型参数约束：

- 要限制类型参数为类型 `T` 子类型，使用表达式 `? extends T `，基于继承体系形象地称为“上界”（upper bound）。
- 要限制类型 `T` 为类型参数的父类型，使用表达式 ` ? super T` ，称为“下界”（lower bound）。

举例来说，泛型的 `Arrays.sort` 的函数签名如下：

```java
public static <T> void sort(
            T[] a, int from, int to,
            Comparator<? super T> c)
```
这个函数对 `T` 类型的数组进行排序， `Comparator<? super T>` 限制了其类型参数的下界为 `T`，意味着这个 `Comparator` 必须能够消费（consume）类型 `T` 的变量进行比较。假设 `T` 是 `Integer`， 且 `Integer` 是 `Number` 的子类型，那么就可以使用一个 `Comparator<Number>` 对象作为最后一个参数（ `Number super Integer`）。

## Go Interface: ad-hoc x subtyping

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

## Type Class: ad-hoc x parametric

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


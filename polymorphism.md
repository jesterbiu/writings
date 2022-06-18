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

在编程语言理论里，函数重载是**特设多态**（ad-hoc polymorphism）的主要形式。特设多态定义是，同一个值能与多个单独定义、异构的类型相关联 [8]。换句话说，对于一个多态的值，其关联的每个类型都需要单独、专门地去实现，因此而得名“特设”。

显然，当前大部分编程语言内置的运算符都属于特设多态。在相等函数一例中，相等函数名 `equal` 是一个多态的值；而比较整型、字符串和列表的函数类型则是分别单独定义的类型；其实现也是异构的，比如整型的比较可能只需要一条机器指令 `cmp`，但字符串的比较则可能基于专门设计的算法。

从实现角度来说，函数重载是一种编程语言内置的派发机制  [Wiki: ad-hoc poly]。编译器内部通常重命名不同的重载版本，在获取重载函数的实参类型后，把被重载函数名派发到合适的重载版本上去，而每个重载版本本身是单态的函数。

## Parametric Polymorphism
通用的数据结构或算法通常与其操作的元素类型无关，无论元素类型如何，它们都会对元素执行同样的操作——比如数组和排序算法，同一份定义应当可用于几乎任意的元素类型。

比如数组就通常被实现为一个多态的类型，它可能是整型数组、字符串数组或数组的数组等。然而，与每个类型单独实现的特设多态正好相反，这个场景需要一种对于每个类型都使用同一份定义的多态。

至此，便引出了**参数多态**（parametric polymorphism）的概念 [ML paper, 1975]。这种多态通过把具体类型参数化，使得编程语言能够表达与具体类型无关的程序；待使用时，再根据上下文提供的类型实参，对参数多态的程序进行实例化（instantiation）得到最终的具体类型。

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

面向对象编程鼓励针对对象的接口进行编程，而不是对象的具体类型。例如，假设有一个函数 `process_file` 需要通过读取字符串来处理一个文件，且 `process_file` 不关注数据的源头，它只需要读取到字符串即可。那么，可以如此编写该函数：

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

由于 `string_stream` 的接口是 `file` 或 `http_response` 接口的子集，是不是可以把后两者的对象当作 `string_stream` 对象使用呢？在程序语法上，这应当是安全的，因为 `process_file` 将用到的方法仅仅是 `file` 或 `http_response` 的部分方法。这种基于接口覆盖的对象替换，引出了面向对象编程中的“多态”：子类型多态（subtyping）。

子类型多态定义为一种类型间的可替代关系：假设 S 是 T 的子类型（记作 S <: T），那么 S 类型的对象可以被安全地用于期待 T 类型对象的上下文中。从感性认识上来说，如果类型 S 的接口包含类型 T 的接口，则 S 是 T 的子类型。

类继承（inheritance）和子类型多态是不同的概念。根据编程语言的具体情况，类继承可能同时继承了基类的接口和实现，但只有继承接口是为了实现子类型多态，而继承实现是为了代码复用（code reuse）。因此，上面的例子中没有让 `file` 或 `http_response` 去“继承” `string_stream`，而仅仅是拥有同样的方法，以说明子类型多态的本质。

在编程语言概念上，继承跟子类型多态没有必然的关系。目前，C++ 或 Java 的类继承就是同时继承实现和接口的；相比类继承，Java 的 `implements interface` 机制更倾向于单纯地实现子类型多态；然而，Go 的 `interface` 实现了子类型多态，但它不是”继承“。

子类型多态所要求的“替换”有两层含义：一是在代码的语法上，子类型确实拥有父类型的接口，编程语言理论对子类型定义通常只要求这一层；二是在其语义上，子类型的确实现了接口所声称的行为——这第二层的概念有个更为人所知的名字：里氏替换原则（Liskov substitution principle），由 Barbara Liskov 的一场演讲中而为人所知 [Liskov 1986]。

例如，Java 的 `Comparator` 接口的文档中，对该接口 `compare()` 方法的语义有详细的描述，要求了分别在什么情况下返回正数和负数、方法必须具备传递性等 [5]。然而，在大部分情况下，编译器或其他检查工具无从证明代码是否正确地实现了上述语义（完全的证明等价于 the halting problem），而只能校验程序语法的正确性。因此，接口语义一般以文档方式记载，是一种“君子协议”。

面向对象设计先驱 Alan Kay 曾在一场演讲中表达过如下观点：依赖程序员能按照接口的语义去编写代码，是面向对象编程最脆弱的地方之一；仅仅在语法上满足接口的输入、输出类型不能满足他对“面向对象编程”的定义 [4]。

## Subtyping is Cross-Cutting

目前，我们已经归纳出了三种重要的基本多态类型。然而，从它们的定义可以看出，它们彼此不是互斥的。实际上，几乎所有的编程语言中都要处理不同的多态是如何交互的问题。其中，又数子类型多态与其他类型的交互最为典型。

在参数化多态（泛型）中，代码能够不依赖具体类型编写具备一个前提，那就是其作用的类型都具备共同结构或接口。比如一个泛型的排序算法，它其实假设了被排序的元素是存在某种偏序关系的。因此，在保证类型安全的前提下，需要限制元素的类型为可以比较顺序的类型。

而在面向对象编程中，对象可以进行比较的性质通常以方法的形式表达，且子类型多态的定义正好就基于了类型间的共同接口（子类型会拥有父类型的接口）。因此，可以利用子类型关系来表达参数类型的要求或限制——它的正式名称是 *bounded quantification* [8]。

Java 的泛型便利用了子类型关系来进行类型参数约束：

- 要限制类型参数为类型 `T` 子类型，使用表达式 `? extends T `，基于继承体系形象地称为“上界”（upper bound）。
- 要限制类型 `T` 为类型参数的父类型，使用表达式 ` ? super T` ，称为“下界”（lower bound）。

举例来说，泛型的 `sort` 的函数签名声明如下，它对元素类型为 `T` 的数组进行排序：

```java
<T> void sort(T[] a, Comparator<? super T> c)
```
参数 `Comparator c` 的用意是，调用者可以自主选择 `T` 对象的比较函数的实现，因此把比较函数包装成类对象作为一个参数传入。而它的类型约束 `? super T` 限制了其类型参数的下界为 `T`，即，类型实参需为 `T` 的父类型。因此，`Comparator` 必然能够比较子类型 `T` 的对象，因为 `T` 作为子类型，能够替代父类行对象被使用。

例如，设 `T = Integer` 且 `Integer <: Number` ，那么就可以用一个 `Comparator<Number>` 对象来调用 `sort`。填入类型实参后，`sort` 的函数签名是：

```Java
sort(Integer[] a, Comparator<Number> c)
```

*可选内容*：subtype 还会遇到另一个问题：covariance 与 contravariance。Variance 指的是：

> how subtyping between more complex types relates to subtyping between their components? 蹩脚的翻译：复合类型之间的 subtyping 关系与构成它们的元素的 subtyping 有什么联系？

- covariance：`Integer extends Number`，那么`ArrayList<Number>` 与 `ArrayList<Integer>`具备超类子类关系吗？
- contravariance：`Function<Number, String>` 是 `Function<Integer, String>` 的子类？如果有个地方期待后者，我可以塞一个前者进去正常运行？

## Go Interface: ad-hoc x subtyping

Go 语言也具备子类型多态的特性，但并不基于类继承，而是特设的。Go 语言的子类型多态基于 `interface` 实现，与基于类继承的面向对象语言相同，一个 `interface` 类型也是由一组方法签名所定义；如果某个类型实现了一个 `interface` 中的所有方法，那么这个类型的对象就可以被当作该 `interface` 对象来使用 。

以 `error interface` 为例，几乎所有的 Go 程序都会用到它：

```
type error interface {
	Error() string
}
```

我们可以这样实现一个具体的错误类型：

```
type baseErr struct {
	msg string
}

func (e baseErr) Error() string { 
	return e.msg 
}

func NewError(msg string) error {
	return baseErr{msg}
}
```

可以看到，`NewError` 能直接把 `baseErr` 对象当作 `error` 对象使用（作为返回值）。这是因为，Go 与类继承语言不同，只要实现 `Error() string` 的对象都可以被当作 `error` 对象使用，而不要求这个对象与 `error` 有继承关系。

Go 编译器在编译时会检查 `baseErr` 是否实现了 `error` 的方法集，来判断前者是否为 `error` 的子类型，而不需要代码来声明它们存在子类型关系。这种基于类型的结构或定义来决定类型之间的相容性或相等关系的类型系统，被称为 *structural type system* [6]。俗称为 *duck typing* 的类型特性，其实就是指在运行时检查类型相容性的 structural type system。

与之相对地，基于声明决定类型相容性的类型系统则被称为 *nominal type system*，比如 Java 就属于基于声明的。同样的情况下，Java 要求代码显式声明 `class err implements error` 。

 从另一方面看，Go的 `interface` 不仅是子类型多态，还是一种特设多态。假设在上述基本的 `baseErr` 之上，我们再添加一种 `messageErr`类型，可以为已有的 `error` 对象添加错误描述：

```
type messageError struct { 
	err error
	msg string 
}

func (e messageError) Error() string {
	return e.msg + ": " + e.err.Error()
}

func WithMessage(err error, msg string) error {
	return messageError{err, msg}
}
```

可以看到，`error`、`baseErr` 和 `messageErr` 三个类型间没有任何显式的关系，共性只是后两者拥有 `error` 的接口。

在子类型多态的实现上，Go 可能比 C++ 或 Java 都更清晰：它没有引入类继承的概念，明确地区分了代码复用和接口继承是两种不同的情况。静态类型的类继承语言自上世纪九十年代统治工业界以来，陆续出现了大量对于类继承的反思。类继承是一种侵入式的特性，如果子类可以访问并依赖父类的实现，便造成彼此实现上的耦合。这种耦合导致了“The fragile base class problem” ，即难以修改 superclass 的实现，但又不破坏子类 [10]。

在 C++ 社区，其实一直有使用非侵入式的技巧来实现 subtype 多态的主张，而非遵循传统地显式地使用继承特性，比如 Sean Parent 的一系列重磅级演讲（其中比较详细的一场见 [11]）。但其实早在90年代末就有 paper 总结过这种模式 [12]；我也写过一系列文章讨论过其实现 [13]。

## Typeclass: ad-hoc x parametric

不止子类型多态会与其他多态产生化学反应，特设多态和参数多态的碰撞也有不一样的火花。上文在参数多态的讨论中提到过，参数多态的函数依赖实际类型具备共同的接口，因此需要某种形式的类型约束。诸如 Java 或 C# 一类的面向对象的语言选择了使用子类型多态和类继承来做类型约束；而 1989 年的一篇论文提出了 *typeclass* 的概念，利用了特设多态来处理这个问题，并首先在知名的函数式语言 Haskell 中实现 [Ad-hoc 1989]。

以 Haskell 编写的 `member` 函数为例，它返回指定元素是否在输入的列表中，这要求元素必须能够进行相等比较；这是一个泛型函数，其中 `a` 是类型参数：

```haskell
class Eq a where
	(==)   :: a -> a -> Bool
	
member :: Eq a => [a] -> a -> Bool
```

我们先定义名为 `Eq` 的 typeclass，其包含一个可以进行相等比较的函数。Typeclass 的定义由一组函数签名构成，就像面向对象中的接口定义。使用时，函数声明里 `=>` 符号前面的表达式就是约束： `Eq a`，意味着 `member` 的参数类型 `a` 必须为 `Eq` 的成员。

然而，typeclass 虽然也叫”class“，但它并不是类型（ `a` 并非 `Eq` 的子类型）。它作用的领域可类比 Java 中类型约束表达式 `T extends Eq`（Haskell 和 Java 同为静态类型，这些表达式的求值位于编译期）；它的操作数（operand）为类型本身，而非类型实例化的对象。

再来看如何为一个类型实现 typeclass。Haskell 中用 `instance` 关键字来开始实现 typeclass，以下是列表（ `[a]`）的 `Eq ` 实现：

```haskell
instance Eq a => Eq [a] where
	[] == []     = True
	[] == y:ys   = False
	x:xs == []   = False
	x:xs == y:ys = (x==y) & (xs==ys)
```

Typeclass 的本质，是结构化的重载函数（因此说它是用特设多态解决参数多态的问题）。重载函数的实现与被重载的类型的定义是解耦的，因此也可以随时为一个类型实现某个 typeclass，无论这个类型是否定义于当前的代码模块。但是单个重载函数无法表达接口，而 typeclass 通过把重载函数组在一起并命名，从而可以定义接口，但又不额外引入类型之间的依赖关系。

要简单地实现 typeclass 也非常直接，编译器只需要为使用了 typeclass 约束的函数隐式地添加一个函数表的参数，并根据实参类型设置表中的函数。在编译时，可以根据实参推断出具体需要 typeclass 的哪个 `instance` ，因此泛型函数可以完全地保留实参类型信息。

更进一步地，编译期拥有完整类型信息，将提供更多的编译优化机会和实现的弹性。比如，编译器可以：

- 每个一个泛型函数都用同一份代码，在函数表中内嵌实参类型信息，并根据这些信息去申请内存、使用变量，这样做在编译产物体积和编译速度上都有优势。
- 编译器可以针对每一个不同的实参类型单独生成一个函数，并函数表内联，从而给编译器更多的上下文信息做优化，这样在代码的运行速度上更有优势。

介绍一点八卦：typeclass paper 的作者之一 Philip Wadler 参与过 Java 泛型和 Go 泛型的工作。Go 泛型设计过程中，由于没有形如 `Object` 类 top-type 的存在，保留了上述实现的弹性。Go第一个泛型版本（Go 1.18）选择了既要给泛型函数传入包含函数变量和类型信息的字典，又会在考虑 GC 策略的前提下为不同大小的类型生成特化的版本，具体的讨论见 [15]。

## C++ Concept

在 typeclass 实现讨论中提到，如果在编译时通过类型实参来确定泛型函数实际需要使用的函数，便能够针对性地生成特化版本，以代码体积上升的代价换取运行时的性能。在主流语言中，最为人熟知的例子莫过于 C++ 的模板 `template`，它类似于语言内置的宏或代码生成特性：对于一个 `template`，编译器会为其每一组不同的实参单独生成一份代码。

C++ 的模板被广泛地用于泛型函数的实现，但却长期缺乏类型约束的设施，导致 C++ 模板编程的必备技能是学会名为 *SFINAE* 的 idiom —— SFINAE 利用了（exploit）语言的编译规则，根据类型实参来选择是否进行模板实例化。

C++20 标准正式地引入了用于约束 `template` 参数的特性 concepts。然而，其语法和功能都有 typeclass 的影子。比如，还是限制类型参数需要支持 `==` 操作，用 `concept` 实现如下：

```c++
template <typename T>
concept eq = requires(T a, T b) {
    { a == b } -> std::convertible_to<bool>;
};
// concept长得很像一个啰嗦的typeclass！
// class Eq a where 
//   (==) :: a -> a -> Bool
```

进一步地，再定义一个 `concept less` 用于要求类型实现偏序比较，并将同时满足 `eq` 和 `less` 的类型称为 `comparable`：

```C++
template <typename T>
concept less = requires(T a, T b) {
    { a < b } -> std::convertible_to<bool>;
};

template <typename T>
concept comparable = eq<T> && less<T>;
```

可见 concept 也有易用的组合性（typeclass 或 Go `interface`  也具备这个有点）。

再来看 concept 的实现。由于 concepts 同样不是类型而将类型作为操作数，我们也可以像实现 typeclass 那样，以特设多态的方式来实现 `concept`，其实就是写重载函数。假设有一个 `person` 类，我们为其“实现” `comparable` 如下：

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

可以看到，concepts 和 typeclass 的对偶，就像 Go 和 Java 各自 `interface` 的对偶：模板之于类型有着类似 structural typing 的性质，代码中不需要提及需要满足的 `concept`；但 typeclass 的 `instance` 需要提及 typeclass 的名字。

但在另一方面，由于 concepts 在模板里才起作用，与其他语言特性便有割裂感，其使用体验便不及 typeclass 那样流畅。事实上，C++ 社区中似乎长期有一条以 template 为界的线，让语言生态两极分化。C++11 起，标准委员会花了很大努力提高模板编程的易用性，以融合模板与其他语言特性。

## 参考链接

[1] Wikipedia - Polymorphism (computer science). https://en.wikipedia.org/wiki/Polymorphism_(computer_science) <br/>[2] Go spec - Comparison operators. https://go.dev/ref/spec#Comparison_operators <br/>[3] Wikipedia - Liskov substitution principle. https://en.wikipedia.org/wiki/Liskov_substitution_principle <br/>[4] Seminar with Alan Kay on Object Oriented Programming (VPRI 0246). https://youtu.be/QjJaFG63Hlo<br/>[5] Comparator (Java Platform SE 8). https://docs.oracle.com/javase/8/docs/api/java/util/Comparator.html<br/>[6] Wikipedia - Structural type system. https://en.wikipedia.org/wiki/Structural_type_system<br/>[7] "Gang of Four" - Design Patterns: Elements of Reusable Object-Oriented Software.<br/>[8] Benjamin Pierce - Types And Programming Languages.<br/>[9] Mikhajlov, Leonid; Sekerinski, Emil - A Study of The Fragile Base Class Problem. http://www.cas.mcmaster.ca/~emil/Publications_files/MikhajlovSekerinski98FragileBaseClassProblem.pdf<br/>[10] Allen Holub - Why extends is evil. https://www.infoworld.com/article/2073649/why-extends-is-evil.html<br/>[11] Sean Parent - Better Code: Runtime Polymorphism. https://youtu.be/QGcVXgEVMJg<br/>[12] Chris Cleeland, Douglas C. Schmidt - External Polymorphism. https://www.dre.vanderbilt.edu/~schmidt/PDF/C++-EP.pdf<br/>[13] 深入浅出C++类型擦除（1） - 知乎. https://zhuanlan.zhihu.com/p/351291649<br/>[14] <br/>[15] Go generic https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries-go1.18.md

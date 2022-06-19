# 多态



## Ad-hoc Polymorphic

起初，编程语言中并没有所谓“类型”的概念，程序员看到的数据都是存储器中的比特串（bit string）。此时，唯一接近类型的概念，就是“字”（word）——机器寄存器大小的定长比特串；比特串只是程序或数据的**二进制表示**（representation），其具体的含义则只能从对比特串的**解释**（interpretation）中获得。

然而，程序员会将串解读为字符、数字、指针或者指令等不同种类的数据，这些数据都有各自的用途或行为。这种分类的开始，也标志着类型系统的演化的开始。

一个“类型”的概念，主要包括对比特串的解释及其操作的限制。比如，如果一个比特串 `1100001` 被记为整型变量，那么应该按整型的编码去解释它的值，使用二进制补码可得这个串对应的十进制整型值为 97。如果按照其他类型解释它，则可能导致程序产生非预期的行为：使用 ASCII 编码按字符类型解释，会得到英文字符  `a`；按指针类型解释，可以得到地址为 `0x61`——但使用这个值去寻址可能会导致内存错误。

如今，大部分编程语言都有类型系统，类型系统很大程度上保证了程序是“类型安全”的，即程序对于一个值的解释是前后一致的，不会在值上执行其类型不允许的操作。像上述把整型直接当指针使用的作法就属于类型不安全的行为。

然而，一些编程语言的类型系统只允许每个值有唯一的类型（如 Pascal, C），无论这个值是函数或是函数的参数类型。这种语言被称为**单态的**（monomorphic）。但在实际编程中，经常遇到需要为不同的类型实现同一种操作的情况，比如变量的相等性比较。

如果每个函数只允许拥有一种类型，每个需要比较相等的类型就需要各编写一个不重名的函数，代码相对冗长。比如，为整型、字符串和列表类型编写相等比较，可能就有如下的三个不同的函数声明：

```
// 单态的相等比较
bool equal_int(int a, int b)
bool equal_str(string a, string b)
bool equal_lst(list a, list b)
```

显然，上述的做法显得啰嗦而死板了。对于这个问题，一种解决方案是**函数重载**（function overloading）。函数重载机制能把一个函数名关联到多个类型不同的实现，由编译器根据调用上下文选择合适的版本进行调用：

```
// 基于函数重载的相等比较
bool equal(int a, int b)
bool equal(string a, string b)
bool equal(list a, list b)
```

因为上述的 `equal` 虽然是只是一个符号，在重载后却可以被解释为不同的函数类型，所以它是**多态的**（polymorphic）[1] ——函数重载在几乎不影响类型系统的前提下，为编程语言引入了多态性。

大部分编程语言以重载操作符的形式（操作符重载本质就是函数重载），在语言里内置了字节、整型及浮点数等原始数据类型的相等比较，并以库的形式提供字符串、列表等数据结构的相等比较。比如，在 C++ 中为一个类型实现相等比较的惯例作法为其重载 `==` 操作符；更有甚者，在类型系统设计上一切从简的 Go 语言，把字符串重载的 `==` 也内置到语言里去了。

在编程语言理论里，函数重载是**特设多态**（ad-hoc polymorphism）的主要形式。特设多态定义的是，多态函数在一组不同类型的参数上有定义，对于不同类型的实参，会调用不同的函数实现，从而根据入参类型获得不同的行为，且这些实现是彼此独立的 [2,3]。换句话说，对于一个特设多态函数，其关联的每个类型都需要单独、专门地去实现，因此而得名“特设”。

显然，当前大部分编程语言内置的运算符都属于特设多态。在 `equal` 一例中，比较整型、字符串和列表的函数则是根据入参类型区分的独立实现。整型的比较可能只需要一条机器指令 `cmp`，但字符串的比较则可能基于专门设计的算法。

从实现角度来说，函数重载是一种编程语言内置的派发机制  [3]。编译器内部通常重命名不同的重载版本，在获取重载函数的实参类型后，把被重载函数名派发到合适的重载版本上去，而每个重载版本本身是单态的函数。

## Parametric Polymorphism
通用的数据结构或算法通常与其操作的元素类型无关，无论元素类型如何，它们都会对元素执行同样的操作——比如数组或排序算法，同一份实现应当可用于几乎任意的元素类型。以数组为例，它可能是整型数组、字符串数组或数组的数组等。然而，与每个类型单独实现的特设多态正好相反，这个场景需要一种对于每个类型都使用同一份定义的多态。

为此，编程语言中出现了**参数多态**（parametric polymorphism）的概念 [2]。这种多态通过把具体类型参数化，使得编程语言能够表达与具体类型无关的程序；待使用时，再根据上下文提供的类型实参，对参数多态的程序进行实例化（instantiation）得到最终的具体类型。

```
// 使用尖括号添加类型参数
type array<T>
void sort(array<T>)

// 使用时，实例化参数多态的类型
array<int> arr
sort(arr)
```

在实现方面，为已有的编程语言引入参数多态通常需要不小的工作量，支持参数多态，需要向类型系统中添加类型实例化的算法：对于给定的参数多态值和上下文，尝试实例化出合法的类型。而在具体的实现中，所谓“合法类型”的定义与编程语言的规范和实现有深刻的关系。以 Java 为例，由于其参数多态通过“类型擦除”实现，所以其类型实参不得是原生类型，形如 `ArrayList<int>` 的类型就属于非法类型 [5]。

> 在面向对象语言中，通常把参数多态称为“泛型”（generics），这两个词的含义是基本一致的。

像 Java、C# 和最近的 Go 都是在发布之后才添加对泛型的支持，因为这涉及语言类型系统各方面的修改。上述三个语言添加泛型的时间分别是 2004 年的 J2SE 5.0（发布 9 年后），2005 年的 C# 2.0 （发布 3 年后），2022 年的 Go 1.18 （发布 12 年后）。然而，参数多态并非是什么 exotic 的特性，它早在 1970 年代就在 ML 语言中被实现了，ML 的贡献者们发布了论文专门讨论这个问题 [4]。

另一方面，参数多态的内存管理也是个难题。由于参数多态对不同类型一视同仁，那么编译器或运行时环境如何确定程序的内存布局呢？

- 有垃圾回收器的语言，通常倾向于用指针（或引用）来管理参数类型对象，并通过添加对象头等结构来保留类型信息，从而消除不同的类型对象大小不同的影响，让函数在实例化后拥有近乎相同的内存布局。更有甚者，如 Java 的泛型容器，无论元素类型都使用同一份字节码。
- 无垃圾回收器的语言，更倾向于给每个不同类型实参实例化一个版本，从而根据内存布局。比如，C++ 基于模板（template）来实现参数多态时，编译器会对每一组不同的模板实参实例化一个具体的类型或函数，从而避免使用指针来覆盖不同的类型——因为指针不可地避免涉及对象的生命周期管理乃至堆分配，这些问题引入的复杂度和开销是巨大的。

## Subtyping

面向对象编程是基于“对象”概念的一种编程范式，这种对象捆绑了数据和操作这些数据的函数（称为方法），程序则由对象之间通过方法进行的交互构成。这种做法屏蔽了数据的实现细节，从而提高程序的封装性，以对抗软件日益增长的复杂度。

常见的面向对象语言基于“类”（class-based）的概念，类是对象的定义，包括对象的字段和方法，同时也是类型定义。类的导出方法（或公共方法）构成的集合又叫称为“接口”（interface），因为类的导出方法决定了如何使用该类的对象。

面向对象编程鼓励针对对象的接口进行编程，而不是对象的具体类型。例如，假设有一个函数 `process_file` 需要通过读取字符流来处理一个文件，且 `process_file` 不关注数据的源头，它只需要读取到字符串即可。那么，可以如此编写该函数：

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

由于字符流可能来自磁盘文件，也可能来自网络传输，因此用不同的类来处理这两种数据源：

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

由于 `string_stream` 的接口是 `file` 或 `http_response` 接口的子集，是不是可以把后两者的对象当作 `string_stream` 对象使用呢？在程序语法上，这应当是安全的，因为 `process_file` 会用到的方法，是 `file` 或 `http_response` 的方法的子集。这种基于接口兼容的对象替换，就是面向对象编程中常说的“多态”：**子类型多态**（subtyping）。

子类型多态定义为一种类型间的可替代关系：如果 S 类型的对象可以被安全地用于期待 T 类型对象的上下文中， 那么称 S 是 T 的子类型（记作 S <: T）[2]。从感性认识上来说，如果类型 S 的接口包含类型 T 的接口，则 S 是 T 的子类型。

类继承（inheritance）和子类型多态是不同的概念。根据编程语言的具体情况，类继承可能同时继承了基类的接口和实现，但只有继承接口是为了实现子类型多态，而继承实现是为了代码复用（code reuse）[6]。因此，上面的伪代码没有让 `file` 或 `http_response` 去“继承” `string_stream`，而仅仅是拥有同样的方法，以说明子类型多态的本质。

虽然，在以 C++ 或 Java 为代表主流面向对象语言中，“类继承”是同时继承实现和接口的，但在编程语言概念上，继承跟子类型多态没有必然的关系。相比类继承，Java 的 `implements interface` 机制更倾向于单纯地实现子类型多态。此外，Go 语言的 `interface` 实现了子类型多态，但它并不是通过”继承“实现的。

子类型多态所要求的“替换”有两层含义：一是在代码的语法上，子类型确实拥有父类型的接口，编程语言理论对子类型定义通常只要求这一层；二是在其语义上，子类型的确实现了接口所声称的行为——这第二层的概念有个更为人所知的名字：里氏替换原则（Liskov substitution principle） [7]。

例如，Java 的 `Comparator` 接口的文档中，对该接口 `compare()` 方法的语义有详细的描述，要求了分别在什么情况下返回正数和负数、方法必须具备传递性等 [8]。然而，在大部分情况下，编译器或其他检查工具无从证明代码是否正确地实现了上述语义（完全的证明等价于 the halting problem），而只能校验程序语法的正确性。因此，接口语义一般以文档方式记载，是一种“君子协议”。

面向对象设计先驱 Alan Kay 曾在一场演讲中表达过如下观点：依赖程序员能按照接口的语义去编写代码，是面向对象编程最脆弱的地方之一；仅仅在语法上满足接口的输入、输出类型不能满足他对“面向对象编程”的定义 [9]。

## Java Bounded Type Parameter

目前，我们已经归纳出了三种重要的基本多态类型。然而，从它们的定义可以看出，它们彼此不是互斥的。实际上，几乎所有的编程语言中都要处理好不同的多态的交互。其中，又数子类型多态与其他类型的交互最为典型。

对于一个参数多态函数，其合法的类型实参都具备某种共同结构或接口，是它不依赖具体类型而编写的前提。比如一个泛型的排序函数，它其实假设了被排序的元素是存在某种偏序关系的。因此，在保证类型安全的前提下，排序函数需要限制元素的类型为可以比较顺序的类型。

在面向对象编程中，对象可以进行比较的性质通常以方法的形式表达，而子类型多态的定义正好就基于了类型间的共同接口（子类型拥有父类型的接口）。因此，可以利用子类型关系来表达参数类型的限制，这种构建的正式名称叫 ***bounded quantification*** [2]。

Java 的泛型便利用了子类型关系来进行类型参数约束：

- 要限制类型参数为类型 `T` 子类型，使用表达式 `? extends T `，基于继承体系形象地称为“上界”（upper bound）。
- 要限制类型 `T` 为类型参数的父类型，使用表达式 ` ? super T` ，称为“下界”（lower bound）。

举例来说，泛型的 `sort` 的函数签名声明如下，它对元素类型为 `T` 的数组进行排序：

```java
<T> void sort(T[] a, Comparator<? super T> c)
```
参数 `Comparator c` 的用意是，调用者可以自定义 `T` 对象的比较函数的实现，因此把比较函数包装成类对象作为一个参数传入。而它的类型约束 `? super T` 限制了其类型参数的下界为 `T`，即类型实参需为 `T` 的父类型；又因为 `T` 是作为子类型，所以`Comparator` 便必然能够比较子类型 `T` 的对象。

例如，设 `T = Integer` 且 `Integer <: Number` ，那么就可以用一个 `Comparator<Number>` 对象来调用 `sort`。填入类型实参后，`sort` 的函数签名是：

```Java
sort(Integer[] a, Comparator<Number> c)
```

更进一步地， `Comparator<Number>`  和 `Comparator<? super Integer>` 的关系又是什么呢？对于泛型类，如果作为参数的类型之间存在子类型关系，意味着实例化形成的泛型类之间存在子类型关系吗？

比如，已知 Java 中存在 `Integer <: Number`，那么 `ArrayList<Integer> <: ArrayList<Number>` 是否成立？有些编程语言（如 OCaml）是允许后者成立的，这种组合类型间的子类型关系与元素子类型关系方向一致的情况，被称为 ***covariance***。这种情况下，Java 不允许它们直接建立子类型关系，而是需要使用 bounded quantification 来间接地建立：`ArrayList<Integer> <: ArrayList<? extends Number>`。

然而，上述规则似乎有瑕疵：假设有 `Function<Integer, String>`，且存在一个类型 `T` 满足`Integer <: T <: Number`，那么 `Function<Integer, String> <: Function<? extends Number, String>` 就不成立。

恰恰相反，由于 `Integer` 作为 `Function` 的入参类型，且存在 `Integer <: Number`，那么 `Function<Number, String>` 应当能够”替代“  `Function<Integer, String>` 使用才对——输入`Integer` 对象，但把它当成 `Number` 使用，是符合子类型关系的。这种组合类型之间的可替代关系与元素类型的相反，被称为 ***contravariance***，泛型函数的入参类型约束问题是描述该概念的典例。

因此，如果用 bounded quantification 来描述 `Function` 的入参约束，应当为 `? super Integer`，最终得出 `sort` 一例中存在 `Comparator<Number> <: Comparator<? super Integer>` 的关系。

## Go Interface

Go 语言也具备子类型多态的特性，但并不基于类继承，而是特设的。Go 语言的子类型多态基于 `interface` 实现，与基于类继承的面向对象语言相同，一个 `interface` 类型也是由一组方法签名所定义；如果某个类型实现了一个 `interface` 中的所有方法，那么这个类型的对象就可以被当作该 `interface` 对象来使用 。

以 `error interface` 为例，几乎所有的 Go 程序都会用到它：

```go
type error interface {
	Error() string
}
```

我们可以这样实现一个具体的错误类型：

```go
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

可以看到，`NewError` 能直接把 `baseErr` 对象当作 `error` 对象使用（作为返回值）。这是因为，Go 与类继承语言不同，只要实现 `Error() string` 的对象都可以被当作 `error` 对象使用，而不要求这个对象与 `error` 有任何继承关系。Go 依靠编译器在编译时检查 `baseErr` 是否实现了 `error` 的方法集，来判断前者是否为 `error` 的子类型，而不需要代码来声明它们存在子类型关系。

这种基于类型的结构或定义来决定类型之间的相容性或相等关系的类型系统，被称为 ***structural type system*** [10]。俗称为 ***duck typing*** 的类型特性，其实就是指在运行时检查类型相容性的 structural type system。

与之相对地，基于声明决定类型相容性的类型系统则被称为 ***nominal type system***，比如 Java 就属于基于声明的。同样的情况下，Java 要求代码显式声明 `class err implements error` 。

 从另一方面看，Go的 `interface` 不仅是子类型多态，还是一种特设多态。假设在上述基本的 `baseErr` 之上，我们再添加一种 `messageErr`类型，可以为已有的 `error` 对象添加错误描述：

```go
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

可以看到，`baseErr` 与 `messageErr` 的 `Error() string` 的实现是彼此独立的，因此 Go `interface` 也被认为是某种特设多态——`baseErr` 或 `messageErr` 可以视作是”重载“了 `error` 的方法。

在子类型多态的实现上，Go 可能比 C++ 或 Java 都更清晰：它没有引入类继承的概念，明确地区分了代码复用和接口继承是两种不同的情况。静态类型的类继承语言自上世纪九十年代统治工业界以来，陆续出现了大量对于类继承的反思。类继承是一种侵入式的特性，如果子类可以访问并依赖父类的实现，便造成彼此实现上的耦合。这种耦合导致了“The fragile base class problem” ，即难以修改基类的实现却不破坏子类 [11]。

在 C++ 社区，其实也一直有使用非侵入式的技巧来实现子类型多态的主张，而摈弃传统的继承特性，其效果就很像 Go `interface`。比如 Sean Parent 的一系列重磅级演讲 [13]，甚至早在 1990 年代末就有 paper 总结过这种模式 [14]。

## Haskell Typeclass

不止子类型多态会与其他多态产生化学反应，特设多态和参数多态的碰撞也有不一样的火花。上文在参数多态的讨论中提到过，参数多态的函数依赖实际类型具备共同的接口，因此需要某种形式的类型约束。诸如 Java 或 C# 一类的面向对象的语言选择了使用子类型多态和类继承来做类型约束；而 1989 年的一篇论文提出了 *typeclass* 的概念，利用了特设多态来处理这个问题，并首先在知名的函数式语言 Haskell 中实现 [16]。

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

介绍一点八卦：typeclass paper 的作者之一 Philip Wadler 参与过 Java 泛型和 Go 泛型的工作。Go 泛型设计过程中，由于没有形如 `Object` 类 top-type 的存在，保留了上述实现的弹性。Go第一个泛型版本（Go 1.18）选择了既要给泛型函数传入包含函数变量和类型信息的字典，又会在考虑 GC 策略的前提下为不同大小的类型生成特化的版本，具体的讨论见 [17]。

## C++ Concept

在 typeclass 实现讨论中提到，如果在编译时通过类型实参来确定泛型函数实际需要使用的函数，便能够针对性地生成特化版本，以代码体积上升的代价换取运行时的性能。在主流语言中，最为人熟知的例子莫过于 C++ 的模板 `template`，它类似于语言内置的宏或代码生成特性：对于一个 `template`，编译器会为其每一组不同的实参单独生成一份代码。

C++ 的模板被广泛地用于泛型函数的实现，但却长期缺乏类型约束的设施，导致 C++ 模板编程的必备技能是学会名为 *SFINAE* 的 idiom —— SFINAE 利用了（exploit）语言的编译规则，根据类型实参来选择是否进行模板实例化 [18]。

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

## 结语

On Understanding Types, Data Abstraction, and Polymorphism

## References

[1] Wikipedia - Polymorphism (computer science). https://en.wikipedia.org/wiki/Polymorphism_(computer_science) <br/>[2] Benjamin Pierce - Types And Programming Languages.<br/>[3] Wikipedia - Ad hoc polymorphism. https://en.wikipedia.org/wiki/Ad_hoc_polymorphism<br/>[4] R. Milner, L. Morris, M. Newey - A Logic for Computable Functions with Reflexive and Polymorphic Types. <br/>[5] Java Generics: Past, Present and Futurit. https://youtu.be/LEAoMMEIUXk<br/>[6] "Gang of Four" - Design Patterns: Elements of Reusable Object-Oriented Software.<br/>[7] Wikipedia - Liskov substitution principle. https://en.wikipedia.org/wiki/Liskov_substitution_principle <br/>[8] Comparator (Java Platform SE 8). https://docs.oracle.com/javase/8/docs/api/java/util/Comparator.html<br/>[9] Seminar with Alan Kay on Object Oriented Programming (VPRI 0246). https://youtu.be/QjJaFG63Hlo<br/>[10] Wikipedia - Structural type system. https://en.wikipedia.org/wiki/Structural_type_system<br/>[11] Mikhajlov, Leonid; Sekerinski, Emil - A Study of The Fragile Base Class Problem. http://www.cas.mcmaster.ca/~emil/Publications_files/MikhajlovSekerinski98FragileBaseClassProblem.pdf<br/>[12] Allen Holub - Why extends is evil. https://www.infoworld.com/article/2073649/why-extends-is-evil.html<br/>[13] Sean Parent - Better Code: Runtime Polymorphism. https://youtu.be/QGcVXgEVMJg<br/>[14] Chris Cleeland, Douglas C. Schmidt - External Polymorphism. https://www.dre.vanderbilt.edu/~schmidt/PDF/C++-EP.pdf<br/>[15] 深入浅出C++类型擦除（1） - 知乎. https://zhuanlan.zhihu.com/p/351291649<br/>[16] Philip Wadler, Stephen Blott - How to Make Ad-hoc Polymorphism Less Ad-hoc. https://dl.acm.org/doi/pdf/10.1145/75277.75283<br/>[17] Go generic https://github.com/golang/proposal/blob/master/design/generics-implementation-dictionaries-go1.18.md<br/>[18] Walter E. Brown - Modern Template Metaprogramming: A Compendium, Part I. https://youtu.be/Am2is2QCvxY<br/>

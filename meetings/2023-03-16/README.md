# gcc 符号版本管理

Speakers: Weiliang Qu

<https://zhuanlan.zhihu.com/p/314912277>

为什么要有符号版本（symbol versioning），主要是为了兼容。当某个 so 更新后，其中的部分函数（符号）行为发生了变更时，很可能导致依赖该 so 的 binary 需要重新编译链接才能继续正常使用。通过 symbol versioning 技术，可以提供完整的向后兼容，binary 无需感知所依赖 so 的版本更新。

## version 相关的 section

在使用 symbol versioning 的 ELF 文件中，有至多三个相关的 sections：.gnu.version，.gnu.version_r，.gnu.version_d，其中 .gnu.version_r 和 .gnu.version_d 是可选的
gnu.version（version symbol section）。dynamic table 中的 DT_VERSYM tag 指向该 section
.gnu.version_r（version requirement section）。dynamic table 中的 DT_VERNEED/DT_VERNEEDNUM tags 标记该 section。
.gnu.version_d（version definition section）。dynamic table 中的 DT_VERDEF/DT_VERDEFNUM tags 标记该 section。

看一个具体的例子

```c
readelf -a /bin/ls
...
Symbol table '.dynsym' contains 129 entries:
   Num:    Value          Size Type    Bind   Vis      Ndx Name
     0: 0000000000000000     0 NOTYPE  LOCAL  DEFAULT  UND
     1: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND __ctype_toupper_loc@GLIBC_2.3 (2)
     2: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND __uflow@GLIBC_2.2.5 (3)
     3: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND getenv@GLIBC_2.2.5 (3)
     4: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND cap_to_text
...
    56: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND acl_extended_file@ACL_1.0 (7)
    57: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND closedir@GLIBC_2.2.5 (3)
...
   123: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND __sprintf_chk@GLIBC_2.3.4 (4)
   124: 000000000061b600     0 NOTYPE  GLOBAL DEFAULT   25 _edata
   125: 000000000061c320     0 NOTYPE  GLOBAL DEFAULT   26 _end
   126: 000000000061b600     0 NOTYPE  GLOBAL DEFAULT   26 __bss_start
   127: 0000000000402190     0 FUNC    GLOBAL DEFAULT   11 _init
   128: 0000000000412a3c     0 FUNC    GLOBAL DEFAULT   14 _fini

Version symbols section '.gnu.version' contains 129 entries:
 Addr: 000000000040145a  Offset: 0x00145a  Link: 5 (.dynsym)
  000:   0 (*local*)       2 (GLIBC_2.3)     3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  004:   0 (*local*)       3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  008:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   4 (GLIBC_2.3.4)   3 (GLIBC_2.2.5)
  00c:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  010:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  014:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  018:   3 (GLIBC_2.2.5)   5 (GLIBC_2.17)    3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  01c:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  020:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  024:   3 (GLIBC_2.2.5)   6 (GLIBC_2.4)     3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  028:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  02c:   3 (GLIBC_2.2.5)   0 (*local*)       3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  030:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  034:   7 (ACL_1.0)       3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  038:   7 (ACL_1.0)       3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  03c:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   0 (*local*)
  040:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  044:   3 (GLIBC_2.2.5)   4 (GLIBC_2.3.4)   3 (GLIBC_2.2.5)   0 (*local*)
  048:   8 (GLIBC_2.14)    3 (GLIBC_2.2.5)   0 (*local*)       3 (GLIBC_2.2.5)
  04c:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  050:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  054:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  058:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   7 (ACL_1.0)       3 (GLIBC_2.2.5)
  05c:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  060:   4 (GLIBC_2.3.4)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  064:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  068:   0 (*local*)       0 (*local*)       3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  06c:   3 (GLIBC_2.2.5)   0 (*local*)       3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  070:   3 (GLIBC_2.2.5)   4 (GLIBC_2.3.4)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  074:   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)   3 (GLIBC_2.2.5)
  078:   2 (GLIBC_2.3)     2 (GLIBC_2.3)     3 (GLIBC_2.2.5)   4 (GLIBC_2.3.4)
  07c:   1 (*global*)      1 (*global*)      1 (*global*)      1 (*global*)
  080:   1 (*global*)

Version needs section '.gnu.version_r' contains 2 entries:
 Addr: 0x0000000000401560  Offset: 0x001560  Link: 6 (.dynstr)
  000000: Version: 1  File: libacl.so.1  Cnt: 1
  0x0010:   Name: ACL_1.0  Flags: none  Version: 7
  0x0020: Version: 1  File: libc.so.6  Cnt: 6
  0x0030:   Name: GLIBC_2.14  Flags: none  Version: 8
  0x0040:   Name: GLIBC_2.4  Flags: none  Version: 6
  0x0050:   Name: GLIBC_2.17  Flags: none  Version: 5
  0x0060:   Name: GLIBC_2.3.4  Flags: none  Version: 4
  0x0070:   Name: GLIBC_2.2.5  Flags: none  Version: 3
  0x0080:   Name: GLIBC_2.3  Flags: none  Version: 2
```

.gnu.version 一共有 129 个表项，对应着 .dynsym 的表项（符号）。.gnu.version_r 一共 2 个表项，记录了未定义符号的外部依赖（依赖 libacl.so.1 和 libc.so.6），其中，libacl.so.1 只涉及一个符号版本，libc.so.6 则涉及 6 个符号版本。所有涉及的符号版本，都有一个编号，从 2 ~ 8。

version index:
version index 0 被称为 VER_NDX_LOCAL。version index 为 0 的符号的 binding 将会更改为 STB_LOCAL。 version index 1 被称为 VER_NDX_GLOBAL，没有特殊作用。index 2 到 0xffef 用于其他 versions。

.gnu.version_r 记录了 version Name 和 index 的对应关系，.gnu.version 记录的则是 version index
.dynsym 每一项的 Name，如果后面跟着 "@"，说明是一个 versioned symbol，比如 "getenv@GLIBC_2.2.5 (3)"，表示 getenv 这个未定义（UND）符号使用的是 GLIBC_2.2.5 这个版本，版本 GLIBC_2.2.5 在 /bin/ls 中对应的 version index = 3

## versioned symbols 定义

symbol versinoning 只适用于动态库（so），定义 versioned symbols 后，需要在链接时通过 --version-script 指定 version script（参考 https://sourceware.org/binutils/docs/ld/VERSION.html）。

定义 versioned symbols 有两种方式：
符号名不包括 "@"，被 version script 的规则匹配获得 version
符号名包括 "@"，由 symver 定义。但实际上，这种方式无法支持动态链接库的。当指定 -shared 选项时，ld 会提示 version node not found for symbol 错误（参考 https://maskray.me/blog/2020-11-26-all-about-symbol-versioning）

由于符号版本不是 C 语言的标准用法，所以 gcc 使用了汇编器的一个特殊指示，也就是 .symver 指示，配合内嵌汇编的方式定义符号的版本。
**asm**(".symver original_foo,foo@");  
**asm**(".symver old_foo,foo@VERS_1.1");  
**asm**(".symver old_foo1,foo@VERS_1.2");  
**asm**(".symver new_foo,foo@@VERS_2.0");

例子中定义了 foo 的四个版本，其中的 symver 的第一个参数为源代码中真正定义的实现，而之后的 foo 则为对外公开的符号。foo@@VERS_2.0 为缺省版本。

## 链接器符号解析

链接器符号解析规则：
定义的 foo 可以满足未定义的 foo（非多版本符号）
定义的 foo@v1 可以满足未定义的 foo@v1
定义的 foo@@v1 可以同时满足未定义的 foo 和 foo@v1

若存在多个 default version 的定义（如 foo@@v1，foo@@v2)，触发 duplicate definition error。通常一个符号有零或一个 default version（@@）定义，任意个 non-default version（@）定义。

## 指定符号版本

当编译环境的 libc 版本比运行时环境的 libc 版本要高时，指定符号版本可以避免运行时找不到符号的错误。以 memcpy 函数为例

```c
[qwl@localhost ~]$ readelf -s /lib64/libc-2.17.so | grep memcpy
   335: 00000000000a8d00     9 FUNC    WEAK   DEFAULT   12 wmemcpy@@GLIBC_2.2.5
  1012: 0000000000117190    20 FUNC    GLOBAL DEFAULT   12 __wmemcpy_chk@@GLIBC_2.4
  1130: 0000000000094bc0    85 IFUNC   GLOBAL DEFAULT   12 memcpy@@GLIBC_2.14
  1132: 000000000008f950    75 IFUNC   GLOBAL DEFAULT   12 memcpy@GLIBC_2.2.5
...
```

libc-2.17 提供了两个版本 GLIBC_2.2.5 和 GLIBC_2.14 的 memcpy 实现，其中缺省为 GLIBC_2.14 版本。也就是说，如果在代码中不指定 memcpy 的版本，编译链接后 binary 中 UND memcpy 指向的是 GLIBC_2.14 版本。

```c
Symbol table '.dynsym' contains 5 entries:
   Num:    Value          Size Type    Bind   Vis      Ndx Name
     0: 0000000000000000     0 NOTYPE  LOCAL  DEFAULT  UND
     1: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND puts@GLIBC_2.2.5 (2)
     2: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND __libc_start_main@GLIBC_2.2.5 (2)
     3: 0000000000000000     0 NOTYPE  WEAK   DEFAULT  UND __gmon_start__
     4: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND memcpy@GLIBC_2.14 (3)
```

如果运行环境的 libc 版本没有提供 memcpy@GLIBC_2.14，binary 运行时将会报错。所以我们使用 symver 指定了要使用的版本：

```c
// mem.c
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

int main()
{
        char buf[32];
        __asm__(".symver memcpy, memcpy@GLIBC_2.2.5");
        memcpy(buf, "testestesttest", sizeof(buf) - 1);
        printf("%s\n", buf);
        return 0;
}

[qwl@localhost c]$ gcc -g -o mem mem.c
[qwl@localhost c]$ readelf -V -s mem

Symbol table '.dynsym' contains 5 entries:
   Num:    Value          Size Type    Bind   Vis      Ndx Name
     0: 0000000000000000     0 NOTYPE  LOCAL  DEFAULT  UND
     1: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND memcpy@GLIBC_2.2.5 (2)
     2: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND puts@GLIBC_2.2.5 (2)
     3: 0000000000000000     0 FUNC    GLOBAL DEFAULT  UND __libc_start_main@GLIBC_2.2.5 (2)
     4: 0000000000000000     0 NOTYPE  WEAK   DEFAULT  UND __gmon_start__
...
```

如果指定的版本编译环境找不到，链接时就会提示错误

```c
[qwl@localhost c]$ cat mem.c
...
int main()
{
        char buf[32];
        __asm__(".symver memcpy, memcpy@GLIBC_2.18");
        memcpy(buf, "testestesttest", sizeof(buf) - 1);
        printf("%s\n", buf);
        return 0;
}

[qwl@localhost c]$ gcc -g -o mem mem.c
/tmp/ccNbQwAc.o: In function `main':
/home/qwl/src/c/mem.c:9: undefined reference to `memcpy@GLIBC_2.18'
collect2: error: ld returned 1 exit status
```

glibc 编译时创建的 libdl.map（用于 version-script），生成 libdl.so 时使用

gcc -shared -static-libgcc -Wl,-O1 -Wl,-z,defs -Wl,-dynamic-linker=/xxx/glibc/install-glibc/lib/ld-linux-x86-64.so.2 -B/xxx/glibc/build-glibc/csu/ -Wl,--version-script=/xxx/glibc/build-glibc/libdl.map -Wl,-soname=libdl.so.2 -Wl,-z,combreloc -Wl,-z,relro -Wl,--hash-style=both -L/xxx/glibc/build-glibc -L/xxx/glibc/build-glibc/math -L/xxx/glibc/build-glibc/elf -L/xxx/glibc/build-glibc/dlfcn -L/xxx/glibc/build-glibc/nss -L/xxx/glibc/build-glibc/nis -L/xxx/glibc/build-glibc/rt -L/xxx/glibc/build-glibc/resolv -L/xxx/glibc/build-glibc/crypt -L/xxx/glibc/build-glibc/mathvec -L/xxx/glibc/build-glibc/support -L/xxx/glibc/build-glibc/nptl -Wl,-rpath-link=/xxx/glibc/build-glibc:/xxx/glibc/build-glibc/math:/xxx/glibc/build-glibc/elf:/xxx/glibc/build-glibc/dlfcn:/xxx/glibc/build-glibc/nss:/xxx/glibc/build-glibc/nis:/xxx/glibc/build-glibc/rt:/xxx/glibc/build-glibc/resolv:/xxx/glibc/build-glibc/crypt:/xxx/glibc/build-glibc/mathvec:/xxx/glibc/build-glibc/support:/xxx/glibc/build-glibc/nptl -o /xxx/glibc/build-glibc/dlfcn/libdl.so -T /xxx/glibc/build-glibc/shlib.lds /xxx/glibc/build-glibc/csu/abi-note.o -Wl,--whole-archive /xxx/glibc/build-glibc/dlfcn/libdl_pic.a -Wl,--no-whole-archive -Wl,--start-group /xxx/glibc/build-glibc/libc.so /xxx/glibc/build-glibc/libc_nonshared.a -Wl,--as-needed /xxx/glibc/build-glibc/elf/ld.so -Wl,--no-as-needed -Wl,--end-group

```c
GLIBC_2.2.5 {
  global:
    dladdr; dlclose; dlerror; dlopen; dlsym;
    dlopen; dlvsym;
  local:
    *;
};
GLIBC_2.3.3 {
  global:
    dladdr1; dlinfo;
} GLIBC_2.2.5;
GLIBC_2.3.4 {
  global:
    dlmopen;
} GLIBC_2.3.3;
GLIBC_PRIVATE {
  global:
    _dlfcn_hook;
} GLIBC_2.3.4;
```

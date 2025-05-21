# Overview

Bitly is a simple expression evaluation REPL, supporting (**in order of precedence**):

**Prefix Operators**
- Unary Negate (`-4`)

**Infix Operators**
- Bitwise And (`3 & 1`) and Or (`2 | 1`)
- Bit Shift Left (`1 << 3`) and Right (`8 >> 3`)
- Exponentiation (`2 ^ 2`)
- Modulation (`6 % 3`)
- Multiplication (`1 * 2`)
- Division (`4 / 2`)
- Addition (`2 + 2`)
- Subtraction (`2 - 2`)

Bitly also supports variables:
```go
> var x = 3 + 8
[x]=>[11]
> x + 4
Result[base-10] => 15
```

Binary (`0b`-prefix) and hex (`0x`-prefix) notations can be used interchangeably:
```go
> var x = 0b1 << 0xa + 4
[x]=>[1028]
```

# Description

I like to utilize minimally-sized integers to encode configuration/options and state where possible. During these experiments I often find myself wanting to quickly check my mental math for things like bit shifts. This want usually translates to heading to Google and using multiple tools (because for whatever reason a single tool doesn't seem to conveniently support all of these operations) to get the job done. This workflow sucks, so.. Bitly!

# CLI Usage

`// TODO`

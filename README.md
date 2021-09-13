# TLPS (Toy Language Processing System)

[![Go](https://github.com/goropikari/tlps/actions/workflows/go.yml/badge.svg)](https://github.com/goropikari/tlps/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goropikari/tlps)](https://goreportcard.com/report/github.com/goropikari/tlps)

TLPS is a toy language processing system for me to study.
TLPS adopts [Off-side rule](https://en.wikipedia.org/wiki/Off-side_rule) like python.



```
git clone https://github.com/goropikari/tlps
cd tlps
docker build -t tlps .
docker run -it tlps # launch REPL
```

Architecture is based on jlox (tree-walk interpreter) by [munificent/craftinginterpreters](https://github.com/munificent/craftinginterpreters).


```
// declare variable
var x = 10; // terminal ';' is optional

var こんにちは = "Hello World"
print(こんにちは) // => Hello World

// if statement
if expr:
  statements
elseif expr:
  statements
else:
  statements

// while loop
while expr:
  statements

// for loop
for var i = 0; i < 5; i = i + 1:
  print(i)

// function
fun fib(n):
  if n <= 1:
    return n
  return fib(n-1) + fib(n-2)

// closure
fun makeCounter():
  var i = 0
  fun count():
    i = i + 1
    return i
  return count

var counter = makeCounter()
print(counter()) // => 1
print(counter()) // => 2

// class
// there is no class variable
class Hoge:
  pass

class Hoge:
  hoge(x, y):
    pass

class Hoge:
  init(x, y):
    this.x = x
    this.y = y

class Hoge:
  piyo(x, y):
    print(x + y)

var h = Hoge()
h.name = "hoge piyo" // instance variable can define anytime
print(h.name) // => hoge piyo

// inheritance
class A:
  method():
    return "a"

class B(A):
  method():
    return "b"
  methodA():
    return super.method()

class C(B):
  pass

print(C().methodA()) // => a

// include another file
include "another.tlps" // path is relative path from the file which describe include statement

// indentation can be used for if branch, loop body and function body
var x = 1
  var y = 1  // indentation error
```

# Todo

- [x] escape sequence
- [x] detect IndentationError
- [x] import another file
  - [ ] detect circular import
- [ ] support varargs
- [ ] support IO

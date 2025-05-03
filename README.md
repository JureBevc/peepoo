# 💩 Peepee poopoo Language

The most sophisticated, vowel-based, programming language ever conceived.

> A language consisting strictly of the letter p separated by vowels. Yes, even variable names. You're welcome.
## Content

- [What Is This?](#🚽-what-is-this)
- [Quick Guide](#quick-guide)
    - [Syntax Rules](#📝-syntax-rules)
        - [Assigning Variables](#📦-assigning-variables)
        - [Operators](#➕-operators)
        - [Printing](#📤-printing)
        - [Conditionals](#❗-conditionals)
        - [Loops](#🔁-loops)
        - [Lists](#📋-lists)
        - [Functions](#🧙‍♂️-functions)
        - [Built-in functions](#🧝‍♂️-built-in-functions)
    - [Short Examples](#🧠-short-examples)
    - [Longer Examples](#🧠🧠-longer-examples)
    - [Run and Build](#🏃-run-and-build)



## 🚽 What Is This?

Peepee poopoo Language is an interpreted language made entirely from the letter `p` and a few brave vowels. Everything — variables, values, keywords — follows the sacred `p + vowel + p + vowel...` pattern.


# Quick Guide

## 📝 Syntax Rules

- Only **p** and vowels (`a, e, i, o, u`) are allowed.
- **Variable names**: UPPERCASE, must follow peepee poopoo pattern.
- **Commands**: lowercase, same pattern.
- **Numbers**: Written in peepee poopoo binary.  
  - Example:  
    - `po` → `0`  
    - `pi` → `1`  
    - `pipo` → `10` → decimal `2`  
    - `pipopo` → `100` → decimal `4`
- **Characters**: Written in (base 5) peepee poopoo style.
    - Example:
        - `a` → `papopupi`
        - `b` → `papopupo`
        - See `-encode` and `-decode` options in the [Run and Build](#🏃-run-and-build) section.

### 📦 Assigning Variables

```
PEE pe pipo
```

`PEE` is now 2.


### ➕ Operators

Math is expressed with unique operator keywords:

| Operation | Keyword   |
|-----------|-----------|
| +         | `pu`      |
| -         | `puu`     |
| *         | `pupu`    |
| /         | `puupuu`  |

Example:
```
PEE pe pi pu pipo
```

Stores `1 + 2` into `PEE`.


### 📤 Printing

Use `paa` to print without newline, `paapa` to print with newline.

```
paapa PEE
```

### ❗ Conditionals

Use `pii` to start an `if` block and `piipii` to close it. Block runs only if condition ≠ 0.

```
pii pipo
    paapa PEE
piipii
```

Prints `2` because `pipo` is 2.

### 🔁 Loops

Use `pepo` to start a loop and `pope` to end it. The loop variable auto-increments from 0 to the upper bound (exclusive).

```
pepo PEE po pipopo
    paapa PEE
pope
```

Prints `0` to `3`.

### 📋 Lists

Define a list of values by listing the values between two `pepe` keywords. This example shows a list `[0, 1, 2]` being defined and stored into the variable `PA`.

```
PA pe pepe po pi pipo pepe
```

Use `pepepa` to append another value to the list.

```
PA pepepa pipi 
```

Use `pepepi` to read or write a value to the list. 

```
PAPI pe PA pepepi po
PA pepepi po pe pi
```

Use `pepepo` to pop a value from a list. This removes the value at a given index and returns the removed value.

```
PAPI pe PA pepepo po
```

Use `pepepe` to get the length of a list.

```
PAPI pe pepepe PA
```




### 🧙‍♂️ Functions

Define a function with `poo`, end with `poopoo`. First word is the function name, everything that follows is a parameter.

```
poo PAPOPE PA poo
    paapa PA
    peepee PA
poopoo
```
This will define a function called `PAPOPE` that accepts a variable `PA`, which will get printed with `paapa` and returned with `peepee`.

A function is called like this:
```
pee PAPOPE pi pee
```
This wil call the function `PAPOPE` with parameter `pi` which is equal to `1`.

### 🧝‍♂️ Built-in Functions

Some functions are already predefined. Note that user-defined functions and variables are always written in uppercase letters, whereas built-in functions use a mix of uppercase and lowercase letters.

Built-in functions are not called in the same way as user-defined functions, because they are not wrapped between two `pee` keywords.

| Function description  | Name      | Example usage |
|-----------------------|-----------|------------------------------------|
| Read from stdin       | `PIpi`    | `PI pe PIpi`                       |
| Read from file        | `PIPIpi`  | `PI pe PIPIpi PA`                  |
| Character to int      | `POpi`    | `PI pe POpi papipupi`              |


## 🧠 Short Examples

### Add Two Numbers and Print

```
PA pe pi
    PE pe pipo
    PI pe PA pu PE
paapa PI
```
Outputs: `3`
### Loop and Print Squares

```
pepo PI po pipopo
    PA pe PI pupu PI
    paa PA
pope
```

Outputs: `0149`

## 🧠🧠 Longer Examples

1. [Hello world example](./examples/helloworld.peepoo)
2. [Fibonacci example](./examples/fibonacci.peepoo)
3. [Palindrome example](./examples/palindrome.peepoo)
4. [Sorting example](./examples/sort.peepoo)

## 🏃 Run and Build

Run without compiling:
```
go run . input.peepoo
```

Compile and run program:
```
go build -o peepoo main.go
./peepoo input.peepoo
```

Encode string:
```
./peepoo -encode "Hello world!"

papipupi papupape papupepo papupepo papupipe papepepi papupopu papupipe papupipu papupepo papupapa papepepo
```

Decode string:
```
./peepoo -encode "papipupi papupape papupepo papupepo papupipe papepepi papupopu papupipe papupipu papupepo papupapa papepepo"

Hello world!
```
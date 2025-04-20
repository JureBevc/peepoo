# 💩 Peepee peepoo Language

The most sophisticated, vowel-based, programming language ever conceived.

> A language consisting strictly of the letter p separated by vowels. Yes, even variable names. You're welcome.


## 🚽 What Is This?

Peepee peepoo Language is an interpreted language made entirely from the letter `p` and a few brave vowels. Everything — variables, values, keywords — follows the sacred `p + vowel + p + vowel...` pattern.


## 📝 Syntax Rules

- Only **p** and vowels (`a, e, i, o, u`) are allowed.
- **Variable names**: UPPERCASE, must follow peepee peepoo pattern.
- **Commands**: lowercase, same pattern.
- **Numbers**: Written in peepee peepoo binary.  
  - Example:  
    - `po` → `0`  
    - `pi` → `1`  
    - `pipo` → `10` → decimal `2`  
    - `pipopo` → `100` → decimal `4`


## 📦 Assigning Variables

```
PEE pe pipo
```

`PEE` is now 2.


## ➕ Operators

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


## 📤 Printing

Use `paa` to print without newline, `paapa` to print with newline.

```
paapa PEE
```

## ❗ Conditionals

Use `pii` to start an `if` block and `piipii` to close it. Block runs only if condition ≠ 0.

```
pii pipo
    paapa PEE
piipii
```

Prints `2` because `pipo` is 2.

## 🔁 Loops

Use `pepo` to start a loop and `pope` to end it. The loop variable auto-increments from 0 to the upper bound (exclusive).

```
pepo PEE po pipopo
    paapa PEE
pope
```

Prints `0` to `3`.


## 🧙‍♂️ Functions

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

[Fibonacci example](./examples/fibonacci.poopoo)

## 🏃 Run/Build from source

Run without compiling:
```
go run . input.peepoo
```

Compile and run:
```
go build -o peepoo main.go
./peepoo input.peepoo
```
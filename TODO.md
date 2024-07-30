# Checklist

## Features
- [ ] Lexer (30%)
- [ ] Parser (30%)
- [ ] Optimisation Pipeline (0%)
- [ ] Interpreter (15%)
- [x] REPL (100%)

## Language Support

### Variable Declarations
- [x] `var` support
- [x] `let` support
- [x] `const` support

### Function Declarations and Calls
- [ ] Named functions
- [ ] Anonymous functions (function expressions)
- [ ] Arrow functions

### Control Flow
- [x] `if` statements
- [x] `else` statements
- [x] `else if` statements
- [x] `switch` statements
  - [ ] without block scope
- [ ] `for` loops
- [ ] `for...of` loops
- [ ] `for...in` loops
- [x] `while` loops
- [ ] `do...while` loops

### Operators
- [x] Arithmetic operators (`+`, `-`, `*`, `/`, `%`)
- [x] Assignment operators (`=`, `+=`, `-=`, `*=`, `/=`)
- [x] Comparison operators (`==`, `===`, `!=`, `!==`, `<`, `>`, `<=`, `>=`)
- [x] Logical operators (`&&`, `||`, `!`)
- [x] Bitwise operators (`&`, `|`, `^`, `<<`, `>>`, `>>>`)
- [x] Unary operators (`++`, `--`, `typeof`, `delete`)
- [x] Ternary operator (`?:`)

### Objects and Arrays
- [ ] Object creation
- [ ] Object property access
- [x] Array creation
- [x] Array indexing
- [ ] Array length property

### Built-in Functions and Objects
- [x] `console.log`
- [x] `Math` object (`Math.sqrt`, `Math.pow`, etc.)

### String Manipulation
- [x] String literals
- [ ] String methods (`length`, `substring`, `toUpperCase`, `toLowerCase`, etc.)

### Error Handling
- [ ] `try` statements
- [ ] `catch` statements
- [ ] `finally` statements
- [ ] Throwing errors

### Scope and Closures
- [ ] Lexical scoping
- [ ] Closure support

### ES6 Features
- [ ] Template literals
- [ ] Destructuring assignment
- [ ] Spread/rest operators

### Type Coercion and Conversion
- [ ] Implicit type conversions
- [ ] Explicit type conversions

### Promises and Asynchronous Programming
- [ ] Basic support for `Promise` objects

### Modules (optional)
- [ ] Basic support for `import`
- [ ] Basic support for `export`

## Notes

- Make an optimisation pipeline and include common techs:
  - Constant Folding
  - Dead Code Elimination
  - Inlining
  - Common Subexpression Elimination, Strength Reduction, ...
  - ...
- Array access and member access in inifix operations are broken. Need to fix that.
- ~~Pointers to token type in lexer instead of copies for performance reasons~~
- ~~Reassignments of variables~~
- ~~Fixing infix operators generally lol~~
- `===` triple operators seem to be broken
- Proper function declarations and calls (no arrow functions)
- String concatenation and interpolation (`Hello ${name}`)
- What happens if you try to use operators on the wrong types?
- Making `const` actually constant. Currently it's just a keyword that doesn't do anything.
- ~~`else if` statements~~
- `for` loops
  - for loops
  - for...in loops
  - for...of loops
- Basic array support?
  - Accessing array elements (arr[index])
  - Array methods (push, pop, shift, unshift)
- Basic json support?
  - Accessing object properties (obj.prop and obj["prop"])
  - Adding properties to objects
- Instead storing line number and text in token structs just store positions in input and look them as needed?
- Some sort of block scope. Perhaps change `env` to a stack of environments (scopes)?
  - Block scope (for let and const)
  - Function scope
  - Global scope
- `...` spread and rest operators
    - Spread operator for arrays and objects (...)
    - Rest parameters in function definitions

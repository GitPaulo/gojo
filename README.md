# gojo

Go javascript engine-o.

```sh
go run main.go
```

Use environment variables to configure the engine.

- `GOJO_INPUT_FILE` - Set file path to use it as input (`default=input_program.js`).
- `GOJO_VERBOSE` - Set to `true` to enable verbose logging.
- `GOJO_MEGA_VERBOSE` - Set to `true` to enable EVEN MORE logging.
- `GOJO_REPL_MODE` - Set to `true` to enable REPL mode.

### Tests

To run lexer, parser and interpreter tests:

```sh
go test ./...
```

## Checkpoints

> Yo! This is a for fun project and not intended to ever be finished. o7

The goal is to have a resonably fast and correct javascript engine for the core features of the language only.

Checkout the [TODO.md](TODO.md).

## Parts by example:

For the code `var x = 1 + 1;`

1. Lexer (FSM): Produces a stream of tokens `var, x, =, 1, +, 1, ;`
2. Parser (Recursive Descent): Produces an AST `Program -> Statements[] -> VariableDeclaration(x = 
   BinaryExpression(1 + 1))`
3. Interpreter (Go Runtime): Process the AST and stores the value in the environment `x = 2`

![img.png](img.png)

## References

- ECMAScript spec - https://tc39.es/ecma262/
- JSConf rust JS engine talk - https://youtu.be/_uD2pijcSi4
- Writing a tokenizer (not great performance wise but a start) - https://dev.to/ndesmic/writing-a-tokenizer-1j85
- General simple lexer - https://codemaster138.github.io/blog/creating-an-interpreter-part-1-lexer/
- Writing an interpreter in go - https://edu.anarcho-copy.org/Programming%20Languages/Go/writing%20an%20INTERPRETER%20in%20go.pdf
- Acorn js js parser - https://github.com/acornjs/acorn
- Standard compliant parser - https://github.com/jquery/esprima
- AST Explore - https://astexplorer.net/
- Recursive Descent Parsing: https://www.cs.rochester.edu/u/nelson/courses/csc_173/grammars/parsing.html#:~:text=Recursive-descent%20parsing%20is%20one,non-terminal%20with%20a%20procedure
- Goja: https://github.com/dop251/goja (great minds think alike)


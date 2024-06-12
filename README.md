# gojo

Go javascript engine.

```
go run main.go
```

Use environment variables to configure the engine.
- `GOJO_VERBOSE` - Set to `true` to enable verbose logging.
- `GOJO_TEST` - Set to `true` to run tests only.
- `GOJO_INPUT_FILE` - Set file path to use it as input.

## Checkpoints

> Yo! This is for learning purposes and not intended to ever be finished.

- Lexer: About halfway to compliance with the ECMAScript spec.
- Parser: Basic parser that can parse control statements variables functions and binary expressions.
- Interpreter (& Go Runtime): Very basic linking between the parser and the interpreter. Only supports variable declaration and binary operators.

## Parts by example:

For the code `var x = 1 + 1;`

1. Lexer: Produces a stream of tokens `var, x, =, 1, +, 1, ;`
2. Parser: Produces an AST `Program -> Statements[] -> VariableDeclaration(x = BinaryExpression(1 + 1))`
3. Interpreter: Process the AST and stores the value in the environment `x = 2`

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


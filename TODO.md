## TODO

- ~~Pointers to token type in lexer instead of copies for performance reasons~~
- ~~Reassignments of variables~~
- Instead storing line number and text in token structs just store positions in input and look them as needed?
- Some sort of block scope. Perhaps change `env` to a stack of environments (scopes)?
- Fixing infix operators generally lol

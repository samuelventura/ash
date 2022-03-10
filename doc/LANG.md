# ash

Very focused and opinionated automation scripting language.

## Goals

- Object style API
- Not general purpose
- High level abstractions
- Backend/frontend glue
- WASM/JS/GO targets
- Text data flow
- Self healing 
   
## Strategy

- Native as much as possible (WASM/JS/GO)
- Have a native API that feels scriptable
- Add a glueing script on top of it

## Principles

- Check and defer
- Least surprise
- Let is crash
- Bubble error

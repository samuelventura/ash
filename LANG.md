# ash

Very focused and opinionated automation scripting language.

## Architecture

- BE/FE interaction
  - `Tree` BE->FE
  - `Event` FE->BE
  - `Queue` BE->FE
- BE abstractions
  - `Queue` FIFO/PubSub
  - `Tree` Exchange tree
  - `Loop` Execution loop
  - `View` User interface

## Goals

- Object style API
- Dynamically typed
- Not general purpose
- High level abstractions
- Backend/frontend glue
- WASM/JS/GO targets
- Text data flow
- Self healing 

## Strategy

- Do what is natural then refine
- Native as much as possible (WASM/JS/GO)
- Have a native API that feels scriptable
- Add a glueing script on top of it

## Principles

- Check and defer
- Least surprise
- Let is crash
- Bubble error

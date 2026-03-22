---
one: value 1
two: value 2
top:
    nested: nested value
list:
    - item 1
    - item 2
---

One: {{.one}}
Two: {{.two}}
Top Nested: {{.top.nested}}
List item 1: {{index .list 0}}
List item 2: {{index .list 1}}

- [X] Step 1
- [X] Step 2
  - [X] Sub Step 1
  - [ ] Sub Step 2
- [X] Step 3 value 1
- [ ] Step 3
- [ ] Step 4
---
one:
two:
top:
    nested:
list: []
---

One: {{.one}}
Two: {{.two}}
Top Nested: {{.top.nested}}
List item 1: {{index .list 0}}
List item 2: {{index .list 1}}

- [ ] Step 1
- [ ] Step 2
  - [ ] Sub Step 1
  - [ ] Sub Step 2
- [ ] Step 3 {{.one}}
- [ ] Step 3
- [ ] Step 4
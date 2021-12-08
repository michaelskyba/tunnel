# tunnel
Tunnel is a simple SM-2 implementation. It's a Go rewrite of
[scripture](https://github.com/michaelskyba/scripture), a similar project I
wrote in Shell. My goal is for tunnel to be faster than scripture and to follow 
the UNIX philosophy better. Being written in Go should bring other benefits,
such as a cleaner, more maintainable codebase and an automated test system
(without an added dependency).

## Status
Currently, tunnel is unusable. Use scripture for the time being and check
back for updates soon.

## Usage
### shovel
Shovel is a provided wrapper around tunnel that simplifies the most common usage.
The only syntax is ``shovel <deck_file>``, which launches an interactive review session.
Shovel is not designed to follow the UNIX philosophy; it's just an example of
a practical application of tunnel commands.

### Documentation of individual tunnel commands

#### format_deck
``format_deck`` will format a deck file for you. More specifically, it will create
required fields for cards and renumber their IDs to stay in order. ``format_deck``
also checks for invalid lines, printing an error if it finds one.

``format_deck`` can either take a deck filename as an argument or the deck 
contents as standard input. It will print to standard output.

```sh
~ $ cat deck
a	b
c	d
~ $ tunnel format_deck deck
1	a	b	2.5	0	2021-04-01
2	c	d	2.5	0	2021-04-01
~ $ tunnel format_deck < deck
1	a	b	2.5	0	2021-04-01
2	c	d	2.5	0	2021-04-01
~ $ echo e >> deck
~ $ tunnel format_deck deck
Invalid line at index 3: e
~ $ echo $? # Prints the exit code
1
```

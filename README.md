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
Shovel is a provided wrapper around tunnel that simplifies the most common
usage. The only syntax is ``shovel <deck_file>``, which launches an interactive
review session. Shovel is not designed to follow the UNIX philosophy; it's
just an example of a practical application of tunnel commands. If you choose to
use shovel, you don't have to worry about anything else, at least until you decide
that shovel is too limited to fit your needs.

### Tutorial
The basic tunnel "process" looks like this:
1. Install
```sh
~ $ git clone https://github.com/michaelskyba/tunnel
~ $ cd tunnel
~/tunnel $ go build tunnel.go
~/tunnel $ su -c "cp tunnel /usr/local/bin/"
```

2. Create a deck file
Read the [deck file](#deck_file) section to understand the syntax. ``example_deck``,
an example deck file, is provided, which I will use here.

3. Use ``new_cards`` to format new cards
```sh
~/tunnel $ tunnel new_cards deck
```

4. Use ``due`` to determine which cards need to be reviewed
```sh
~/tunnel $ tunnel due deck
1
3
4
```

5. TODO

### Documentation of individual tunnel commands

#### new_cards
```sh
~ $ cat deck
a	b
c	d
~ $ tunnel new_cards < deck
a	b	0	2.5	0	2021-04-01
c	d	0	2.5	0	2021-04-01
~ $ cat deck
a	b
c	d
~ $ tunnel new_cards deck
~ $ cat deck
a	b	0	2.5	0	2021-04-01
c	d	0	2.5	0	2021-04-01
```

``new_cards`` will add default SM-2 values to new card lines. Specifically, it
appends ``0	2.5	0	2021-04-01``. ``new_cards`` can either take a 
deck filename as an argument or the deck contents as stdin. If you use an 
argument, tunnel will update the specified file. If you use stdin, the 
formatted version will instead be sent to stdout.

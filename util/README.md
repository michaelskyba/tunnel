# util
These are other scripts or files that may be useful to a tunnel user.

## ``shovel``
``shovel`` is explained in the [main README](https://github.com/michaelskyba/tunnel/blob/main/README.md).

## ``descript``
This can be used for converting legacy [scripture](https://github.com/michaelskyba/scripture) to tunnel decks.

## ``parse_cloze``
``parse_cloze`` is used to parse cloze input in tunnel deck files. Wrap anything
you want to be separated into a card in square brackets. Then, pass it to
parse_cloze as stdin and capture the stdout.

```
$ cat input
The [terminal arm] of an angle in [standard position] is placed [in any direction, depending on the measurement of the angle]
$ parse_cloze < input
The [] of an angle in [standard position] is placed [in any direction, depending on the measurement of the angle]	terminal arm
The [terminal arm] of an angle in [] is placed [in any direction, depending on the measurement of the angle]	standard position
The [terminal arm] of an angle in [standard position] is placed []	in any direction, depending on the measurement of the angle
```

## ``decloze``
``decloze`` takes one or more cards parsed using ``parse_cloze`` and reverts
them back to their original cloze input. (I couldn't think of a better name.)

```
$ cat input
The [] of an angle in [standard position] is placed [in any direction, depending on the measurement of the angle]	terminal arm
The [terminal arm] of an angle in [] is placed [in any direction, depending on the measurement of the angle]	standard position
The [terminal arm] of an angle in [standard position] is placed []	in any direction, depending on the measurement of the angle
gcc stands for []	GNU Compiler Collection
$ decloze < input
The [terminal arm] of an angle in [standard position] is placed [in any direction, depending on the measurement of the angle]
gcc stands for [GNU Compiler Collection]
```

## ``tunnel.kak``
A sample configuration for Kakoune is provided in ``tunnel.kak``, which is
similar to what I use.

Here is a screenshot of the higlighting, with the terminal colours set by pywal:

![Kakoune screenshot](https://raw.githubusercontent.com/michaelskyba/tunnel/main/util/kak-screenshot.png)

## ``spread``
```sh
~ $ cat deck
a	b	0	2.5	0	1688434760
c	d	0	2.5	0	1688434760
e	f	0	2.5	0	1688434760
g	h	0	2.5	0	1688434760
~ $ date
1688434761
~ $ spread 2 deck
~ $ cat deck
a	b	0	2.5	0	1688434760
c	d	0	2.5	0	1688564360
e	f	0	2.5	0	1688521160
g	h	0	2.5	0	1688477960
~ $ spread 365 20 deck
~ $ cat deck
a	b	0	2.5	0	1719970760
c	d	0	2.5	0	1720402760
e	f	0	2.5	0	1721266760
g	h	0	2.5	0	1720834760
```

``spread``'s syntax is ``spread [offset] <days> <deck filename>``. The initial
introduction of cards will be spread out across ``<days>`` days. If an offset is
provided, this will be from ``now + <offset>`` to ``now + <offset> + <days>``.
Otherwise, it will be from now to ``now + <days>``. ``spread`` only applies to
cards formatted with ``tunnel new_cards`` which have not been reviewed yet.

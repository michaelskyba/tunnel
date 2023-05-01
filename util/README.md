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

## ``tunnel.kak``
A sample configuration for Kakoune is provided in ``tunnel.kak``, which is
similar to what I use.

Here is a screenshot of the higlighting, with the terminal colours set by pywal:

![Kakoune screenshot](https://raw.githubusercontent.com/michaelskyba/tunnel/main/util/kak-screenshot.png)

## ``spread``
```sh
~ $ cat deck
a	b	0	2.5	0	1682869324
c	d	0	2.5	0	1682869324
e	f	0	2.5	0	1682869324
g	h	0	2.5	0	1682869324
~ $ date
1682869325
~ $ spread 2 deck
~ $ cat deck
a	b	0	2.5	0	1682955725
c	d	0	2.5	0	1682869324
e	f	0	2.5	0	1682955725
g	h	0	2.5	0	1682869324
```

``spread``'s syntax is ``spread <days> <deck filename>``. It will modify the
last review placeholder dates and randomly spread them out over <days> days. As
shown, if, for example, you have four new cards and want to spread them out over
two days, two of the cards will be due today and two will be due tomorrow.

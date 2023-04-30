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

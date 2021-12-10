# tunnel
Tunnel is a simple SM-2 implementation. It's a Go rewrite/remake of
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

#### 1. Create a deck file.
Read the [deck file](#deck-file) section to understand the syntax. ``example_deck``,
an example deck file, is provided, which I will use here.

#### 2. Use ``new_cards`` to format new cards.
```sh
~/tunnel $ tunnel new_cards example_deck
```

#### 3. Use ``due`` to determine which cards need to be reviewed.
```sh
~/tunnel $ tunnel due example_deck
1
3
4
```

#### 4. Use ``review`` to review each due card.
Use ``front`` and ``back`` to see the fronts and backs of the respective cards.
```sh
~/tunnel $ tunnel front 1 example_deck
Symbol: Gold
~/tunnel $ # The user can't remember
~/tunnel $ tunnel back 1 example_deck
Au
~/tunnel $ # Now the answer seems obvious and familiar ("Oh, right, it's Au!")
~/tunnel $ tunnel review 1 2 example_deck
```
```sh
~/tunnel $ tunnel front 3 example_deck
Symbol: Silver
~/tunnel $ # The user thinks for a second and then remembers
~/tunnel $ tunnel back 3 example_deck
Ag
~/tunnel $ tunnel review 3 4 example_deck
```
```sh
~/tunnel $ tunnel front 4 example_deck
Symbol: Carbon
~/tunnel $ # The user remembered instantly
~/tunnel $ tunnel back 4 example_deck
C
~/tunnel $ tunnel review 4 5 example_deck
```
Repeat this process with ``due`` after each set of reviews until there are no
more due cards. You'll have to review a card multiple times if you score below
4, so you can't just use ``due`` once.
```sh
~/tunnel $ tunnel due example_deck
1
```
```sh
~/tunnel $ tunnel front 1 example_deck
Symbol: Gold
~/tunnel $ # The user remembered instantly
~/tunnel $ tunnel back 1 example_deck
Au
~/tunnel $ tunnel review 1 5 example_deck
~/tunnel $ tunnel due example_deck
~/tunnel $ # Reviews are done for this deck!
```
That's it! Understand, though, that these tunnel commands are never supposed
to be run manually like I showed here. This "tutorial" section is meant to
give you a general understanding of the order, so that making your own wrapper
(like shovel, but more fitted to your needs) will have less friction. Read the
documentation for each of the commands involved to learn more.

### Deck file
A deck file is a file containing a deck of cards, each of which will be reviewed.
Deck files are in TSV format. If you prefer e.g. CSV, convert your commas to
tabs before running tunnel and then conver them back to commas afterwards.

Cards are inputted in the syntax ``front	back``. In reviews, you will
look at the front of the card and attempt to recall the back of the card. This
is the only type of card, unlike e.g. Anki, which has many card types. You can
emulate other card types easily, though. For instance, I have a macro in my text
editor that converts
```
[Cellular respiration] happens in [every cell's mitochondrion]
```
to
```
[] happens in [every cell's mitochondrion]	Cellular respiration
[Cellular respiration] happens in []	every cell's mitochondrion
```
, thereby creating a sort of cloze-deletion card type.

If you want a line to act as a comment, don't put any tabs in it. Blank lines
are fine too.

After running ``new_cards``, you will see that ``0	2.5	0	
2021-04-01`` is added. The first 0 is the repetition number, the 2.5 is the
easiness factor, the second 0 is the inter-repetition interval, and the date
is the that of the last review. Since there haven't been any reviews for new
cards, the date is set to an arbitrary past date. As you review your cards,
these values will be updated. If you want to modify a card's front or backside,
feel free to edit the first two values (front and back), but never touch
any of the others (e.g. repetition number).

### Description of individual tunnel commands

#### ``new_cards``
```sh
~ $ cat deck
a	b
c	d
~ $ tunnel new_cards deck
~ $ cat deck
a	b	0	2.5	0	2021-04-01
c	d	0	2.5	0	2021-04-01
```

``new_cards``'s syntax is ``tunnel new_cards <deck filename>``. It will modify
the file as to add default SM-2 values to new card lines. Specifically, it
appends ``0	2.5	0	2021-04-01``.

#### ``due``
```sh
~ $ date
Thu Dec  9 03:04:51 PM EST 2021
~ $ cat letters
a	b	0	2.5	0	2021-04-01
c	d	0	2.5	0	2021-04-01
e	f	5	3	131.554	2021-11-21
h	i	1	2.46	1	2021-12-08
~ $ tunnel due letters
1
2
4
```

``due``'s syntax is ``tunnel due <deck filename>``. It will iterate over the
lines in the deck file and print those that are scheduled for review. This is
determined by looking at the inter-repetition interval and the last review
date. The numbers spit out are the positions of each card line, which you
should provide to ``front``, ``back``, and ``review``. If you add new cards
after ``due``ing but before reviewing, you should re-run ``due``, since the old
numbers may be inaccurate now.

#### ``front`` and ``back``
```sh
~ $ cat letters
a	b	1	2.46	1	2021-12-08
c	d
e
~ $ tunnel front 1 letters
a
~ $ tunnel back 2 letters
d
~ $ tunnel front 3 letters
Error: line 3 is not a valid card
~ $ tunnel back 4 letters
Error: no line 4 in deck
```
The syntax here is ``tunnel <front|back> <card line number> <deck filename>``.
``front`` will print the first tab-separated value in the card, which is the front
of the card, and ``back`` will print the second tab-separated value, which is the back.

#### ``review``
```sh
~ $ date
Thu Dec  9 05:45:13 PM EST 2021
~ $ cat letters
a	b	0	2.5	0	2021-04-01
c	d	5	3	131.554	2021-11-21
e
~ $ tunnel due letters
1
~ $ tunnel review 1 4 letters
~ $ cat letters
a	b	1	2.5	1	2021-12-09
c	d	5	3	131.554	2021-11-21
e
~ $ tunnel review 2 2 letters
Error: card 2 is not due for review.
~ $ tunnel review 3 5 letters
Error: line 3 is not a valid card
~ $ tunnel review 4 0 letters
Error: no line 4 in deck
```
``review``'s syntax is ``tunnel review <card line number> <review grade> <deck filename>``.
This will update the card's SM-2 fields in accordance to the SM-2 algorithm. To see
which cards you need to review, use ``due``.

The grades' meanings are as follows:
```
0: "Total blackout", complete failure to recall the information.
1: Incorrect response, but upon seeing the correct answer it felt familiar.
2: Incorrect response, but upon seeing the correct answer it seemed easy to remember.
3: Correct response, but required significant difficulty to recall.
4: Correct response, after some hesitation.
5: Correct response with perfect recall.
```
(https://en.wikipedia.org/wiki/SuperMemo#Description_of_SM-2_algorithm)

If you need to review a set of cards outside of their regular schedule, copy them into
a new, temporary deck and study that one.

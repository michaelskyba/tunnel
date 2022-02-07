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

## Installation
```sh
git clone https://github.com/michaelskyba/tunnel
cd tunnel
go build tunnel.go
su -c "cp tunnel /usr/local/bin/"
su -c "cp shovel /usr/local/bin/" # Only if you care about shovel
```

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
0
2
3
```

#### 4. Use ``review`` to review each due card.
Use ``front`` and ``back`` to see the fronts and backs of the respective cards.
```sh
~/tunnel $ tunnel front 0 example_deck
Symbol: Gold
~/tunnel $ # The user can't remember
~/tunnel $ tunnel back 0 example_deck
Au
~/tunnel $ # Now the answer seems obvious and familiar ("Oh, right, it's Au!")
~/tunnel $ tunnel review 0 2 example_deck
```
```sh
~/tunnel $ tunnel front 2 example_deck
Symbol: Silver
~/tunnel $ # The user thinks for a second and then remembers
~/tunnel $ tunnel back 2 example_deck
Ag
~/tunnel $ tunnel review 2 4 example_deck
```
```sh
~/tunnel $ tunnel front 3 example_deck
Symbol: Carbon
~/tunnel $ # The user remembered instantly
~/tunnel $ tunnel back 3 example_deck
C
~/tunnel $ tunnel review 3 5 example_deck
```

#### 5. Cycle through ``retry``
After going through the initial set of due cards, use ``retry`` to see which
cards need to be retried. Repeat this process with ``retry`` after each set of
reviews until there are no more cards to review. You'll have to review a card
again if you score below 4, so you can't skip ``retry`` or only use it once.
```sh
~/tunnel $ tunnel retry example_deck
0
```
```sh
~/tunnel $ tunnel front 0 example_deck
Symbol: Gold
~/tunnel $ # The user remembered instantly
~/tunnel $ tunnel back 0 example_deck
Au
~/tunnel $ tunnel review 0 5 example_deck
~/tunnel $ tunnel retry example_deck
~/tunnel $ # Reviews are done for this deck!
```
That's it! Understand, though, that these tunnel commands are never supposed
to be run manually as I showed here. This "tutorial" section is meant to
give you a general understanding of the order, so that making your own wrapper
(like shovel, but more fitted to your needs) will have less friction. Read the
documentation for each of the commands involved to learn more.

### Deck file
A deck file is a file containing a deck of cards, each of which will be reviewed.
Deck files are in TSV format. If you prefer e.g. CSV, convert your commas to
tabs before running tunnel and then convert them back to commas afterwards.

Cards are inputted in the syntax ``front<tab>back``. In reviews, you will
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

After running ``new_cards``, you will see that ``0x2.5x0x1617249600`` (where
``x``s are tabs) is added. The first 0 is the repetition number, the 2.5 is the
easiness factor, the second 0 is the inter-repetition interval, and the epoch
time is that of the last review. Since there haven't been any reviews for
new cards, the date is set to an arbitrary past date. As you review your cards,
these values will be updated. If you want to modify a card's front or backside,
feel free to edit the first two values (front and back), but never touch any of
the others (e.g. repetition number).

### Description of individual tunnel commands
Note that commands only check for validity in the context of their own
functions, so a card in your deck may be invalid even if e.g. ``front``
doesn't tell you that it is.

#### ``new_cards``
```sh
~ $ cat deck
a	b
c	d
~ $ tunnel new_cards deck
~ $ cat deck
a	b	0	2.5	0	1617249600
c	d	0	2.5	0	1617249600
```

``new_cards``'s syntax is ``tunnel new_cards <deck filename>``. It will modify
the file as to add default SM-2 values to new card lines. Specifically, it
appends ``0	2.5	0	1617249600``.

#### ``due``
```sh
~ $ date
Thu Dec  9 03:04:51 PM EST 2021
~ $ cat letters
a	b	0	2.5	0	1617249600
c	d	0	2.5	0	1617249600
e	f	5	3	131.554	1637470800
h	i	1	2.46	1	1638939600
~ $ tunnel due letters
0
1
3
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
a	b	1	2.46	1	1638939600
c	d
e
~ $ tunnel front 0 letters
a
~ $ tunnel back 1 letters
d
~ $ tunnel front 2 letters
Error: line 2 is not a valid card.
~ $ tunnel back 3 letters
Error: no line 3 in deck.
```
The syntax here is ``tunnel <front|back> <card line number> <deck filename>``.
``front`` will print the first tab-separated value in the card, which is the front
of the card, and ``back`` will print the second tab-separated value, which is the back.

#### ``review``
```sh
~ $ date
Thu Dec  9 05:45:13 PM EST 2021
~ $ cat letters
a	b	0	2.5	0	1617249600
c	d	5	3	131.554	1637470800
e
~ $ tunnel due letters
0
~ $ tunnel review 0 4 letters
~ $ cat letters
a	b	1	2.5	1	1639026000
c	d	5	3	131.554	1637470800
e
~ $ tunnel review 1 2 letters
Error: card 1 is not due for review.
~ $ tunnel review 2 5 letters
Error: line 2 is not a valid card.
~ $ tunnel review 3 0 letters
Error: no line 3 in deck.
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

#### ``retry``
```sh
~ $ tunnel retry geography
~ $ tunnel due geography
0
1
2
~ $ tunnel review 0 4 geography
~ $ tunnel review 1 3 geography
~ $ tunnel review 2 2 geography
~ $ tunnel retry geography
1
2
~ $ tunnel review 1 1 geography
~ $ tunnel review 2 4 geography
~ $ tunnel retry geography
1
~ $ tunnel review 1 5 geography
~ $ tunnel retry geography
~ $ # Done with reviews
```

``retry``'s syntax is ``tunnel retry <deck filename>``. After doing your initial review,
SM-2 wants you to retry any cards you gave a grade below 4. The ``retry`` command will
show you these "retry cards" so that you don't have to keep track of them yourself.
After each set of retries, the retry list will be updated. To keep track, tunnel uses
files in ``/tmp/tunnel``.

For the example above, the file would look like this after each command:
- ``tunnel review 1 3 geography``
```
1
```
- ``tunnel review 2 2 geography``
```
1
2
```
- ``tunnel review 1 1 geography``
(As you can see, a dash is used to indicate the next retry cycle)
```
2
-
1
```
- ``tunnel review 2 4 geography``
```
1
```
- ``tunnel review 1 5 geography``
(The file gets deleted)

If the deck file was ``/home/michael/decks/geography``, the retry file's path
would be ``/tmp/tunnel/home/michael/decks/geography``. We need to have the same
chain of directories because a user could have different deck files with the
same filename being reviewed.

Do not start moving lines around in your deck file after starting a review. If you
fail card 1 and thus the retry file contains card 1, there's no way tunnel will know
if you suddenly swap the first and second lines of your file. Then, retry will have
inaccurate information. So, if you want to make modifications, finish all reviews first.

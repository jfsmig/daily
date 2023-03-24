#!/usr/bin/env python3
#
# Generate a high-grade excuse to not attend the Daily Meeting.
#

#
# Copyright (C) 2023 Jean-Francois Smigielski
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
#

from random import choice, seed
from re import sub as substitute
from collections.abc import Iterable

def wrap(x):
    if isinstance(x, str):
        return Term(x)
    if isinstance(x, Iterable):
        return Choice(*x)
    return x

class Term(object):
    def __init__(self, v):
        self.value_ = str(v)
    def expand(self, out):
        out.write(self.value_)

class Choice(object):
    def __init__(self, *nargs):
        self.choices_ = tuple((wrap(x) for x in nargs))
    def expand(self, out):
        assert(len(self.choices_) > 0)
        choice(self.choices_).expand(out)

class Seq(object):
    def __init__(self, *nargs):
        self.items_ = tuple((wrap(x) for x in nargs))
    def expand(self, out):
        assert(len(self.items_)> 0)
        for c in self.items_:
            c.expand(out)
            out.write(' ')

class Writer(object):
    def __init__(self):
        self.tokens_ = list()
    def __str__(self):
        return ' '.join(self.tokens_)
    def write(self, x):
        self.tokens_.append(x)
    
seed(None)

statement = Choice("I cannot attend the daily meeting",
                   "Going to miss the meeting",
                   "Gonna miss the meeting",
                   "I won't be able to attend the daily meeting")

cause = Choice("Feeling sick this morning",
               "Woke up with a pretty bad headache",
               "Not feeling great",
               "I have a dog grooming errand",
               "My wife is stuck in her car",
               "I had a hard time last night",
               "I have to file for a new ID card",
               "I've lost my keys in the river",
               "someone stole my catalytic exhaust pipe",
               "My colonoscopy won't be done yet",
               "My prostate exam doesnt happen as expected",
               "My nan slipped on a dung")

conjonction_cause = Choice("because",
                           "'cos",
                           "since")

conjonction_consequence = Choice("thus",
                                 "then",
                                 "as a consequence",
                                 "therefore")

sentence = Choice(Seq(statement, conjonction_cause, cause),
                  Seq(cause, conjonction_consequence, statement))

out = Writer()
sentence.expand(out)
print(substitute(' +', ' ', str(out).strip()))


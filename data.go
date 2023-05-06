// Copyright (C) 2023 Jean-Francois Smigielski
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	_ "embed"
	"github.com/jfsmig/daily/excuse"
)

//go:embed data/ooo.txt
var oooData string

//go:embed data/meeting.txt
var meetingData string

func newOOO() (excuse.Generator, error) {
	statement := excuse.NewChoice(
		excuse.NewTerm("I'm going to be OOO,"),
		excuse.NewTerm("I need to be OOO today,"),
		excuse.NewTerm("I can't show up today,"))

	cause, err := excuse.ParseStreamString(oooData)
	if err != nil {
		return nil, err
	}

	return excuse.NewSequence(statement, cause), nil
}

func newNoMeeting() (excuse.Generator, error) {
	extenCause, err := excuse.ParseStreamString(meetingData)
	if err != nil {
		return nil, err
	}

	noDailyStatement := excuse.NewChoice(
		excuse.NewTerm("I cannot attend the daily, "),
		excuse.NewTerm("going to miss the meeting, "),
		excuse.NewTerm("gonna miss the meeting, "),
		excuse.NewTerm("No daily meeting for me, "))

	return excuse.NewSequence(noDailyStatement, extenCause), nil
}

func newGenerator() (excuse.Generator, error) {
	items := make([]excuse.Generator, 0)
	if n, err := newNoMeeting(); err != nil {
		return nil, err
	} else {
		items = append(items, n)
	}
	if n, err := newOOO(); err != nil {
		return nil, err
	} else {
		items = append(items, n)
	}
	return excuse.NewChoice(items...), nil
}

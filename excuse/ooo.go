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

package excuse

var oooTodayStatement = NewChoice(
	NewTerm("It's likely I will be OOO today"),
	NewTerm("I'm going to be OOO"),
	NewTerm("I'm forced to stay OOO today"),
	NewTerm("I can't show up today"))

var oooCauseBadLuck = NewChoice(
	NewTerm("I've been sprayed by a skunk and the smell if too hard to let me focus"),
	NewTerm("the dog ate my shoes"),
	NewTerm("a tree fell on my car"),
	NewTerm("my toe is trapped in the bath tap"),
	NewTerm("I've had a sleepless night"),
	NewTerm("my trousers split on the way to work"),
	NewTerm("I've had a hair dye disaster"),
	NewTerm("my curlers burned my hair, and I had to go to the hairdresser"),
	NewTerm("I am stuck in my house because the door's broken"),
)

var oooCauseCotorep = NewChoice(
	NewTerm("I forgot to come back to work after lunch"),
	NewTerm("my brain went to sleep, and I couldn't wake it up"),
	NewTerm("I forgot what day of the week it was"),
	NewTerm("my dog has had a big fright and I don't want to leave him"),
	NewTerm("I drank too much and fell asleep on someone's floor - I don't know where I am"),
	NewTerm("I woke up late and missed my train... again"),
	NewTerm("I'm in A&E as I got a clothes peg stuck on my tongue"),
	NewTerm("I ate some very spicey chicken wings last night -- It'll be best I stay home"),
	NewTerm("I woke up and unexpectedly had to drive my family to another state"),
)

var oooCausePolice = NewChoice(
	NewTerm("my bus broke down and was held up by robbers"),
	NewTerm("I was arrested as a result of mistaken identity"),
	NewTerm("I totaled my wife's jeep in a collision with a cow"),
	NewTerm("a hitman was looking for me"),
	NewTerm("I eloped. They are after, but won't catch me"),
	NewTerm("I had to be there for my partner's grand jury trial"),
	NewTerm("I had to ship my grandmother's bones to her homeland"),
	NewTerm("someone slipped drugs in my drink last night"),
	NewTerm("my car handbrake broke and it rolled down the hill into a lamppost"),
	NewTerm("I'm using a new contact lens solution and my eyes are watering"),
)

var oooCauseMedicalBare = NewChoice(
	NewTerm("I tripped over my dog and was knocked unconscious"),
	NewTerm("I burned my hand on the toaster"),
	NewTerm("I hurt myself bowling"),
	NewTerm("I've got a sore finger"),
	NewTerm("I have a blocked nose"),
	NewTerm("I was spit on by a venomous snake. I admit, this is uncommon, but it hurts"),
	NewTerm("my kids get sick only get sick on weekdays"),
	NewTerm("a can of baked beans landed on my big toe"),
	NewTerm("I've injured myself during sex"),
	NewTerm("my new girlfriend bit me in a delicate place"),
	NewTerm("I swallowed white spirit"),
	NewTerm("I am hallucinating"),
	NewTerm("I was swimming too fast and smacked my head on the poolside, and it bleeds"),
	NewTerm("I've been bitten by an insect and the larvaes start to get out"),
	NewTerm("someone slipped drugs in my drink last night"),
)

var oooCauseMedical = NewSequence(
	oooCauseMedicalBare,
	NewChoice(
		NewTerm(""),
		NewTerm(", I should go the the hospital"),
		NewTerm(", maybe I should go the the hospital"),
		NewTerm(", I am heading to the doctor"),
		NewTerm(", so it seems I need a MD"),
	),
)

var oooCauseVet = NewChoice(
	NewTerm("my fish is sick"),
	NewTerm("my cat puked last night, it kept me up late. I'm too tired"),
)

func NewOOOMedical() (Generator, error) {
	cause := oooCauseMedical
	return NewChoice(
		NewSequence(oooTodayStatement, conjonction_cause, cause),
		NewSequence(cause, conjonction_consequence, oooTodayStatement)), nil
}

func NewOOO() (Generator, error) {
	cause := NewChoice(
		oooCauseMedical,
		oooCauseVet,
		oooCausePolice,
		oooCauseBadLuck,
		oooCauseCotorep)
	return NewChoice(
		NewSequence(
			oooTodayStatement,
			conjonction_cause,
			cause),
		NewSequence(
			cause,
			conjonction_consequence,
			oooTodayStatement)), nil
}

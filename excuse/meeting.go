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

var noDailyStatement = NewChoice(
	NewTerm("I cannot attend the daily meeting"),
	NewTerm("going to miss the meeting"),
	NewTerm("gonna miss the meeting"),
	NewTerm("I won't be able to attend the daily meeting"))

var conjonction_cause = NewChoice(
	NewTerm("because"),
	NewTerm("'cos"),
	NewTerm("since"))

var conjonction_consequence = NewChoice(
	NewTerm("thus"),
	NewTerm("then"),
	NewTerm("as a consequence"),
	NewTerm("therefore"))

var extenCause = NewChoice(
	NewTerm("I'm feeling sick this morning"),
	NewTerm("I'm not feeling great"),
	NewTerm("I am having all the symptoms of a cold"),
	NewTerm("I think I have a light flu"),
	NewTerm("the network is down in the whole area, currently using a weak phone connection"),
	NewTerm("the network is down, I expect someone to check for the cables right during the meeting"),
	NewTerm("network outage here!"),
	NewTerm("power outage here in the area, I am going downtown to fetch some fuel for the generator"),
	NewTerm("due to harsh weather conditions, I suffer an african-grade network stability"),
	NewTerm("it's blizzard day, snow is piling all around here and it affects the network infrastructure"),
	NewTerm("due to the blizzard, the power supply isn't stable enough to let me attend"),
	NewTerm("it's haircut day! The only time slot available is incompatible with the meeting"),
	NewTerm("it's shower day and ... it takes time"),
	NewTerm("I'm waiting for a parcel delivery"),
	NewTerm("I need to go to the parcel pickup point"),
	NewTerm("it's my naturist day, and I wouldn't expose my perfect body to the team."),
	NewTerm("I woke up with a pretty bad headache"),
	NewTerm("I have a dog grooming errand"),
	NewTerm("my wife has a doctor apointment, I need to take her"),
	NewTerm("my wife is stuck in her car and I have to go help her"),
	NewTerm("I had a hard time last night"),
	NewTerm("I have to file for a new ID card"),
	NewTerm("I need to visit the county clerk office to renew my license plates."),
	NewTerm("I've lost my keys in the river"),
	NewTerm("someone stole my catalytic converter"),
	NewTerm("I have a dentist appointment and it will be hard to speak clearly"),
	NewTerm("my colonoscopy won't be done yet"),
	NewTerm("my prostate exam doesnt happen as expected"),
	NewTerm("the network is down in the area"),
	NewTerm("my nan slipped on a dung (long story short...)"))

func NewNoMeeting() (Node, error) {
	return NewChoice(
		NewSequence(noDailyStatement, conjonction_cause, extenCause),
		NewSequence(extenCause, conjonction_consequence, noDailyStatement)), nil
}

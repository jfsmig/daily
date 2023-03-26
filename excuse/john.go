package excuse

var statement = NewChoice(
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

var cause = NewChoice(
	NewTerm("feeling sick this morning"),
	NewTerm("woke up with a pretty bad headache"),
	NewTerm("not feeling great"),
	NewTerm("I have a dog grooming errand"),
	NewTerm("my wife is stuck in her car"),
	NewTerm("I had a hard time last night"),
	NewTerm("I have to file for a new ID card"),
	NewTerm("I've lost my keys in the river"),
	NewTerm("someone stole my catalytic converter"),
	NewTerm("my colonoscopy won't be done yet"),
	NewTerm("my prostate exam doesnt happen as expected"),
	NewTerm("my nan slipped on a dung"))

var sentence = NewChoice(
	NewSequence(statement, conjonction_cause, cause),
	NewSequence(cause, conjonction_consequence, statement))

func NewJohn() (Node, error) {
	return sentence, nil
}

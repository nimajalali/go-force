package sobjects

import (
	"github.com/nimajalali/go-force/forcejson"
	"reflect"
	"testing"
	"time"
)

type Thing struct {
	BaseSObject
	First *Time `force:"First,omitempty"`
	Next  *Time `force:"Next,omitempty"`
}

func TestDate(t *testing.T) {
	in := Thing{
		First: AsTime(time.Now()),
		// Next: leave empty.
	}

	buf, err := forcejson.Marshal(&in)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	t.Logf("%s", string(buf))

	var out Thing
	err = forcejson.Unmarshal(buf, &out)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Errorf("wrong output:\nexpected: %+v\n     got: %+v", in, out)
	}
}

func TestPrecision(t *testing.T) {
	a := time.Now().Truncate(time.Second)
	b := a.Add(time.Millisecond)
	c := b.Add(time.Second)

	// test that milliseconds don't matter
	if !AsTime(a).Time().Equal(AsTime(b).Time()) {
		t.Errorf("sObject dates differ only by ms, SFDC doesn't store milliseconds! a: %s, b: %s", a, b)
	}
	// test that seconds still matter
	if AsTime(a) == AsTime(c) {
		t.Errorf("dates should differ by seconds but are equal! a: %s, b: %s", a, b)
	}
}

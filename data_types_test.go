package groupme

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type DataTypesSuite struct {
	suite.Suite
}

func (s *DataTypesSuite) TestID_Valid_True() {
	var id string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	s.Assert().True(ValidID(id))
}

func (s *DataTypesSuite) TestID_Valid_False() {
	var id string = "`~!@#$%^&*()_-+={[}]:;\"'<,>.?/|\\"
	s.Assert().False(ValidID(id))
}

func (s *DataTypesSuite) TestTimestamp_FromTime() {
	t := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	expected := Timestamp(0)
	actual := FromTime(t)
	s.Assert().EqualValues(expected, actual)
}

func (s *DataTypesSuite) TestTimestamp_ToTime() {
	t := Timestamp(0)
	expected := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	actual := t.ToTime()
	s.Assert().EqualValues(expected, actual)
}

func (s *DataTypesSuite) TestPhoneNumber_Valid_True() {
	var pn PhoneNumber = "+1 0123456789"
	s.Assert().True(pn.Valid())
}

func (s *DataTypesSuite) TestPhoneNumber_Valid_NoPlus() {
	var pn PhoneNumber = "1 0123456789"
	s.Assert().False(pn.Valid())
}

func (s *DataTypesSuite) TestPhoneNumber_Valid_NoSpace() {
	var pn PhoneNumber = "+10123456789"
	s.Assert().False(pn.Valid())
}

func (s *DataTypesSuite) TestPhoneNumber_Valid_BadLength() {
	var pn PhoneNumber = "+1 01234567890"
	s.Assert().False(pn.Valid())
}

func TestDataTypesSuite(t *testing.T) {
	suite.Run(t, new(DataTypesSuite))
}

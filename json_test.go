package groupme

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type JSONSuite struct {
	suite.Suite
}

func (s *JSONSuite) TestGroup_GetMemberByUserID_Match() {
	m := Member{
		UserID: "123",
	}

	g := Group{
		Members: []*Member{&m},
	}

	actual := g.GetMemberByUserID("123")

	s.Require().NotNil(actual)
	s.Assert().Equal(m, *actual)
}

func (s *JSONSuite) TestGroup_GetMemberByUserID_NoMatch() {
	g := Group{
		Members: []*Member{},
	}

	actual := g.GetMemberByUserID("123")

	s.Require().Nil(actual)
}

func (s *JSONSuite) TestGroup_GetMemberByNickname_Match() {
	m := Member{
		Nickname: "Test User",
	}

	g := Group{
		Members: []*Member{&m},
	}

	actual := g.GetMemberByNickname("Test User")

	s.Require().NotNil(actual)
	s.Assert().Equal(m, *actual)
}

func (s *JSONSuite) TestGroup_GetMemberByNickname_NoMatch() {
	g := Group{
		Members: []*Member{},
	}

	actual := g.GetMemberByNickname("Test User")

	s.Require().Nil(actual)
}

func (s *JSONSuite) TestMarshal_NoError() {
	s.Assert().Equal("{}", marshal(&struct{}{}))
}

func (s *JSONSuite) TestMarshal_Error() {
	var c chan struct{}
	s.Assert().Equal("", marshal(c))
}

func TestJSONSuite(t *testing.T) {
	suite.Run(t, new(JSONSuite))
}

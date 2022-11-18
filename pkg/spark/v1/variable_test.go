package spark_v1

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

/************************************************************************/
// TYPES SUITE
/************************************************************************/

type VariableSuite struct {
	suite.Suite
}

/************************************************************************/
// TESTS
/************************************************************************/

func (s *VariableSuite) Test_Create_Var_Returns_Valid_Var() {
	v := NewVar("test", "application/text", "testValue")
	s.Equal("test", v.Name)
	s.Equal("application/text", v.MimeType)
	s.Equal("testValue", v.Value)
}

/************************************************************************/
// SUITE
/************************************************************************/

func TestVariableSuite(t *testing.T) {
	suite.Run(t, new(VariableSuite))
}
package log

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SuiteWrapper struct {
	suite.Suite
}

func TestRunLogSuite(t *testing.T) {
	suite.Run(t, new(SuiteWrapper))
}

func (s *SuiteWrapper) SetupSuite() {
}

func (s *SuiteWrapper) TestWrapper() {
	// given
	SetLogger("DEBUG", "", "", "", 2, false)

	// then
	Print("[LOG TEST] Print Test OK!!!")
	Printf("[LOG TEST] Printf Test %s\n", "OK!!!")
	Println("[LOG TEST] Println Test OK!!!")

	assert.Panics(s.T(), func() { Panic("[LOG TEST] Panic Test OK!!!") })
	assert.Panics(s.T(), func() { Panicf("[LOG TEST] Panicf Test %s\n", "OK!!!") })
	assert.Panics(s.T(), func() { Panicln("[LOG TEST] Panicln Test OK!!!") })

	//Fatal("[LOG TEST] Fatal Test OK!!!")
	//Fatalf("[LOG TEST] Fatalf Test %s\n", "OK!!!")
	//Fatalln("[LOG TEST] Fatalln Test OK!!!")
}

func (s *SuiteWrapper) TestWrapperWithJson() {
	// given
	SetLogger("DEBUG", "", "", "", 2, true)

	// then
	Print("[LOG TEST] Print Test by Json OK!!!")
	Printf("[LOG TEST] Printf Test by Json %s\n", "OK!!!")
	Println("[LOG TEST] Println Test by Json OK!!!")

	assert.Panics(s.T(), func() { Panic("[LOG TEST] Panic Test by Json OK!!!") })
	assert.Panics(s.T(), func() { Panicf("[LOG TEST] Panicf Test by Json %s\n", "OK!!!") })
	assert.Panics(s.T(), func() { Panicln("[LOG TEST] Panicln Test by Json OK!!!") })

	//Fatal("[LOG TEST] Fatal Test by Json OK!!!")
	//Fatalf("[LOG TEST] Fatalf Test by Json %s\n", "OK!!!")
	//Fatalln("[LOG TEST] Fatalln Test by Json OK!!!")
}

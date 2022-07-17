package basecep

import (
	"testing"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(&APISuite{})

type APISuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *APISuite) TestBrCepResultSanitizeShouldCleanCEP(c *gc.C) {
	var brCepResult = &BrAddress{
		Cep: "78-04-8000",
	}

	brCepResult.Sanitize()

	c.Check(brCepResult.Cep, gc.Equals, "78048000")
}

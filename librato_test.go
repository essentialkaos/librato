package librato

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2015 Essential Kaos                         //
//      Essential Kaos Open Source License <http://essentialkaos.com/ekol?en>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"ek/rand"
	"errors"
	. "gopkg.in/check.v1"
	"testing"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type LibratoSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&LibratoSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (ls *LibratoSuite) TestMetricsValidation(c *C) {
	c.Assert(validateMetrics(&Metrics{}), DeepEquals, errors.New("Metrics struct is not initialized"))

	m, e := NewMetrics(time.Minute)

	c.Assert(m, NotNil)
	c.Assert(e, IsNil)
}

func (ls *LibratoSuite) TestAnnotaionsValidation(c *C) {
	c.Assert(validateAnotations(&Annotations{}), DeepEquals, errors.New("Annotations struct is not initialized"))

	a, e := NewAnnotations("")

	c.Assert(a, IsNil)
	c.Assert(e, NotNil)

	c.Assert(e, DeepEquals, errors.New("Annotations must have non-empty property stream"))

	a, e = NewAnnotations("test")

	c.Assert(a, NotNil)
	c.Assert(e, IsNil)
}

func (ls *LibratoSuite) TestGaugeValidation(c *C) {
	c.Assert(validateGauge(&Gauge{}), DeepEquals, errors.New("Gauge property Name can't be empty"))
	c.Assert(validateGauge(&Gauge{Name: rand.String(256)}), DeepEquals, errors.New("Length of gauge property Name must be 255 or fewer characters"))
	c.Assert(validateGauge(&Gauge{Name: "1", Value: "test"}), DeepEquals, errors.New("Gauge property Value can't be non-numeric"))
	c.Assert(validateGauge(&Gauge{Name: "1", Value: 1, Count: "test"}), DeepEquals, errors.New("Gauge property Count can't be non-numeric"))
	c.Assert(validateGauge(&Gauge{Name: "1", Value: 1, Sum: "test"}), DeepEquals, errors.New("Gauge property Sum can't be non-numeric"))
	c.Assert(validateGauge(&Gauge{Name: "1", Value: 1, Min: "test"}), DeepEquals, errors.New("Gauge property Min can't be non-numeric"))
	c.Assert(validateGauge(&Gauge{Name: "1", Value: 1, Max: "test"}), DeepEquals, errors.New("Gauge property Max can't be non-numeric"))
	c.Assert(validateGauge(&Gauge{Name: "1", Value: 1, SumSquares: "test"}), DeepEquals, errors.New("Gauge property SumSquares can't be non-numeric"))

	c.Assert(validateGauge(&Gauge{Name: "1", Value: 1}), IsNil)
}

func (ls *LibratoSuite) TestAnnotaionValidation(c *C) {
	c.Assert(validateAnotation(&Annotation{}), DeepEquals, errors.New("Annotation property Title can't be empty"))

	c.Assert(validateAnotation(&Annotation{Title: "test"}), IsNil)
}

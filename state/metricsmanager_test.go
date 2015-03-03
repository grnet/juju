// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state_test

import (
	"fmt"
	"time"

	testing "github.com/juju/juju/state/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type metricsManagerSuite struct {
	testing.StateSuite
}

var _ = gc.Suite(&metricsManagerSuite{})

func (s *metricsManagerSuite) TestDefaultsWritten(c *gc.C) {
	mm, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(mm.DocID(), gc.Equals, fmt.Sprintf("%s:metricsManagerKey", s.State.EnvironUUID()))
	c.Assert(mm.LastSuccessfulSend(), gc.DeepEquals, time.Time{})
	c.Assert(mm.ConsecutiveErrors(), gc.Equals, 0)
	c.Assert(mm.GracePeriod(), gc.Equals, 24*7*time.Hour)
}

func (s *metricsManagerSuite) TestMetricsManagerCreatesThenReturns(c *gc.C) {
	mm, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	err = mm.IncrementConsecutiveErrors()
	c.Assert(err, jc.ErrorIsNil)
	mm2, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(mm.ConsecutiveErrors(), gc.Equals, mm2.ConsecutiveErrors())
}

func (s *metricsManagerSuite) TestSetMetricsManagerSuccesfulSend(c *gc.C) {
	mm, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	err = mm.IncrementConsecutiveErrors()
	c.Assert(err, jc.ErrorIsNil)
	now := time.Now().Round(time.Second).UTC()
	err = mm.SetMetricsManagerSuccessfulSend(now)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(mm.LastSuccessfulSend(), gc.DeepEquals, now)
	c.Assert(mm.ConsecutiveErrors(), gc.Equals, 0)

	m, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(m.LastSuccessfulSend().Equal(now), jc.IsTrue)
	c.Assert(mm.ConsecutiveErrors(), gc.Equals, 0)
}

func (s *metricsManagerSuite) TestIncrementConsecutiveErrors(c *gc.C) {
	mm, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(mm.ConsecutiveErrors(), gc.Equals, 0)
	err = mm.IncrementConsecutiveErrors()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(mm.ConsecutiveErrors(), gc.Equals, 1)

	m, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(m.ConsecutiveErrors(), gc.Equals, 1)
}

func (s *metricsManagerSuite) TestMeterStatus(c *gc.C) {
	mm, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	code, info := mm.MeterStatus()
	c.Assert(code, gc.Equals, "GREEN")
	c.Assert(info, gc.Equals, "ok")
	now := time.Now()
	err = mm.SetMetricsManagerSuccessfulSend(now)
	c.Assert(err, jc.ErrorIsNil)
	for i := 0; i < 3; i++ {
		err := mm.IncrementConsecutiveErrors()
		c.Assert(err, jc.ErrorIsNil)
	}
	code, info = mm.MeterStatus()
	c.Assert(code, gc.Equals, "AMBER")
	c.Assert(info, gc.Equals, "failed to send metrics")
	err = mm.SetMetricsManagerSuccessfulSend(now.Add(-(24 * 7 * time.Hour)))
	c.Assert(err, jc.ErrorIsNil)
	for i := 0; i < 3; i++ {
		err := mm.IncrementConsecutiveErrors()
		c.Assert(err, jc.ErrorIsNil)
	}
	code, info = mm.MeterStatus()
	c.Assert(code, gc.Equals, "RED")
	c.Assert(info, gc.Equals, "failed to send metrics, exceeded grace period")

	err = mm.SetMetricsManagerSuccessfulSend(now)
	c.Assert(err, jc.ErrorIsNil)
	code, info = mm.MeterStatus()
	c.Assert(code, gc.Equals, "GREEN")
	c.Assert(info, gc.Equals, "ok")
}

// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"time"

	"github.com/juju/errors"
	"github.com/juju/loggo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

var metricsManagerLogger = loggo.GetLogger("juju.state.metricsmanager")

const (
	defaultGracePeriod                      = 7 * 24 * time.Hour // 1 week in hours
	metricsManagerConsecutiveErrorThreshold = 3
	metricsManagerKey                       = "metricsManagerKey"
)

// MetricsManager stores data about the state of the metrics manager
type MetricsManager struct {
	st  *State
	doc metricsManagerDoc
}

type metricsManagerDoc struct {
	DocID              string        `bson:"_id"`
	LastSuccessfulSend time.Time     `bson:"lastsuccessfulsend"`
	ConsecutiveErrors  int           `bson:"consecutiveerrors"`
	GracePeriod        time.Duration `bson:"graceperiod"`
}

// DocID returns the Document id of the MetricsManager.
func (m *MetricsManager) DocID() string {
	return m.doc.DocID
}

// LastSuccessfulSend returns the time of the last successful send.
func (m *MetricsManager) LastSuccessfulSend() time.Time {
	return m.doc.LastSuccessfulSend
}

// ConsecutiveErrors returns the number of consecutive failures.
func (m *MetricsManager) ConsecutiveErrors() int {
	return m.doc.ConsecutiveErrors
}

// GracePeriod returns the current grace period.
func (m *MetricsManager) GracePeriod() time.Duration {
	return m.doc.GracePeriod
}

// MetricsManager returns an existing metricsmanager, or a new one if non exists.
func (st *State) MetricsManager() (*MetricsManager, error) {
	mm, err := st.getMetricsManager()
	if errors.IsNotFound(err) {
		return st.newMetricsManager()
	} else if err != nil {
		return nil, errors.Trace(err)
	}
	return mm, nil
}

func (st *State) newMetricsManager() (*MetricsManager, error) {
	mm := &MetricsManager{
		st: st,
		doc: metricsManagerDoc{
			DocID:              st.docID(metricsManagerKey),
			LastSuccessfulSend: time.Time{},
			ConsecutiveErrors:  0,
			GracePeriod:        defaultGracePeriod,
		}}
	buildTxn := func(attempt int) ([]txn.Op, error) {
		if attempt > 0 {
			coll, closer := st.getCollection(metricsManagerC)
			defer closer()
			n, err := coll.FindId(st.docID(metricsManagerKey)).Count()
			if err != nil {
				return nil, errors.Trace(err)
			}
			if n > 0 {
				return nil, errors.NotFoundf("metrics manager")
			}
		}
		return []txn.Op{{
			C:      metricsManagerC,
			Id:     st.docID(metricsManagerKey),
			Assert: txn.DocMissing,
			Insert: mm.doc,
		},
		}, nil
	}
	err := st.run(buildTxn)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return mm, nil
}

func (st *State) getMetricsManager() (*MetricsManager, error) {
	coll, closer := st.getCollection(metricsManagerC)
	defer closer()
	var doc metricsManagerDoc
	err := coll.FindId(st.docID(metricsManagerKey)).One(&doc)
	if err == mgo.ErrNotFound {
		return nil, errors.NotFoundf("metrics manager")
	} else if err != nil {
		return nil, errors.Trace(err)
	}
	return &MetricsManager{st: st, doc: doc}, nil
}

func (m *MetricsManager) updateMetricsManager(update bson.M) error {
	buildTxn := func(attempt int) ([]txn.Op, error) {
		return []txn.Op{{
			C:      metricsManagerC,
			Id:     m.st.docID(metricsManagerKey),
			Assert: txn.DocExists,
			Update: update,
		}}, nil
	}
	return m.st.run(buildTxn)
}

// SetMetricsManagerSuccessfulSend sets the last successful send time to the input time.
func (m *MetricsManager) SetMetricsManagerSuccessfulSend(t time.Time) error {
	err := m.updateMetricsManager(
		bson.M{"$set": bson.M{
			"lastsuccessfulsend": t.UTC(),
			"consecutiveerrors":  0,
		}},
	)
	if err != nil {
		return errors.Trace(err)
	}
	m.doc.LastSuccessfulSend = t.UTC()
	m.doc.ConsecutiveErrors = 0
	return nil
}

// IncrementConsecutiveErrors adds 1 to the consecutive errors count.
func (m *MetricsManager) IncrementConsecutiveErrors() error {
	err := m.updateMetricsManager(
		bson.M{"$inc": bson.M{"consecutiveerrors": 1}},
	)
	if err != nil {
		return errors.Trace(err)
	}
	m.doc.ConsecutiveErrors++
	return nil
}

func (m *MetricsManager) gracePeriodExceeded() bool {
	now := time.Now()
	t := m.LastSuccessfulSend().Add(m.GracePeriod())
	return t.Before(now) || t.Equal(now)
}

// MeterStatus returns the overall state of the MetricsManager as a meter status summary.
func (m *MetricsManager) MeterStatus() (code, info string) {
	if m.ConsecutiveErrors() < metricsManagerConsecutiveErrorThreshold {
		return string(MeterGreen), "ok"
	}
	if m.gracePeriodExceeded() {
		return string(MeterRed), "failed to send metrics, exceeded grace period"
	}
	return string(MeterAmber), "failed to send metrics"
}

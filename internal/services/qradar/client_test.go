package qradar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var cli = NewQRadarClient("https://qradar.redcarbon.ai", "xxx", true)

func TestRetrieveQRadarOffenseUrl(t *testing.T) {
	t.Skip()

	url := cli.RetrieveOffenseUrl(123)
	assert.NotNil(t, url)

	assert.Equal(t, "https://qradar.redcarbon.ai/console/do/sem/offensesummary?appName=Sem&pageId=OffenseSummary&summaryId=123", *url)
}

func TestRetrieveOffenceEvents(t *testing.T) {
	t.Skip()

	_, err := cli.SearchOffenseEvents(t.Context(), 208080, 1745906350195)

	assert.Nil(t, err)

	assert.FailNow(t, "stop here")
}

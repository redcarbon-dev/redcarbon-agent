package qradar

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetrieveQRadarOffenseUrl(t *testing.T) {
	cli := NewQRadarClient("https://qradar.redcarbon.ai", "xxx", true)

	url := cli.RetrieveOffenseUrl(123)
	assert.NotNil(t, url)

	assert.Equal(t, "https://qradar.redcarbon.ai/console/do/sem/offensesummary?appName=Sem&pageId=OffenseSummary&summaryId=123", *url)
}

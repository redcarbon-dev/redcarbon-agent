package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"pkg.redcarbon.ai/internal/utils"
)

var data = []map[string]string{
	{
		"Column1": "f",
		"Column2": "o",
		"Column3": "o",
	},
	{
		"Column1": "f",
		"Column2": "o",
		"Column3": "o",
	},
	{
		"Column1": "o",
		"Column2": "o",
		"Column3": "f",
	},
}

func TestGetUniqueDataForColumnInMap(t *testing.T) {
	c1 := utils.GetUniqueDataForColumnInMap(data, "Column1")
	assert.Len(t, c1, 2)
	assert.Equal(t, []string{"f", "o"}, c1)

	c2 := utils.GetUniqueDataForColumnInMap(data, "Column2")
	assert.Len(t, c2, 1)
	assert.Equal(t, []string{"o"}, c2)
}

func TestGroupingMapByColumn(t *testing.T) {
	groups := utils.GroupMapByColumn(data, "Column1")
	assert.Len(t, groups, 2)
	assert.Len(t, groups["f"], 2)
	assert.Len(t, groups["o"], 1)

	groups = utils.GroupMapByColumn(data, "Column2")
	assert.Len(t, groups, 1)
	assert.Len(t, groups["f"], 0)
	assert.Len(t, groups["o"], 3)

	groups = utils.GroupMapByColumn(data, "Column3")
	assert.Len(t, groups, 2)
	assert.Len(t, groups["f"], 1)
	assert.Len(t, groups["o"], 2)
}

package fixtures

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetBlanks(t *testing.T) {
	teams := map[int]*teamBucket{
		1: {team: team{Id: 1, Name: "Team A"}},
		2: {team: team{Id: 2, Name: "team B"}},
		3: {team: team{Id: 3, Name: "Team C"}},
	}
	gameweekMap := map[int][]fixture{
		1: {
			fixture{HomeTeamId: 1, AwayTeamId: 2},
		},
	}

	result := getBlanks(gameweekMap, teams)

	assert.Len(t, result[1], 1)
	assert.Equal(t, "Team C", result[1][0].team.Name)
}

func TestGetDoubles(t *testing.T) {
	gameweekMap := map[int][]fixture{
		1: {
			fixture{Id: 1, HomeTeamId: 1, AwayTeamId: 2},
			fixture{Id: 2, HomeTeamId: 1, AwayTeamId: 3},
		},
	}

	result := getDoubles(gameweekMap)

	assert.Len(t, result[1], 2)
	assert.Equal(t, 1, result[1][0].HomeTeamId)
	assert.Equal(t, 1, result[1][1].HomeTeamId)
}

func TestIsDouble(t *testing.T) {
	fixtures := []fixture{
		{Id: 1, HomeTeamId: 1, AwayTeamId: 2},
		{Id: 2, HomeTeamId: 1, AwayTeamId: 3},
	}

	result := isDouble(fixtures[0], fixtures)
	assert.True(t, result)

	nonDoubleFixture := fixture{Id: 3, HomeTeamId: 4, AwayTeamId: 5}
	assert.False(t, isDouble(nonDoubleFixture, fixtures))
}

func TestRemove(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	result := remove(slice, 2)

	assert.Equal(t, []int{1, 3, 4}, result)
}

func TestGetSortedKeysFromMap(t *testing.T) {
	input := map[int]string{
		3: "three",
		1: "one",
		2: "two",
	}

	result := getSortedKeysFromMap(input)

	assert.Equal(t, []int{1, 2, 3}, result)
}

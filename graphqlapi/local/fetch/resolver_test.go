/*                          _       _
 *__      _____  __ ___   ___  __ _| |_ ___
 *\ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
 * \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
 *  \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
 *
 * Copyright © 2016 - 2019 Weaviate. All rights reserved.
 * LICENSE: https://github.com/creativesoftwarefdn/weaviate/blob/develop/LICENSE.md
 * DESIGN & CONCEPT: Bob van Luijt (@bobvanluijt)
 * CONTACT: hello@creativesoftwarefdn.org
 */
package fetch

import (
	"testing"

	"github.com/creativesoftwarefdn/weaviate/contextionary"
	"github.com/creativesoftwarefdn/weaviate/models"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name                                string
	query                               string
	expectedParamsToConnector           *Params
	expectedSearchParamsToContextionary []contextionary.SearchParams
	resolverReturn                      interface{}
	expectedResults                     []result
}

type testCases []testCase

type result struct {
	pathToField   []string
	expectedValue interface{}
}

func Test_Resolve(t *testing.T) {
	t.Parallel()

	tests := testCases{
		testCase{
			name: "single prop: mean",
			query: `
			{
				Fetch {
					Things(where: {
						class: {
							name: "bestclass"
							certainty: 0.8
							keywords: [{value: "foo", weight: 0.9}]
						},
						properties: {
							name: "bestproperty"
							certainty: 0.8
							keywords: [{value: "bar", weight: 0.9}]
							operator: Equal
							valueString: "some-value"
						},
					}) {
						beacon certainty
					}
				}
			}`,
			expectedSearchParamsToContextionary: []contextionary.SearchParams{
				{
					SearchType: contextionary.SearchTypeClass,
					Name:       "bestclass",
					Certainty:  0.8,
					Keywords: models.SemanticSchemaKeywords{{
						Keyword: "foo",
						Weight:  0.9,
					}},
				},
				{
					SearchType: contextionary.SearchTypeProperty,
					Name:       "bestproperty",
					Certainty:  0.8,
					Keywords: models.SemanticSchemaKeywords{{
						Keyword: "bar",
						Weight:  0.9,
					}},
				},
			},
			expectedParamsToConnector: &Params{
				PossibleClassNames: contextionary.SearchResults{
					Type: contextionary.SearchTypeClass,
					Results: []contextionary.SearchResult{{
						Name:      "bestclass",
						Certainty: 0.95,
					}, {
						Name:      "bestclassalternative",
						Certainty: 0.85,
					}},
				},
				Properties: []Property{
					{
						PossibleNames: contextionary.SearchResults{
							Type: contextionary.SearchTypeProperty,
							Results: []contextionary.SearchResult{{
								Name:      "bestproperty",
								Certainty: 0.95,
							}, {
								Name:      "bestpropertyalternative",
								Certainty: 0.85,
							}},
						},
					},
				},
			},
			resolverReturn: []interface{}{
				map[string]interface{}{
					"beacon":    "weaviate://peerName/things/uuid1",
					"certainty": 0.7,
				},
			},
			expectedResults: []result{{
				pathToField: []string{"Fetch", "Things"},
				expectedValue: []interface{}{
					map[string]interface{}{
						"beacon":    "weaviate://peerName/things/uuid1",
						"certainty": 0.7,
					},
				},
			}},
		},
	}

	tests.AssertExtraction(t)
}

func (tests testCases) AssertExtraction(t *testing.T) {
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			c11y := newMockContextionary()

			if len(testCase.expectedSearchParamsToContextionary) != 2 {
				t.Fatalf("test setup incorrect: expectedSearchParamsToContextionary must have len 2")
			}

			c11y.On("SchemaSearch", testCase.expectedSearchParamsToContextionary[0]).Once()
			c11y.On("SchemaSearch", testCase.expectedSearchParamsToContextionary[1]).Once()

			resolver := newMockResolver(c11y)

			resolver.On("LocalFetchKindClass", testCase.expectedParamsToConnector).
				Return(testCase.resolverReturn, nil).Once()

			result := resolver.AssertResolve(t, testCase.query)
			c11y.AssertExpectations(t)

			for _, expectedResult := range testCase.expectedResults {
				value := result.Get(expectedResult.pathToField...).Result

				assert.Equal(t, expectedResult.expectedValue, value)
			}
		})
	}
}
//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2019 SeMI Holding B.V. (registered @ Dutch Chamber of Commerce no 75221632). All rights reserved.
//  LICENSE WEAVIATE OPEN SOURCE: https://www.semi.technology/playbook/playbook/contract-weaviate-OSS.html
//  LICENSE WEAVIATE ENTERPRISE: https://www.semi.technology/playbook/contract-weaviate-enterprise.html
//  CONCEPT: Bob van Luijt (@bobvanluijt)
//  CONTACT: hello@semi.technology
//

package get

import (
	"testing"

	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/entities/schema"
	"github.com/semi-technologies/weaviate/usecases/network/crossrefs"
	"github.com/stretchr/testify/assert"
)

func TestExtractEmptySchema(t *testing.T) {
	schema := &schema.Schema{
		Actions: nil,
		Things:  nil,
	}

	result := extractNetworkRefClassNames(schema)
	assert.Equal(t, []crossrefs.NetworkClass{}, result, "should be an empty list")
}

func TestExtractSchemaWithPrimitiveActions(t *testing.T) {
	schema := &schema.Schema{
		Actions: &models.Schema{
			Classes: []*models.Class{
				&models.Class{
					Class: "BestAction",
					Properties: []*models.Property{
						&models.Property{
							DataType: []string{"string"},
							Name:     "bestStringProp",
						},
					},
				}},
		},
		Things: nil,
	}

	result := extractNetworkRefClassNames(schema)
	assert.Equal(t, []crossrefs.NetworkClass{}, result, "should be an empty list")
}

func TestExtractSchemaWithPrimitiveThings(t *testing.T) {
	schema := &schema.Schema{
		Things: &models.Schema{
			Classes: []*models.Class{
				&models.Class{
					Class: "BestThing",
					Properties: []*models.Property{
						&models.Property{
							DataType: []string{"string"},
							Name:     "bestStringProp",
						},
					},
				}},
		},
		Actions: nil,
	}

	result := extractNetworkRefClassNames(schema)
	assert.Equal(t, []crossrefs.NetworkClass{}, result, "should be an empty list")
}

func TestExtractSchemaWithThingsWithLocalRefs(t *testing.T) {
	schema := &schema.Schema{
		Things: &models.Schema{
			Classes: []*models.Class{
				&models.Class{
					Class: "BestThing",
					Properties: []*models.Property{
						&models.Property{
							DataType: []string{"AnotherFairlyGoodThing"},
							Name:     "BestReference",
						},
					},
				}},
		},
		Actions: nil,
	}

	result := extractNetworkRefClassNames(schema)
	assert.Equal(t, []crossrefs.NetworkClass{}, result, "should be an empty list")
}

func TestExtractSchemaWithThingsWithNetworkRefs(t *testing.T) {
	schema := &schema.Schema{
		Things: &models.Schema{
			Classes: []*models.Class{
				&models.Class{
					Class: "BestThing",
					Properties: []*models.Property{
						&models.Property{
							DataType: []string{"OtherInstance/TheBestThing"},
							Name:     "BestReference",
						},
					},
				},
				&models.Class{
					Class: "WorstThing",
					Properties: []*models.Property{
						&models.Property{
							DataType: []string{"OtherInstance/TheWorstThing"},
							Name:     "WorstReference",
						},
						&models.Property{
							DataType: []string{"OtherInstance/TheMediocreThing"},
							Name:     "MediocreReference",
						},
					},
				},
			},
		},
		Actions: nil,
	}

	result := extractNetworkRefClassNames(schema)
	assert.Equal(t, []crossrefs.NetworkClass{
		{PeerName: "OtherInstance", ClassName: "TheBestThing"},
		{PeerName: "OtherInstance", ClassName: "TheWorstThing"},
		{PeerName: "OtherInstance", ClassName: "TheMediocreThing"},
	}, result, "should find the network classes")
}

func TestExtractSchemaWithActionsWithNetworkRefs(t *testing.T) {
	schema := &schema.Schema{
		Actions: &models.Schema{
			Classes: []*models.Class{
				&models.Class{
					Class: "BestAction",
					Properties: []*models.Property{
						&models.Property{
							DataType: []string{"OtherInstance/TheBestThing"},
							Name:     "BestReference",
						},
					},
				},
				&models.Class{
					Class: "WorstThing",
					Properties: []*models.Property{
						&models.Property{
							DataType: []string{"OtherInstance/TheWorstThing"},
							Name:     "WorstReference",
						},
						&models.Property{
							DataType: []string{"OtherInstance/TheMediocreThing"},
							Name:     "MediocreReference",
						},
					},
				},
			},
		},
		Things: nil,
	}

	result := extractNetworkRefClassNames(schema)
	assert.Equal(t, []crossrefs.NetworkClass{
		{PeerName: "OtherInstance", ClassName: "TheBestThing"},
		{PeerName: "OtherInstance", ClassName: "TheWorstThing"},
		{PeerName: "OtherInstance", ClassName: "TheMediocreThing"},
	}, result, "should find the network classes")
}

func TestExtractSchemaWithDuplicates(t *testing.T) {
	schema := &schema.Schema{
		Actions: &models.Schema{
			Classes: []*models.Class{
				&models.Class{
					Class: "BestAction",
					Properties: []*models.Property{
						&models.Property{
							DataType: []string{"OtherInstance/TheBestThing"},
							Name:     "BestReference",
						},
					},
				},
				&models.Class{
					Class: "WorstThing",
					Properties: []*models.Property{
						&models.Property{
							DataType: []string{"OtherInstance/TheBestThing"},
							Name:     "WorstReference",
						},
						&models.Property{
							DataType: []string{"OtherInstance/TheBestThing"},
							Name:     "MediocreReference",
						},
					},
				},
			},
		},
		Things: &models.Schema{
			Classes: []*models.Class{
				&models.Class{
					Class: "BestThing",
					Properties: []*models.Property{
						&models.Property{
							DataType: []string{"OtherInstance/TheBestThing"},
							Name:     "BestReference",
						},
					},
				},
				&models.Class{
					Class: "WorstThing",
					Properties: []*models.Property{
						&models.Property{
							DataType: []string{"OtherInstance/TheBestThing"},
							Name:     "WorstReference",
						},
						&models.Property{
							DataType: []string{"OtherInstance/TheBestThing"},
							Name:     "MediocreReference",
						},
					},
				},
			},
		},
	}

	result := extractNetworkRefClassNames(schema)
	assert.Equal(t,
		[]crossrefs.NetworkClass{{PeerName: "OtherInstance", ClassName: "TheBestThing"}},
		result, "should remove duplicates")
}

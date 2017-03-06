package v2_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_RequestDetailsViewV1_GetQuery_SortsQueryString(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.RequestDetailsViewV1{
		Query: util.StringToPointer("b=b&a=a"),
	}
	queryString := unit.GetQuery()
	Expect(queryString).ToNot(BeNil())

	Expect(*queryString).To(Equal("a=a&b=b"))
}

func Test_RequestDetailsViewV1_GetQuery_ReturnsNilIfNil(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.RequestDetailsViewV1{
		Query: nil,
	}
	queryString := unit.GetQuery()
	Expect(queryString).To(BeNil())
}

func Test_SimulationViewV1_Upgrade_ReturnsAV2Simulation(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV1{
		v2.DataViewV1{
			RequestResponsePairViewV1: []v2.RequestResponsePairViewV1{
				v2.RequestResponsePairViewV1{
					Request: v2.RequestDetailsViewV1{
						Scheme:      util.StringToPointer("http"),
						Body:        util.StringToPointer("body"),
						Destination: util.StringToPointer("destination"),
						Method:      util.StringToPointer("GET"),
						Path:        util.StringToPointer("/path"),
						Query:       util.StringToPointer("query=query"),
						Headers: map[string][]string{
							"Test": []string{"headers"},
						},
					},
					Response: v2.ResponseDetailsView{
						Status:      200,
						Body:        "body",
						EncodedBody: false,
						Headers: map[string][]string{
							"Test": []string{"headers"},
						},
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion:   "v1",
			HoverflyVersion: "test",
			TimeExported:    "today",
		},
	}

	simulationViewV2 := unit.Upgrade()

	Expect(simulationViewV2.RequestResponsePairs).To(HaveLen(1))

	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Scheme).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("http"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Body).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("body"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Destination).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("destination"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Method).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("GET"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Path).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("/path"),
	}))
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Query).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("query=query"),
	}))
	Expect(simulationViewV2.RequestResponsePairs[0].Request.Headers).To(HaveKeyWithValue("Test", []string{"headers"}))

	Expect(simulationViewV2.RequestResponsePairs[0].Response.Status).To(Equal(200))
	Expect(simulationViewV2.RequestResponsePairs[0].Response.Body).To(Equal("body"))
	Expect(simulationViewV2.RequestResponsePairs[0].Response.EncodedBody).To(BeFalse())
	Expect(simulationViewV2.RequestResponsePairs[0].Response.Headers).To(HaveKeyWithValue("Test", []string{"headers"}))

	Expect(simulationViewV2.SchemaVersion).To(Equal("v2"))
	Expect(simulationViewV2.HoverflyVersion).To(Equal("test"))
	Expect(simulationViewV2.TimeExported).To(Equal("today"))
}

func Test_SimulationViewV1_Upgrade_CanReturnAnIncompleteRequest(t *testing.T) {
	RegisterTestingT(t)

	unit := v2.SimulationViewV1{
		v2.DataViewV1{
			RequestResponsePairViewV1: []v2.RequestResponsePairViewV1{
				v2.RequestResponsePairViewV1{
					Request: v2.RequestDetailsViewV1{
						Method: util.StringToPointer("POST"),
					},
					Response: v2.ResponseDetailsView{
						Status:      200,
						Body:        "body",
						EncodedBody: false,
						Headers: map[string][]string{
							"Test": []string{"headers"},
						},
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion:   "v1",
			HoverflyVersion: "test",
			TimeExported:    "today",
		},
	}

	simulationViewV2 := unit.Upgrade()

	Expect(simulationViewV2.RequestResponsePairs).To(HaveLen(1))

	Expect(simulationViewV2.RequestResponsePairs[0].Request.Scheme).To(BeNil())
	Expect(simulationViewV2.RequestResponsePairs[0].Request.Body).To(BeNil())
	Expect(simulationViewV2.RequestResponsePairs[0].Request.Destination).To(BeNil())
	Expect(*simulationViewV2.RequestResponsePairs[0].Request.Method).To(Equal(v2.RequestFieldMatchersView{
		ExactMatch: util.StringToPointer("POST"),
	}))
	Expect(simulationViewV2.RequestResponsePairs[0].Request.Path).To(BeNil())
	Expect(simulationViewV2.RequestResponsePairs[0].Request.Query).To(BeNil())
	Expect(simulationViewV2.RequestResponsePairs[0].Request.Headers).To(BeNil())

	Expect(simulationViewV2.RequestResponsePairs[0].Response.Status).To(Equal(200))
	Expect(simulationViewV2.RequestResponsePairs[0].Response.Body).To(Equal("body"))
	Expect(simulationViewV2.RequestResponsePairs[0].Response.EncodedBody).To(BeFalse())
	Expect(simulationViewV2.RequestResponsePairs[0].Response.Headers).To(HaveKeyWithValue("Test", []string{"headers"}))
}
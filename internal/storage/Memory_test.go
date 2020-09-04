package storage

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInMemory(t *testing.T) {
	Convey("Given an in-memory storage", t, func() {
		store := InMemory()

		Convey("When no agent is known", func() {
			agents := store.Get("nope", "nope")

			Convey("An empty list is returned", func() {
				So(agents, ShouldHaveLength, 0)
			})
		})

		Convey("When a single agent is registered", func() {
			store.Add("client", "service", Agent{
				AgentID:        "id",
				PrivateAddress: "addr",
			})
			agents := store.Get("client", "service")

			Convey("It is returned", func() {
				So(agents, ShouldHaveLength, 1)
				So(agents[0].AgentID, ShouldEqual, "id")
				So(agents[0].PrivateAddress, ShouldEqual, "addr")
			})
		})

		Convey("When a two agents with the same id are registered", func() {
			store.Add("client", "service", Agent{
				AgentID:        "id",
				PrivateAddress: "addr",
			})
			store.Add("client", "service", Agent{
				AgentID:        "id",
				PrivateAddress: "updated",
			})
			agents := store.Get("client", "service")

			Convey("It is overwritten", func() {
				So(agents, ShouldHaveLength, 1)
				So(agents[0].AgentID, ShouldEqual, "id")
				So(agents[0].PrivateAddress, ShouldEqual, "updated")
			})
		})

		Convey("When a two agents with unique ids are registered", func() {
			store.Add("client", "service", Agent{
				AgentID:        "id",
				PrivateAddress: "addr",
			})
			store.Add("client", "service", Agent{
				AgentID:        "id2",
				PrivateAddress: "updated",
			})
			agents := store.Get("client", "service")

			Convey("They are both registered", func() {
				So(agents, ShouldHaveLength, 2)
			})
		})

		Convey("When a two agents differen client ids are registered", func() {
			store.Add("client1", "service", Agent{
				AgentID:        "id",
				PrivateAddress: "addr",
			})
			store.Add("client2", "service", Agent{
				AgentID:        "id2",
				PrivateAddress: "updated",
			})
			agents := store.Get("client1", "service")

			Convey("Only one is returned", func() {
				So(agents, ShouldHaveLength, 1)
			})
		})
	})
}

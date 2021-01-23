// +build conformance

package tst

import (
	"context"
	"testing"

	"github.com/venturemark/apigengo/pkg/pbf/audience"
	"github.com/venturemark/apigengo/pkg/pbf/timeline"

	"github.com/venturemark/fmz/pkg/client"
)

// Test_Timeline_001 ensures that the lifecycle of timelines is covered from
// creation to deletion.
func Test_Timeline_001(t *testing.T) {
	var err error

	var cli *client.Client
	{
		c := client.Config{}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		defer cli.Connection().Close()
	}

	var ti1 string
	{
		i := &timeline.CreateI{
			Obj: &timeline.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
				Property: &timeline.CreateI_Obj_Property{
					Name: "Marketing Campaign",
				},
			},
		}

		o, err := cli.Timeline().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ti1 = s
	}

	var ti2 string
	{
		i := &timeline.CreateI{
			Obj: &timeline.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
				Property: &timeline.CreateI_Obj_Property{
					Name: "Internal Project",
				},
			},
		}

		o, err := cli.Timeline().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ti2 = s
	}

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"organization.venturemark.co/id": "1",
						"user.venturemark.co/id":         "1",
					},
				},
			},
		}

		o, err := cli.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 2 {
			t.Fatal("there must be two timelines")
		}

		if o.Obj[0].Property.Name != "Internal Project" {
			t.Fatal("name must be Internal Project")
		}
		if o.Obj[1].Property.Name != "Marketing Campaign" {
			t.Fatal("name must be Marketing Campaign")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: &timeline.UpdateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     ti1,
					"user.venturemark.co/id":         "1",
				},
				Property: &timeline.UpdateI_Obj_Property{
					Stat: toStringP("archived"),
				},
			},
		}

		o, err := cli.Timeline().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "updated" {
			t.Fatal("status must be updated")
		}
	}

	{
		i := &timeline.DeleteI{
			Obj: &timeline.DeleteI_Obj{
				Metadata: map[string]string{
					"timeline.venturemark.co/id":     ti1,
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
			},
		}

		o, err := cli.Timeline().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: &timeline.UpdateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     ti2,
					"user.venturemark.co/id":         "1",
				},
				Property: &timeline.UpdateI_Obj_Property{
					Stat: toStringP("archived"),
				},
			},
		}

		o, err := cli.Timeline().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "updated" {
			t.Fatal("status must be updated")
		}
	}

	{
		i := &timeline.DeleteI{
			Obj: &timeline.DeleteI_Obj{
				Metadata: map[string]string{
					"timeline.venturemark.co/id":     ti2,
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
			},
		}

		o, err := cli.Timeline().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"organization.venturemark.co/id": "1",
						"user.venturemark.co/id":         "1",
					},
				},
			},
		}

		o, err := cli.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 0 {
			t.Fatal("there must be zero timelines")
		}
	}
}

// Test_Timeline_002 ensures that timeline names are unique.
func Test_Timeline_002(t *testing.T) {
	var err error

	var cli *client.Client
	{
		c := client.Config{}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		defer cli.Connection().Close()
	}

	var tid string
	{
		i := &timeline.CreateI{
			Obj: &timeline.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
				Property: &timeline.CreateI_Obj_Property{
					Name: "Marketing Campaign",
				},
			},
		}

		o, err := cli.Timeline().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		tid = s
	}

	{
		i := &timeline.CreateI{
			Obj: &timeline.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
				Property: &timeline.CreateI_Obj_Property{
					Name: "Marketing Campaign",
				},
			},
		}

		_, err := cli.Timeline().Create(context.Background(), i)
		if err == nil {
			t.Fatal("name must be unique")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: &timeline.UpdateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     tid,
					"user.venturemark.co/id":         "1",
				},
				Property: &timeline.UpdateI_Obj_Property{
					Stat: toStringP("archived"),
				},
			},
		}

		o, err := cli.Timeline().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "updated" {
			t.Fatal("status must be updated")
		}
	}

	{
		i := &timeline.DeleteI{
			Obj: &timeline.DeleteI_Obj{
				Metadata: map[string]string{
					"timeline.venturemark.co/id":     tid,
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
			},
		}

		o, err := cli.Timeline().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}
}

// Test_Timeline_003 ensures that the timeline state can be updated.
func Test_Timeline_003(t *testing.T) {
	var err error

	var cli *client.Client
	{
		c := client.Config{}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		defer cli.Connection().Close()
	}

	var tid string
	{
		i := &timeline.CreateI{
			Obj: &timeline.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
				Property: &timeline.CreateI_Obj_Property{
					Name: "Marketing Campaign",
				},
			},
		}

		o, err := cli.Timeline().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		tid = s
	}

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"organization.venturemark.co/id": "1",
						"user.venturemark.co/id":         "1",
					},
				},
			},
		}

		o, err := cli.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one timeline")
		}

		if o.Obj[0].Property.Desc != "" {
			t.Fatal("desc must be empty")
		}
		if o.Obj[0].Property.Name != "Marketing Campaign" {
			t.Fatal("name must be Internal Project")
		}
		if o.Obj[0].Property.Stat != "active" {
			t.Fatal("stat must be active")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: &timeline.UpdateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     tid,
					"user.venturemark.co/id":         "1",
				},
				Property: &timeline.UpdateI_Obj_Property{
					Stat: toStringP("archived"),
				},
			},
		}

		o, err := cli.Timeline().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "updated" {
			t.Fatal("status must be updated")
		}
	}

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"organization.venturemark.co/id": "1",
						"user.venturemark.co/id":         "1",
					},
				},
			},
		}

		o, err := cli.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one timeline")
		}

		if o.Obj[0].Property.Desc != "" {
			t.Fatal("desc must be empty")
		}
		if o.Obj[0].Property.Name != "Marketing Campaign" {
			t.Fatal("name must be Internal Project")
		}
		if o.Obj[0].Property.Stat != "archived" {
			t.Fatal("stat must be archived")
		}
	}

	{
		i := &timeline.DeleteI{
			Obj: &timeline.DeleteI_Obj{
				Metadata: map[string]string{
					"timeline.venturemark.co/id":     tid,
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
			},
		}

		o, err := cli.Timeline().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}
}

// Test_Timeline_004 ensures that timelines not having the archived state cannot
// be deleted.
func Test_Timeline_004(t *testing.T) {
	var err error

	var cli *client.Client
	{
		c := client.Config{}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		defer cli.Connection().Close()
	}

	var tid string
	{
		i := &timeline.CreateI{
			Obj: &timeline.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
				Property: &timeline.CreateI_Obj_Property{
					Name: "Marketing Campaign",
				},
			},
		}

		o, err := cli.Timeline().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		tid = s
	}

	{
		i := &timeline.DeleteI{
			Obj: &timeline.DeleteI_Obj{
				Metadata: map[string]string{
					"timeline.venturemark.co/id":     tid,
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
			},
		}

		_, err := cli.Timeline().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("state must be archived for deletion")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: &timeline.UpdateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     tid,
					"user.venturemark.co/id":         "1",
				},
				Property: &timeline.UpdateI_Obj_Property{
					Stat: toStringP("archived"),
				},
			},
		}

		o, err := cli.Timeline().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "updated" {
			t.Fatal("status must be updated")
		}
	}

	{
		i := &timeline.DeleteI{
			Obj: &timeline.DeleteI_Obj{
				Metadata: map[string]string{
					"timeline.venturemark.co/id":     tid,
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
			},
		}

		o, err := cli.Timeline().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}
}

// Test_Timeline_005 ensures that timelines can be shown to users of audiences
// allowed to see certain timelines. Users should only be able to see a timeline
// if they are part of the audience the timeline is configured with.
func Test_Timeline_005(t *testing.T) {
	var err error

	var cli *client.Client
	{
		c := client.Config{}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		defer cli.Connection().Close()
	}

	var ui1 string
	var ui2 string
	{
		ui1 = "1"
		ui2 = "2"
	}

	var ti1 string
	{
		i := &timeline.CreateI{
			Obj: &timeline.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         ui1,
				},
				Property: &timeline.CreateI_Obj_Property{
					Name: "Marketing Campaign",
				},
			},
		}

		o, err := cli.Timeline().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ti1 = s
	}

	var ti2 string
	{
		i := &timeline.CreateI{
			Obj: &timeline.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         ui2,
				},
				Property: &timeline.CreateI_Obj_Property{
					Name: "Internal Project",
				},
			},
		}

		o, err := cli.Timeline().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ti2 = s
	}

	var ai1 string
	{
		i := &audience.CreateI{
			Obj: &audience.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         ui1,
				},
				Property: &audience.CreateI_Obj_Property{
					Name: "Employees",
					Tmln: []string{
						ti1,
					},
					User: []string{
						"1",
					},
				},
			},
		}

		o, err := cli.Audience().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["audience.venturemark.co/id"]
		if !ok {
			t.Fatal("audience ID must not be empty")
		}

		ai1 = s
	}

	var ai2 string
	{
		i := &audience.CreateI{
			Obj: &audience.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         ui2,
				},
				Property: &audience.CreateI_Obj_Property{
					Name: "Vendors",
					Tmln: []string{
						ti2,
					},
					User: []string{
						"2",
					},
				},
			},
		}

		o, err := cli.Audience().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["audience.venturemark.co/id"]
		if !ok {
			t.Fatal("audience ID must not be empty")
		}

		ai2 = s
	}

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"organization.venturemark.co/id": "1",
						"user.venturemark.co/id":         ui1,
					},
				},
			},
		}

		o, err := cli.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		// Once we have a permission model in place there must only be one
		// timeline.
		if len(o.Obj) != 2 {
			t.Fatal("there must be two timelines")
		}
	}

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"organization.venturemark.co/id": "1",
						"user.venturemark.co/id":         ui2,
					},
				},
			},
		}

		o, err := cli.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		// Once we have a permission model in place there must only be one
		// timeline.
		if len(o.Obj) != 2 {
			t.Fatal("there must be two timelines")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: &timeline.UpdateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     ti1,
					"user.venturemark.co/id":         ui1,
				},
				Property: &timeline.UpdateI_Obj_Property{
					Stat: toStringP("archived"),
				},
			},
		}

		o, err := cli.Timeline().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "updated" {
			t.Fatal("status must be updated")
		}
	}

	{
		i := &timeline.DeleteI{
			Obj: &timeline.DeleteI_Obj{
				Metadata: map[string]string{
					"timeline.venturemark.co/id":     ti1,
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         ui1,
				},
			},
		}

		o, err := cli.Timeline().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: &timeline.UpdateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     ti2,
					"user.venturemark.co/id":         ui2,
				},
				Property: &timeline.UpdateI_Obj_Property{
					Stat: toStringP("archived"),
				},
			},
		}

		o, err := cli.Timeline().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "updated" {
			t.Fatal("status must be updated")
		}
	}

	{
		i := &timeline.DeleteI{
			Obj: &timeline.DeleteI_Obj{
				Metadata: map[string]string{
					"timeline.venturemark.co/id":     ti2,
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         ui2,
				},
			},
		}

		o, err := cli.Timeline().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &audience.DeleteI{
			Obj: &audience.DeleteI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     ai1,
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         ui1,
				},
			},
		}

		o, err := cli.Audience().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["audience.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &audience.DeleteI{
			Obj: &audience.DeleteI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     ai2,
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         ui2,
				},
			},
		}

		o, err := cli.Audience().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["audience.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}
}

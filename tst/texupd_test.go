// +build conformance

package tst

import (
	"context"
	"testing"

	"github.com/venturemark/apigengo/pkg/pbf/texupd"
	"github.com/venturemark/apigengo/pkg/pbf/timeline"
	"github.com/venturemark/apigengo/pkg/pbf/update"

	"github.com/venturemark/fmz/pkg/client"
)

func Test_TexUpd_Lifecycle(t *testing.T) {
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
					"audience.venturemark.co/id":     "1",
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
			t.Fatal("timeline ID must not be empty")
		}

		tid = s
	}

	var ui1 string
	{
		i := &texupd.CreateI{
			Obj: &texupd.CreateI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     "1",
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     tid,
					"user.venturemark.co/id":         "1",
				},
				Property: &texupd.CreateI_Obj_Property{
					Text: "Lorem ipsum 1",
				},
			},
		}

		o, err := cli.TexUpd().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["update.venturemark.co/id"]
		if !ok {
			t.Fatal("texupd ID must not be empty")
		}

		ui1 = s
	}

	var ui2 string
	{
		i := &texupd.CreateI{
			Obj: &texupd.CreateI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     "1",
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     tid,
					"user.venturemark.co/id":         "1",
				},
				Property: &texupd.CreateI_Obj_Property{
					Text: "Lorem ipsum 2",
				},
			},
		}

		o, err := cli.TexUpd().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["update.venturemark.co/id"]
		if !ok {
			t.Fatal("texupd ID must not be empty")
		}

		ui2 = s
	}

	{
		i := &texupd.DeleteI{
			Obj: &texupd.DeleteI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     "1",
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     tid,
					"update.venturemark.co/id":       ui1,
					"user.venturemark.co/id":         "1",
				},
			},
		}

		o, err := cli.TexUpd().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["update.venturemark.co/status"]
		if !ok {
			t.Fatal("texupd status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("texupd status must be deleted")
		}
	}

	{
		i := &texupd.DeleteI{
			Obj: &texupd.DeleteI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     "1",
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     tid,
					"update.venturemark.co/id":       ui2,
					"user.venturemark.co/id":         "1",
				},
			},
		}

		o, err := cli.TexUpd().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["update.venturemark.co/status"]
		if !ok {
			t.Fatal("texupd status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("texupd status must be deleted")
		}
	}

	{
		i := &texupd.CreateI{
			Obj: &texupd.CreateI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     "1",
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     tid,
					"user.venturemark.co/id":         "1",
				},
				Property: &texupd.CreateI_Obj_Property{
					Text: "Lorem ipsum 1",
				},
			},
		}

		o, err := cli.TexUpd().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["update.venturemark.co/id"]
		if !ok {
			t.Fatal("texupd ID must not be empty")
		}

		ui1 = s
	}

	{
		i := &texupd.CreateI{
			Obj: &texupd.CreateI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     "1",
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     tid,
					"user.venturemark.co/id":         "1",
				},
				Property: &texupd.CreateI_Obj_Property{
					Text: "Lorem ipsum 2",
				},
			},
		}

		o, err := cli.TexUpd().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["update.venturemark.co/id"]
		if !ok {
			t.Fatal("texupd ID must not be empty")
		}

		ui2 = s
	}

	{
		i := &update.SearchI{
			Obj: []*update.SearchI_Obj{
				{
					Metadata: map[string]string{
						"audience.venturemark.co/id":     "1",
						"organization.venturemark.co/id": "1",
						"timeline.venturemark.co/id":     tid,
						"user.venturemark.co/id":         "1",
					},
				},
			},
		}

		o, err := cli.Update().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 2 {
			t.Fatal("there must be two updates")
		}

		// We show the latest updates first when watching at a timeline.
		if o.Obj[0].Property.Text != "Lorem ipsum 2" {
			t.Fatal("texupd name must be Lorem ipsum 1")
		}
		if o.Obj[1].Property.Text != "Lorem ipsum 1" {
			t.Fatal("texupd name must be Lorem ipsum 2")
		}
	}

	{
		i := &texupd.DeleteI{
			Obj: &texupd.DeleteI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     "1",
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     tid,
					"update.venturemark.co/id":       ui1,
					"user.venturemark.co/id":         "1",
				},
			},
		}

		o, err := cli.TexUpd().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["update.venturemark.co/status"]
		if !ok {
			t.Fatal("texupd status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("texupd status must be deleted")
		}
	}

	{
		i := &texupd.DeleteI{
			Obj: &texupd.DeleteI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     "1",
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     tid,
					"update.venturemark.co/id":       ui2,
					"user.venturemark.co/id":         "1",
				},
			},
		}

		o, err := cli.TexUpd().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["update.venturemark.co/status"]
		if !ok {
			t.Fatal("texupd status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("texupd status must be deleted")
		}
	}

	{
		i := &timeline.DeleteI{
			Obj: &timeline.DeleteI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     "1",
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
			t.Fatal("timeline status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("timeline status must be deleted")
		}
	}
}
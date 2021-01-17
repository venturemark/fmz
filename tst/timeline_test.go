// +build conformance

package tst

import (
	"context"
	"testing"

	"github.com/venturemark/apigengo/pkg/pbf/timeline"

	"github.com/venturemark/fmz/pkg/client"
)

func Test_Timeline_Lifecycle(t *testing.T) {
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

		_, err := cli.Timeline().Create(context.Background(), i)
		if err == nil {
			t.Fatal("timeline name must be unique")
		}
	}

	{
		i := &timeline.DeleteI{
			Obj: &timeline.DeleteI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     "1",
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     tid,
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

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"audience.venturemark.co/id":     "1",
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

		if o.Obj[0].Property.Name != "Marketing Campaign" {
			t.Fatal("timeline name must be Marketing Campaign")
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

// +build conformance

package tst

import (
	"context"
	"testing"

	"github.com/venturemark/apigengo/pkg/pbf/audience"

	"github.com/venturemark/fmz/pkg/client"
)

func Test_Audience_Lifecycle(t *testing.T) {
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

	var aid string
	{
		i := &audience.CreateI{
			Obj: &audience.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
				Property: &audience.CreateI_Obj_Property{
					Name: "Employees",
					User: []string{
						"xh3b4sd",
						"marcoelli",
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

		aid = s
	}

	{
		i := &audience.CreateI{
			Obj: &audience.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
				Property: &audience.CreateI_Obj_Property{
					Name: "Employees",
					User: []string{
						"xh3b4sd",
						"marcoelli",
					},
				},
			},
		}

		_, err := cli.Audience().Create(context.Background(), i)
		if err == nil {
			t.Fatal("audience name must be unique")
		}
	}

	{
		i := &audience.DeleteI{
			Obj: &audience.DeleteI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     aid,
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
			},
		}

		o, err := cli.Audience().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["audience.venturemark.co/status"]
		if !ok {
			t.Fatal("audience status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("audience status must be deleted")
		}
	}

	{
		i := &audience.CreateI{
			Obj: &audience.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
				Property: &audience.CreateI_Obj_Property{
					Name: "Employees",
					User: []string{
						"xh3b4sd",
						"marcoelli",
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

		aid = s
	}

	{
		i := &audience.SearchI{
			Obj: []*audience.SearchI_Obj{
				{
					Metadata: map[string]string{
						"organization.venturemark.co/id": "1",
						"user.venturemark.co/id":         "1",
					},
				},
			},
		}

		o, err := cli.Audience().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one audience")
		}

		if o.Obj[0].Property.Name != "Employees" {
			t.Fatal("audience name must be Employees")
		}
	}

	{
		i := &audience.DeleteI{
			Obj: &audience.DeleteI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     aid,
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
			},
		}

		o, err := cli.Audience().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["audience.venturemark.co/status"]
		if !ok {
			t.Fatal("audience status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("audience status must be deleted")
		}
	}
}

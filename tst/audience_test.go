// +build conformance

package tst

import (
	"context"
	"testing"

	"github.com/venturemark/apigengo/pkg/pbf/audience"

	"github.com/venturemark/fmz/pkg/client"
)

// Test_Audience_001 ensures that the lifecycle of audiences is covered from
// creation to deletion.
func Test_Audience_001(t *testing.T) {
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

	var ai1 string
	{
		i := &audience.CreateI{
			Obj: &audience.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
				Property: &audience.CreateI_Obj_Property{
					Name: "Employees",
					Tmln: []string{
						"foo",
						"bar",
					},
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

		ai1 = s
	}

	var ai2 string
	{
		i := &audience.CreateI{
			Obj: &audience.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"user.venturemark.co/id":         "1",
				},
				Property: &audience.CreateI_Obj_Property{
					Name: "Vendors",
					Tmln: []string{
						"foo",
						"bar",
					},
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

		ai2 = s
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

		if len(o.Obj) != 2 {
			t.Fatal("there must be two audiences")
		}

		if o.Obj[0].Property.Name != "Vendors" {
			t.Fatal("audience name must be Vendors")
		}
		if o.Obj[1].Property.Name != "Employees" {
			t.Fatal("audience name must be Employees")
		}
	}

	{
		i := &audience.DeleteI{
			Obj: &audience.DeleteI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     ai1,
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
		i := &audience.DeleteI{
			Obj: &audience.DeleteI_Obj{
				Metadata: map[string]string{
					"audience.venturemark.co/id":     ai2,
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

		if len(o.Obj) != 0 {
			t.Fatal("there must be zero audiences")
		}
	}
}

// Test_Audience_002 ensures that audience names are unique.
func Test_Audience_002(t *testing.T) {
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
					Tmln: []string{
						"foo",
						"bar",
					},
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
					Tmln: []string{
						"foo",
						"bar",
					},
					User: []string{
						"foo",
						"bar",
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
}

// Test_Audience_003 ensures that audiences cannot be created
// without timelines.
func Test_Audience_003(t *testing.T) {
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
			t.Fatal("timelines must not be empty")
		}
	}
}

// Test_Audience_004 is a temporary test that ensures audiences can be created
// with zero users.
func Test_Audience_004(t *testing.T) {
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
					Tmln: []string{
						"foo",
						"bar",
					},
					User: []string{},
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

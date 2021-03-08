// +build conformance

package tst

import (
	"context"
	"testing"

	"github.com/venturemark/apigengo/pkg/pbf/audience"

	"github.com/venturemark/cfm/pkg/client"
	"github.com/venturemark/cfm/pkg/to"
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

		err = cli.Redigo().Purge()
		if err != nil {
			t.Fatal(err)
		}

		defer cli.Grpc().Close()
	}

	var ai1 string
	{
		i := &audience.CreateI{
			Obj: []*audience.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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
			},
		}

		o, err := cli.Audience().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ai1 = s
	}

	var ai2 string
	{
		i := &audience.CreateI{
			Obj: []*audience.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
					Property: &audience.CreateI_Obj_Property{
						Name: "Investors",
						Tmln: []string{
							"bar",
							"baz",
						},
						User: []string{
							"marcoelli",
							"xh3b4sd",
						},
					},
				},
			},
		}

		o, err := cli.Audience().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ai2 = s
	}

	{
		i := &audience.SearchI{
			Obj: []*audience.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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

		{
			s, ok := o.Obj[0].Metadata["audience.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ai2 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[0].Property.Name != "Investors" {
				t.Fatal("name must be Investors")
			}
			if o.Obj[0].Property.Tmln[0] != "bar" {
				t.Fatal("timeline must include bar")
			}
			if o.Obj[0].Property.Tmln[1] != "baz" {
				t.Fatal("timeline must include baz")
			}
			if o.Obj[0].Property.User[0] != "marcoelli" {
				t.Fatal("user must include marcoelli")
			}
			if o.Obj[0].Property.User[1] != "xh3b4sd" {
				t.Fatal("user must include xh3b4sd")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["audience.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ai1 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[1].Property.Name != "Employees" {
				t.Fatal("name must be Employees")
			}
			if o.Obj[1].Property.Tmln[0] != "foo" {
				t.Fatal("timeline must include foo")
			}
			if o.Obj[1].Property.Tmln[1] != "bar" {
				t.Fatal("timeline must include bar")
			}
			if o.Obj[1].Property.User[0] != "xh3b4sd" {
				t.Fatal("user must include xh3b4sd")
			}
			if o.Obj[1].Property.User[1] != "marcoelli" {
				t.Fatal("user must include marcoelli")
			}
		}
	}

	{
		i := &audience.DeleteI{
			Obj: []*audience.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"audience.venturemark.co/id": ai1,
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cli.Audience().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &audience.DeleteI{
			Obj: []*audience.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"audience.venturemark.co/id": ai2,
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cli.Audience().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &audience.SearchI{
			Obj: []*audience.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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

		err = cli.Redigo().Purge()
		if err != nil {
			t.Fatal(err)
		}

		defer cli.Grpc().Close()
	}

	var aid string
	{
		i := &audience.CreateI{
			Obj: []*audience.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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
			},
		}

		o, err := cli.Audience().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		aid = s
	}

	{
		i := &audience.CreateI{
			Obj: []*audience.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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
			},
		}

		_, err := cli.Audience().Create(context.Background(), i)
		if err == nil {
			t.Fatal("name must be unique")
		}
	}

	{
		i := &audience.DeleteI{
			Obj: []*audience.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"audience.venturemark.co/id": aid,
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cli.Audience().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
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

		err = cli.Redigo().Purge()
		if err != nil {
			t.Fatal(err)
		}

		defer cli.Grpc().Close()
	}

	{
		i := &audience.CreateI{
			Obj: []*audience.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
					Property: &audience.CreateI_Obj_Property{
						Name: "Employees",
						User: []string{
							"xh3b4sd",
							"marcoelli",
						},
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

		err = cli.Redigo().Purge()
		if err != nil {
			t.Fatal(err)
		}

		defer cli.Grpc().Close()
	}

	var aid string
	{
		i := &audience.CreateI{
			Obj: []*audience.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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
			},
		}

		o, err := cli.Audience().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		aid = s
	}

	{
		i := &audience.DeleteI{
			Obj: []*audience.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"audience.venturemark.co/id": aid,
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cli.Audience().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}
}

// Test_Audience_005 ensures that updating audiences via JSON-Patch methods
// works as expected.
func Test_Audience_005(t *testing.T) {
	var err error

	var cli *client.Client
	{
		c := client.Config{}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Redigo().Purge()
		if err != nil {
			t.Fatal(err)
		}

		defer cli.Grpc().Close()
	}

	var ai1 string
	{
		i := &audience.CreateI{
			Obj: []*audience.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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
			},
		}

		o, err := cli.Audience().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ai1 = s
	}

	var ai2 string
	{
		i := &audience.CreateI{
			Obj: []*audience.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
					Property: &audience.CreateI_Obj_Property{
						Name: "Investors",
						Tmln: []string{
							"bar",
							"baz",
						},
						User: []string{
							"marcoelli",
							"xh3b4sd",
						},
					},
				},
			},
		}

		o, err := cli.Audience().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ai2 = s
	}

	{
		i := &audience.SearchI{
			Obj: []*audience.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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

		if o.Obj[0].Property.Name != "Investors" {
			t.Fatal("name must be Investors")
		}
		if o.Obj[0].Property.Tmln[0] != "bar" {
			t.Fatal("timeline must include bar")
		}
		if o.Obj[0].Property.Tmln[1] != "baz" {
			t.Fatal("timeline must include baz")
		}
		if o.Obj[0].Property.User[0] != "marcoelli" {
			t.Fatal("user must include marcoelli")
		}
		if o.Obj[0].Property.User[1] != "xh3b4sd" {
			t.Fatal("user must include xh3b4sd")
		}
		if o.Obj[1].Property.Name != "Employees" {
			t.Fatal("name must be Employees")
		}
		if o.Obj[1].Property.Tmln[0] != "foo" {
			t.Fatal("timeline must include foo")
		}
		if o.Obj[1].Property.Tmln[1] != "bar" {
			t.Fatal("timeline must include bar")
		}
		if o.Obj[1].Property.User[0] != "xh3b4sd" {
			t.Fatal("user must include xh3b4sd")
		}
		if o.Obj[1].Property.User[1] != "marcoelli" {
			t.Fatal("user must include marcoelli")
		}
	}

	{
		i := &audience.UpdateI{
			Obj: []*audience.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"audience.venturemark.co/id": ai1,
						"venture.venturemark.co/id":  "1",
					},
					Jsnpatch: []*audience.UpdateI_Obj_Jsnpatch{
						{
							Ope: "replace",
							Pat: "/obj/property/name",
							Val: to.StringP("replaced name"),
						},
						{
							Ope: "remove",
							Pat: "/obj/property/user/0",
						},
						{
							Ope: "add",
							Pat: "/obj/property/user/-",
							Val: to.StringP("added user"),
						},
					},
				},
			},
		}

		o, err := cli.Audience().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "updated" {
			t.Fatal("status must be updated")
		}
	}

	{
		i := &audience.SearchI{
			Obj: []*audience.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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

		if o.Obj[0].Property.Name != "Investors" {
			t.Fatal("name must be Investors")
		}
		if o.Obj[0].Property.Tmln[0] != "bar" {
			t.Fatal("timeline must include bar")
		}
		if o.Obj[0].Property.Tmln[1] != "baz" {
			t.Fatal("timeline must include baz")
		}
		if o.Obj[0].Property.User[0] != "marcoelli" {
			t.Fatal("user must include marcoelli")
		}
		if o.Obj[0].Property.User[1] != "xh3b4sd" {
			t.Fatal("user must include xh3b4sd")
		}
		if o.Obj[1].Property.Name != "replaced name" {
			t.Fatal("name must be replaced name")
		}
		if o.Obj[1].Property.Tmln[0] != "foo" {
			t.Fatal("timeline must include foo")
		}
		if o.Obj[1].Property.Tmln[1] != "bar" {
			t.Fatal("timeline must include bar")
		}
		if o.Obj[1].Property.User[0] != "marcoelli" {
			t.Fatal("user must include marcoelli")
		}
		if o.Obj[1].Property.User[1] != "added user" {
			t.Fatal("user must include added user")
		}
	}

	{
		i := &audience.DeleteI{
			Obj: []*audience.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"audience.venturemark.co/id": ai1,
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cli.Audience().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &audience.DeleteI{
			Obj: []*audience.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"audience.venturemark.co/id": ai2,
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cli.Audience().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &audience.SearchI{
			Obj: []*audience.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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

// Test_Audience_006 ensures that deleting audience resources which do not exist
// returns an error.
func Test_Audience_006(t *testing.T) {
	var err error

	var cli *client.Client
	{
		c := client.Config{}

		cli, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cli.Redigo().Purge()
		if err != nil {
			t.Fatal(err)
		}

		defer cli.Grpc().Close()
	}

	{
		i := &audience.DeleteI{
			Obj: []*audience.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"audience.venturemark.co/id": "1",
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		_, err := cli.Audience().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}
}

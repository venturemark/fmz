// +build conformance

package tst

import (
	"context"
	"testing"

	"github.com/venturemark/apigengo/pkg/pbf/role"

	"github.com/venturemark/cfm/pkg/client"
	"github.com/venturemark/cfm/pkg/to"
)

// Test_Role_001 ensures that the lifecycle of roles is covered from
// creation to deletion.
func Test_Role_001(t *testing.T) {
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

	var ri1 string
	{
		i := &role.CreateI{
			Obj: []*role.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
						"subject.venturemark.co/id": "1",
					},
					Property: &role.CreateI_Obj_Property{
						Kin: "owner",
						Res: "venture",
					},
				},
			},
		}

		o, err := cli.Role().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one role")
		}

		s, ok := o.Obj[0].Metadata["role.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ri1 = s
	}

	var ri2 string
	{
		i := &role.CreateI{
			Obj: []*role.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
						"subject.venturemark.co/id": "2",
					},
					Property: &role.CreateI_Obj_Property{
						Kin: "member",
						Res: "venture",
					},
				},
			},
		}

		o, err := cli.Role().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one role")
		}

		s, ok := o.Obj[0].Metadata["role.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ri2 = s
	}

	{
		i := &role.SearchI{
			Obj: []*role.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
				},
			},
		}

		o, err := cli.Role().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 2 {
			t.Fatal("there must be two roles")
		}

		{
			s, ok := o.Obj[0].Metadata["role.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ri2 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[0].Property.Kin != "member" {
				t.Fatal("kin must be member")
			}
			if o.Obj[0].Property.Res != "venture" {
				t.Fatal("res must be venture")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["role.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ri1 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[1].Property.Kin != "owner" {
				t.Fatal("kin must be owner")
			}
			if o.Obj[1].Property.Res != "venture" {
				t.Fatal("res must be venture")
			}
		}
	}

	{
		i := &role.UpdateI{
			Obj: []*role.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"role.venturemark.co/id":    ri2,
						"venture.venturemark.co/id": "1",
					},
					Jsnpatch: []*role.UpdateI_Obj_Jsnpatch{
						{
							Ope: "replace",
							Pat: "/obj/property/kin",
							Val: to.StringP("owner"),
						},
					},
				},
			},
		}

		o, err := cli.Role().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one role")
		}

		s, ok := o.Obj[0].Metadata["role.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "updated" {
			t.Fatal("status must be updated")
		}
	}

	{
		i := &role.SearchI{
			Obj: []*role.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
				},
			},
		}

		o, err := cli.Role().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 2 {
			t.Fatal("there must be two roles")
		}

		{
			s, ok := o.Obj[0].Metadata["role.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ri2 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[0].Property.Kin != "owner" {
				t.Fatal("kin must be owner")
			}
			if o.Obj[0].Property.Res != "venture" {
				t.Fatal("res must be venture")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["role.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ri1 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[1].Property.Kin != "owner" {
				t.Fatal("kin must be owner")
			}
			if o.Obj[1].Property.Res != "venture" {
				t.Fatal("res must be venture")
			}
		}
	}

	{
		i := &role.DeleteI{
			Obj: []*role.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"role.venturemark.co/id":    ri1,
						"venture.venturemark.co/id": "1",
					},
				},
			},
		}

		o, err := cli.Role().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one role")
		}

		s, ok := o.Obj[0].Metadata["role.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &role.DeleteI{
			Obj: []*role.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"role.venturemark.co/id":    ri2,
						"venture.venturemark.co/id": "1",
					},
				},
			},
		}

		o, err := cli.Role().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one role")
		}

		s, ok := o.Obj[0].Metadata["role.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &role.SearchI{
			Obj: []*role.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
				},
			},
		}

		o, err := cli.Role().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 0 {
			t.Fatal("there must be zero roles")
		}
	}
}

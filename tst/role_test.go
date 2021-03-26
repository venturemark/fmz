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

	var su1 string
	var su2 string
	{
		su1 = "1"
		su2 = "2"
	}

	var ro1 string
	{
		i := &role.CreateI{
			Obj: []*role.CreateI_Obj{
				{
					Metadata: map[string]string{
						"resource.venturemark.co/kind": "venture",
						"role.venturemark.co/kind":     "owner",
						"subject.venturemark.co/id":    su1,
						"venture.venturemark.co/id":    "1",
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

		ro1 = s
	}

	var ro2 string
	{
		i := &role.CreateI{
			Obj: []*role.CreateI_Obj{
				{
					Metadata: map[string]string{
						"resource.venturemark.co/kind": "venture",
						"role.venturemark.co/kind":     "member",
						"subject.venturemark.co/id":    su2,
						"venture.venturemark.co/id":    "1",
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

		ro2 = s
	}

	{
		i := &role.SearchI{
			Obj: []*role.SearchI_Obj{
				{
					Metadata: map[string]string{
						"resource.venturemark.co/kind": "venture",
						"venture.venturemark.co/id":    "1",
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
			s, ok := o.Obj[0].Metadata["resource.venturemark.co/kind"]
			if !ok {
				t.Fatal("kind must not be empty")
			}
			if s != "venture" {
				t.Fatal("kind must match across actions")
			}
		}

		{
			s, ok := o.Obj[0].Metadata["role.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ro2 {
				t.Fatal("id must match across actions")
			}
		}

		{
			s, ok := o.Obj[0].Metadata["role.venturemark.co/kind"]
			if !ok {
				t.Fatal("kind must not be empty")
			}
			if s != "member" {
				t.Fatal("kind must match across actions")
			}
		}

		{
			s, ok := o.Obj[0].Metadata["subject.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != su2 {
				t.Fatal("id must match across actions")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["resource.venturemark.co/kind"]
			if !ok {
				t.Fatal("kind must not be empty")
			}
			if s != "venture" {
				t.Fatal("kind must match across actions")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["role.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ro1 {
				t.Fatal("id must match across actions")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["role.venturemark.co/kind"]
			if !ok {
				t.Fatal("kind must not be empty")
			}
			if s != "owner" {
				t.Fatal("kind must match across actions")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["subject.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != su1 {
				t.Fatal("id must match across actions")
			}
		}
	}

	{
		i := &role.UpdateI{
			Obj: []*role.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"resource.venturemark.co/kind": "venture",
						"role.venturemark.co/id":       ro2,
						"venture.venturemark.co/id":    "1",
					},
					Jsnpatch: []*role.UpdateI_Obj_Jsnpatch{
						{
							Ope: "replace",
							Pat: "/obj/metadata/role.venturemark.co~1kind",
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

		{
			s, ok := o.Obj[0].Metadata["role.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}

			if s != ro2 {
				t.Fatal("id must match across actions")
			}
		}

		{
			s, ok := o.Obj[0].Metadata["role.venturemark.co/status"]
			if !ok {
				t.Fatal("status must not be empty")
			}

			if s != "updated" {
				t.Fatal("status must be updated")
			}
		}
	}

	{
		i := &role.SearchI{
			Obj: []*role.SearchI_Obj{
				{
					Metadata: map[string]string{
						"resource.venturemark.co/kind": "venture",
						"venture.venturemark.co/id":    "1",
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
			s, ok := o.Obj[0].Metadata["resource.venturemark.co/kind"]
			if !ok {
				t.Fatal("kind must not be empty")
			}
			if s != "venture" {
				t.Fatal("kind must match across actions")
			}
		}

		{
			s, ok := o.Obj[0].Metadata["role.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ro2 {
				t.Fatal("id must match across actions")
			}
		}

		{
			s, ok := o.Obj[0].Metadata["role.venturemark.co/kind"]
			if !ok {
				t.Fatal("kind must not be empty")
			}
			if s != "owner" {
				t.Fatal("kind must match across actions")
			}
		}

		{
			s, ok := o.Obj[0].Metadata["subject.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != su2 {
				t.Fatal("id must match across actions")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["resource.venturemark.co/kind"]
			if !ok {
				t.Fatal("kind must not be empty")
			}
			if s != "venture" {
				t.Fatal("kind must match across actions")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["role.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ro1 {
				t.Fatal("id must match across actions")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["role.venturemark.co/kind"]
			if !ok {
				t.Fatal("kind must not be empty")
			}
			if s != "owner" {
				t.Fatal("kind must match across actions")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["subject.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != su1 {
				t.Fatal("id must match across actions")
			}
		}
	}

	{
		i := &role.DeleteI{
			Obj: []*role.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"resource.venturemark.co/kind": "venture",
						"role.venturemark.co/id":       ro1,
						"venture.venturemark.co/id":    "1",
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
						"resource.venturemark.co/kind": "venture",
						"role.venturemark.co/id":       ro2,
						"venture.venturemark.co/id":    "1",
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
						"resource.venturemark.co/kind": "venture",
						"venture.venturemark.co/id":    "1",
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

// Test_Role_002 ensures that deleting role resources which do not exist returns
// an error.
func Test_Role_002(t *testing.T) {
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
		i := &role.DeleteI{
			Obj: []*role.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"resource.venturemark.co/kind": "venture",
						"role.venturemark.co/id":       "1",
						"venture.venturemark.co/id":    "1",
					},
				},
			},
		}

		_, err := cli.Role().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}
}

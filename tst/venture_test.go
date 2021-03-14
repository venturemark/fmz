// +build conformance

package tst

import (
	"context"
	"testing"

	"github.com/venturemark/apigengo/pkg/pbf/role"
	"github.com/venturemark/apigengo/pkg/pbf/venture"

	"github.com/venturemark/cfm/pkg/client"
	"github.com/venturemark/cfm/pkg/oauth"
)

// Test_Venture_001 ensures that the lifecycle of ventures is covered from
// creation to deletion.
func Test_Venture_001(t *testing.T) {
	var err error

	var cr1 *oauth.Insecure
	var cr2 *oauth.Insecure
	{
		cr1 = oauth.NewInsecureOne()
		cr2 = oauth.NewInsecureTwo()
	}

	var cl1 *client.Client
	{
		c := client.Config{
			Credentials: cr1,
		}

		cl1, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cl1.Redigo().Purge()
		if err != nil {
			t.Fatal(err)
		}

		defer cl1.Grpc().Close()
	}

	var cl2 *client.Client
	{
		c := client.Config{
			Credentials: cr2,
		}

		cl2, err = client.New(c)
		if err != nil {
			t.Fatal(err)
		}

		err = cl2.Redigo().Purge()
		if err != nil {
			t.Fatal(err)
		}

		defer cl2.Grpc().Close()
	}

	var ve1 string
	{
		i := &venture.CreateI{
			Obj: []*venture.CreateI_Obj{
				{
					Property: &venture.CreateI_Obj_Property{
						Name: "IBM",
					},
				},
			},
		}

		o, err := cl1.Venture().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["venture.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ve1 = s
	}

	var ve2 string
	{
		i := &venture.CreateI{
			Obj: []*venture.CreateI_Obj{
				{
					Property: &venture.CreateI_Obj_Property{
						Name: "GME",
					},
				},
			},
		}

		o, err := cl1.Venture().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["venture.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ve2 = s
	}

	{
		i := &role.CreateI{
			Obj: []*role.CreateI_Obj{
				{
					Metadata: map[string]string{
						"resource.venturemark.co/kind": "venture",
						"role.venturemark.co/kind":     "member",
						"subject.venturemark.co/id":    cr2.User(),
						"venture.venturemark.co/id":    ve2,
					},
				},
			},
		}

		o, err := cl1.Role().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one role")
		}

		_, ok := o.Obj[0].Metadata["role.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}
	}

	{
		i := &venture.SearchI{
			Obj: []*venture.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": ve1,
					},
				},
			},
		}

		o, err := cl1.Venture().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one venture")
		}

		{
			s, ok := o.Obj[0].Metadata["venture.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ve1 {
				t.Fatal("id must match across actions")
			}
		}

		{
			if o.Obj[0].Property.Name != "IBM" {
				t.Fatal("name must be IBM")
			}
		}
	}

	{
		i := &venture.SearchI{
			Obj: []*venture.SearchI_Obj{
				{
					Metadata: map[string]string{
						"subject.venturemark.co/id": cr1.User(),
					},
				},
			},
		}

		o, err := cl1.Venture().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 2 {
			t.Fatal("there must be two ventures")
		}

		{
			s, ok := o.Obj[0].Metadata["venture.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ve2 {
				t.Fatal("id must match across actions")
			}
		}

		{
			if o.Obj[0].Property.Name != "GME" {
				t.Fatal("name must be GME")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["venture.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ve1 {
				t.Fatal("id must match across actions")
			}
		}

		{
			if o.Obj[1].Property.Name != "IBM" {
				t.Fatal("name must be IBM")
			}
		}
	}

	{
		i := &venture.SearchI{
			Obj: []*venture.SearchI_Obj{
				{
					Metadata: map[string]string{
						"subject.venturemark.co/id": cr2.User(),
					},
				},
			},
		}

		o, err := cl1.Venture().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one venture")
		}

		{
			s, ok := o.Obj[0].Metadata["venture.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ve2 {
				t.Fatal("id must match across actions")
			}
		}

		{
			if o.Obj[0].Property.Name != "GME" {
				t.Fatal("name must be GME")
			}
		}
	}

	{
		i := &venture.DeleteI{
			Obj: []*venture.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": ve1,
					},
				},
			},
		}

		_, err := cl2.Venture().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}

	{
		i := &venture.DeleteI{
			Obj: []*venture.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": ve2,
					},
				},
			},
		}

		_, err := cl2.Venture().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}

	{
		i := &venture.DeleteI{
			Obj: []*venture.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": ve1,
					},
				},
			},
		}

		o, err := cl1.Venture().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["venture.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &venture.DeleteI{
			Obj: []*venture.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": ve2,
					},
				},
			},
		}

		o, err := cl1.Venture().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["venture.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &venture.SearchI{
			Obj: []*venture.SearchI_Obj{
				{
					Metadata: map[string]string{
						"subject.venturemark.co/id": cr1.User(),
					},
				},
			},
		}

		o, err := cl1.Venture().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 0 {
			t.Fatal("there must be zero ventures")
		}
	}

	{
		i := &venture.SearchI{
			Obj: []*venture.SearchI_Obj{
				{
					Metadata: map[string]string{
						"subject.venturemark.co/id": cr2.User(),
					},
				},
			},
		}

		o, err := cl1.Venture().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 0 {
			t.Fatal("there must be zero ventures")
		}
	}
}

// Test_Venture_002 ensures that deleting venture resources which do not exist
// returns an error.
func Test_Venture_002(t *testing.T) {
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
		i := &venture.DeleteI{
			Obj: []*venture.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
				},
			},
		}

		_, err := cli.Venture().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}
}

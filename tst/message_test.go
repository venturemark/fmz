// +build conformance

package tst

import (
	"context"
	"testing"

	"github.com/venturemark/apigengo/pkg/pbf/message"
	"github.com/venturemark/apigengo/pkg/pbf/role"
	"github.com/venturemark/apigengo/pkg/pbf/timeline"
	"github.com/venturemark/apigengo/pkg/pbf/user"
	"github.com/venturemark/apigengo/pkg/pbf/venture"

	"github.com/venturemark/cfm/pkg/client"
	"github.com/venturemark/cfm/pkg/oauth"
)

// Test_Message_001 ensures that the lifecycle of messages is covered from
// creation to deletion.
func Test_Message_001(t *testing.T) {
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

	var us1 string
	{
		i := &user.CreateI{
			Obj: []*user.CreateI_Obj{
				{
					Property: &user.CreateI_Obj_Property{
						Name: "marcojelli",
					},
				},
			},
		}

		o, err := cl1.User().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		us1 = s
	}

	var us2 string
	{
		i := &user.CreateI{
			Obj: []*user.CreateI_Obj{
				{
					Property: &user.CreateI_Obj_Property{
						Name: "disreszi",
					},
				},
			},
		}

		o, err := cl2.User().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		us2 = s
	}

	var vei string
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

		vei = s
	}

	var tii string
	{
		i := &timeline.CreateI{
			Obj: []*timeline.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
					Property: &timeline.CreateI_Obj_Property{
						Name: "Marketing Campaign",
					},
				},
			},
		}

		o, err := cl1.Timeline().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["timeline.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		tii = s
	}

	{
		i := &role.CreateI{
			Obj: []*role.CreateI_Obj{
				{
					Metadata: map[string]string{
						"resource.venturemark.co/kind": "timeline",
						"role.venturemark.co/kind":     "member",
						"subject.venturemark.co/id":    us2,
						"timeline.venturemark.co/id":   tii,
						"venture.venturemark.co/id":    vei,
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

	var mi1 string
	{
		i := &message.CreateI{
			Obj: []*message.CreateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  vei,
					},
					Property: &message.CreateI_Obj_Property{
						Text: "Lorem ipsum 1",
					},
				},
			},
		}

		o, err := cl1.Message().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["message.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		mi1 = s
	}

	var mi2 string
	{
		i := &message.CreateI{
			Obj: []*message.CreateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  vei,
					},
					Property: &message.CreateI_Obj_Property{
						Text: "Lorem ipsum 2",
					},
				},
			},
		}

		o, err := cl2.Message().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["message.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		mi2 = s
	}

	{
		i := &message.SearchI{
			Obj: []*message.SearchI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  vei,
					},
				},
			},
		}

		o, err := cl1.Message().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 2 {
			t.Fatal("there must be two messages")
		}

		{
			s, ok := o.Obj[0].Metadata["message.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != mi2 {
				t.Fatal("id must match across actions")
			}
		}

		{
			s, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != us2 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[0].Property.Text != "Lorem ipsum 2" {
				t.Fatal("text must be Lorem ipsum 1")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["message.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != mi1 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[1].Property.Text != "Lorem ipsum 1" {
				t.Fatal("text must be Lorem ipsum 2")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != us1 {
				t.Fatal("id must match across actions")
			}
		}
	}

	{
		i := &message.SearchI{
			Obj: []*message.SearchI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  vei,
					},
				},
			},
		}

		o, err := cl2.Message().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 2 {
			t.Fatal("there must be two messages")
		}

		{
			if o.Obj[0].Property.Text != "Lorem ipsum 2" {
				t.Fatal("text must be Lorem ipsum 1")
			}
			s, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != us2 {
				t.Fatal("id must match")
			}
		}

		{
			if o.Obj[1].Property.Text != "Lorem ipsum 1" {
				t.Fatal("text must be Lorem ipsum 2")
			}
			s, ok := o.Obj[1].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != us1 {
				t.Fatal("id must match")
			}
		}
	}

	{
		i := &message.DeleteI{
			Obj: []*message.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"message.venturemark.co/id":  mi1,
						"timeline.venturemark.co/id": tii,
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  vei,
					},
				},
			},
		}

		o, err := cl1.Message().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["message.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &message.DeleteI{
			Obj: []*message.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"message.venturemark.co/id":  mi2,
						"timeline.venturemark.co/id": tii,
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  vei,
					},
				},
			},
		}

		o, err := cl2.Message().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["message.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &message.SearchI{
			Obj: []*message.SearchI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  vei,
					},
				},
			},
		}

		o, err := cl1.Message().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 0 {
			t.Fatal("there must be zero messages")
		}
	}

	{
		i := &message.SearchI{
			Obj: []*message.SearchI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  vei,
					},
				},
			},
		}

		o, err := cl2.Message().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 0 {
			t.Fatal("there must be zero messages")
		}
	}
}

// Test_Message_002 ensures that deleting message resources which do not exist
// returns an error.
func Test_Message_002(t *testing.T) {
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
		i := &message.DeleteI{
			Obj: []*message.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"message.venturemark.co/id":  "1",
						"timeline.venturemark.co/id": "1",
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		_, err := cli.Message().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}
}

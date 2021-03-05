// +build conformance

package tst

import (
	"context"
	"testing"

	"github.com/venturemark/apigengo/pkg/pbf/message"

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

	var mi1 string
	{
		i := &message.CreateI{
			Obj: &message.CreateI_Obj{
				Metadata: map[string]string{
					"timeline.venturemark.co/id": "1",
					"update.venturemark.co/id":   "1",
					"venture.venturemark.co/id":  "1",
				},
				Property: &message.CreateI_Obj_Property{
					Text: "Lorem ipsum 1",
				},
			},
		}

		o, err := cl1.Message().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["message.venturemark.co/id"]
		if !ok {
			t.Fatal("message ID must not be empty")
		}

		mi1 = s
	}

	var mi2 string
	{
		i := &message.CreateI{
			Obj: &message.CreateI_Obj{
				Metadata: map[string]string{
					"timeline.venturemark.co/id": "1",
					"update.venturemark.co/id":   "1",
					"venture.venturemark.co/id":  "1",
				},
				Property: &message.CreateI_Obj_Property{
					Text: "Lorem ipsum 2",
				},
			},
		}

		o, err := cl2.Message().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["message.venturemark.co/id"]
		if !ok {
			t.Fatal("message ID must not be empty")
		}

		mi2 = s
	}

	{
		i := &message.SearchI{
			Obj: []*message.SearchI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": "1",
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  "1",
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
			if s != cr2.User() {
				t.Fatal("id must match across actions")
			}
			if o.Obj[0].Property.Text != "Lorem ipsum 2" {
				t.Fatal("message text must be Lorem ipsum 1")
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
				t.Fatal("message text must be Lorem ipsum 2")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != cr1.User() {
				t.Fatal("id must match across actions")
			}
		}
	}

	{
		i := &message.SearchI{
			Obj: []*message.SearchI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": "1",
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  "1",
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
				t.Fatal("message text must be Lorem ipsum 1")
			}
			s, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != cr2.User() {
				t.Fatal("id must match")
			}
		}

		{
			if o.Obj[1].Property.Text != "Lorem ipsum 1" {
				t.Fatal("message text must be Lorem ipsum 2")
			}
			s, ok := o.Obj[1].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != cr1.User() {
				t.Fatal("id must match")
			}
		}
	}

	{
		i := &message.DeleteI{
			Obj: &message.DeleteI_Obj{
				Metadata: map[string]string{
					"message.venturemark.co/id":  mi1,
					"timeline.venturemark.co/id": "1",
					"update.venturemark.co/id":   "1",
					"venture.venturemark.co/id":  "1",
				},
			},
		}

		o, err := cl1.Message().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["message.venturemark.co/status"]
		if !ok {
			t.Fatal("message status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("message status must be deleted")
		}
	}

	{
		i := &message.DeleteI{
			Obj: &message.DeleteI_Obj{
				Metadata: map[string]string{
					"message.venturemark.co/id":  mi2,
					"timeline.venturemark.co/id": "1",
					"update.venturemark.co/id":   "1",
					"venture.venturemark.co/id":  "1",
				},
			},
		}

		o, err := cl2.Message().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj.Metadata["message.venturemark.co/status"]
		if !ok {
			t.Fatal("message status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("message status must be deleted")
		}
	}

	{
		i := &message.SearchI{
			Obj: []*message.SearchI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": "1",
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  "1",
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
						"timeline.venturemark.co/id": "1",
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  "1",
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

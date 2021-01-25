// +build conformance

package tst

import (
	"context"
	"testing"

	"github.com/venturemark/apigengo/pkg/pbf/message"

	"github.com/venturemark/fmz/pkg/client"
)

func Test_Message_Lifecycle(t *testing.T) {
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

	var mi1 string
	{
		i := &message.CreateI{
			Obj: &message.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     "1",
					"update.venturemark.co/id":       "1",
					"user.venturemark.co/id":         "1",
				},
				Property: &message.CreateI_Obj_Property{
					Text: "Lorem ipsum 1",
				},
			},
		}

		o, err := cli.Message().Create(context.Background(), i)
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
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     "1",
					"update.venturemark.co/id":       "1",
					"user.venturemark.co/id":         "1",
				},
				Property: &message.CreateI_Obj_Property{
					Text: "Lorem ipsum 2",
				},
			},
		}

		o, err := cli.Message().Create(context.Background(), i)
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
						"organization.venturemark.co/id": "1",
						"timeline.venturemark.co/id":     "1",
						"update.venturemark.co/id":       "1",
						"user.venturemark.co/id":         "1",
					},
				},
			},
		}

		o, err := cli.Message().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 2 {
			t.Fatal("there must be two messages")
		}

		if o.Obj[0].Property.Text != "Lorem ipsum 2" {
			t.Fatal("message text must be Lorem ipsum 1")
		}
		if o.Obj[1].Property.Text != "Lorem ipsum 1" {
			t.Fatal("message text must be Lorem ipsum 2")
		}
	}

	{
		i := &message.DeleteI{
			Obj: &message.DeleteI_Obj{
				Metadata: map[string]string{
					"message.venturemark.co/id":      mi1,
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     "1",
					"update.venturemark.co/id":       "1",
					"user.venturemark.co/id":         "1",
				},
			},
		}

		o, err := cli.Message().Delete(context.Background(), i)
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
					"message.venturemark.co/id":      mi2,
					"organization.venturemark.co/id": "1",
					"timeline.venturemark.co/id":     "1",
					"update.venturemark.co/id":       "1",
					"user.venturemark.co/id":         "1",
				},
			},
		}

		o, err := cli.Message().Delete(context.Background(), i)
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
						"organization.venturemark.co/id": "1",
						"timeline.venturemark.co/id":     "1",
						"update.venturemark.co/id":       "1",
						"user.venturemark.co/id":         "1",
					},
				},
			},
		}

		o, err := cli.Message().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 0 {
			t.Fatal("there must be zero messages")
		}
	}
}

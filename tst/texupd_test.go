// +build conformance

package tst

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/venturemark/apigengo/pkg/pbf/message"
	"github.com/venturemark/apigengo/pkg/pbf/texupd"
	"github.com/venturemark/apigengo/pkg/pbf/timeline"
	"github.com/venturemark/apigengo/pkg/pbf/update"
	"github.com/venturemark/apigengo/pkg/pbf/venture"
	"github.com/xh3b4sd/budget"
	"github.com/xh3b4sd/tracer"

	"github.com/venturemark/cfm/pkg/client"
	"github.com/venturemark/cfm/pkg/to"
)

// Test_TexUpd_001 ensures that the lifecycle of text updates is covered from
// creation to deletion.
func Test_TexUpd_001(t *testing.T) {
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

		o, err := cli.Venture().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["venture.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ve1 = s
	}

	var tii string
	{
		i := &timeline.CreateI{
			Obj: []*timeline.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": ve1,
					},
					Property: &timeline.CreateI_Obj_Property{
						Name: "Marketing Campaign",
					},
				},
			},
		}

		o, err := cli.Timeline().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["timeline.venturemark.co/id"]
		if !ok {
			t.Fatal("timeline ID must not be empty")
		}

		tii = s
	}

	var up1 string
	{
		i := &texupd.CreateI{
			Obj: []*texupd.CreateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"venture.venturemark.co/id":  ve1,
					},
					Property: &texupd.CreateI_Obj_Property{
						Text: "Lorem ipsum 1",
					},
				},
			},
		}

		o, err := cli.TexUpd().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["update.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		up1 = s
	}

	var up2 string
	{
		i := &texupd.CreateI{
			Obj: []*texupd.CreateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"venture.venturemark.co/id":  ve1,
					},
					Property: &texupd.CreateI_Obj_Property{
						Text: "Lorem ipsum 2",
					},
				},
			},
		}

		o, err := cli.TexUpd().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["update.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		up2 = s
	}

	{
		i := &update.SearchI{
			Obj: []*update.SearchI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"venture.venturemark.co/id":  ve1,
					},
				},
			},
		}

		o, err := cli.Update().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 2 {
			t.Fatal("there must be two updates")
		}

		{
			s, ok := o.Obj[0].Metadata["update.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != up2 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[0].Property.Text != "Lorem ipsum 2" {
				t.Fatal("text must be Lorem ipsum 2")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["update.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != up1 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[1].Property.Text != "Lorem ipsum 1" {
				t.Fatal("text must be Lorem ipsum 1")
			}
		}
	}

	{
		i := &texupd.DeleteI{
			Obj: []*texupd.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"update.venturemark.co/id":   up1,
						"venture.venturemark.co/id":  ve1,
					},
				},
			},
		}

		o, err := cli.TexUpd().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["update.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &texupd.DeleteI{
			Obj: []*texupd.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"update.venturemark.co/id":   up2,
						"venture.venturemark.co/id":  ve1,
					},
				},
			},
		}

		o, err := cli.TexUpd().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["update.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: []*timeline.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"venture.venturemark.co/id":  ve1,
					},
					Jsnpatch: []*timeline.UpdateI_Obj_Jsnpatch{
						{
							Ope: "replace",
							Pat: "/obj/property/stat",
							Val: to.StringP("archived"),
						},
					},
				},
			},
		}

		o, err := cli.Timeline().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("timeline status must not be empty")
		}

		if s != "updated" {
			t.Fatal("timeline status must be updated")
		}
	}

	{
		i := &timeline.DeleteI{
			Obj: []*timeline.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"venture.venturemark.co/id":  ve1,
					},
				},
			},
		}

		o, err := cli.Timeline().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("timeline status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("timeline status must be deleted")
		}
	}

	{
		i := &update.SearchI{
			Obj: []*update.SearchI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"venture.venturemark.co/id":  ve1,
					},
				},
			},
		}

		o, err := cli.Update().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 0 {
			t.Fatal("there must be zero updates")
		}
	}
}

// Test_TexUpd_002 ensures text updates can only be created for timelines that
// already exist.
func Test_TexUpd_002(t *testing.T) {
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
		i := &texupd.CreateI{
			Obj: []*texupd.CreateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": "1",
						"venture.venturemark.co/id":  "1",
					},
					Property: &texupd.CreateI_Obj_Property{
						Text: "Lorem ipsum 1",
					},
				},
			},
		}

		_, err := cli.TexUpd().Create(context.Background(), i)
		if err == nil {
			t.Fatal("update must not be created without timeline")
		}
	}
}

// Test_TexUpd_003 ensures that the cascaded deletion of updates is working as
// ecpected.
func Test_TexUpd_003(t *testing.T) {
	var err error

	var b budget.Interface
	{
		c := budget.ConstantConfig{
			Budget:   9,
			Duration: 5 * time.Second,
		}

		b, err = budget.NewConstant(c)
		if err != nil {
			panic(err)
		}
	}

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

		o, err := cli.Venture().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["venture.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		vei = s
	}

	var ti1 string
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

		o, err := cli.Timeline().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["timeline.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ti1 = s
	}

	var ti2 string
	{
		i := &timeline.CreateI{
			Obj: []*timeline.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
					Property: &timeline.CreateI_Obj_Property{
						Name: "Internal Project",
					},
				},
			},
		}

		o, err := cli.Timeline().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["timeline.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ti2 = s
	}

	var up1 string
	{
		i := &texupd.CreateI{
			Obj: []*texupd.CreateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti1,
						"venture.venturemark.co/id":  vei,
					},
					Property: &texupd.CreateI_Obj_Property{
						Text: "Lorem ipsum 1",
					},
				},
			},
		}

		o, err := cli.TexUpd().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["update.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		up1 = s
	}

	var up2 string
	{
		i := &texupd.CreateI{
			Obj: []*texupd.CreateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti2,
						"venture.venturemark.co/id":  vei,
					},
					Property: &texupd.CreateI_Obj_Property{
						Text: "Lorem ipsum 2",
					},
				},
			},
		}

		o, err := cli.TexUpd().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["update.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		up2 = s
	}

	{
		i := &message.CreateI{
			Obj: []*message.CreateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti1,
						"update.venturemark.co/id":   up1,
						"venture.venturemark.co/id":  vei,
					},
					Property: &message.CreateI_Obj_Property{
						Text: "Lorem ipsum 1",
					},
				},
			},
		}

		o, err := cli.Message().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		_, ok := o.Obj[0].Metadata["message.venturemark.co/id"]
		if !ok {
			t.Fatal("message ID must not be empty")
		}
	}

	var me2 string
	{
		i := &message.CreateI{
			Obj: []*message.CreateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti2,
						"update.venturemark.co/id":   up2,
						"venture.venturemark.co/id":  vei,
					},
					Property: &message.CreateI_Obj_Property{
						Text: "Lorem ipsum 2",
					},
				},
			},
		}

		o, err := cli.Message().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["message.venturemark.co/id"]
		if !ok {
			t.Fatal("message ID must not be empty")
		}

		me2 = s
	}

	{
		i := &texupd.DeleteI{
			Obj: []*texupd.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti1,
						"update.venturemark.co/id":   up1,
						"venture.venturemark.co/id":  vei,
					},
				},
			},
		}

		o, err := cli.TexUpd().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["update.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	// Here we check for deleted messages, updates, timelines to be gone after
	// the apiworker cleaned them up eventually. Since the deletion is an
	// asynchronous process we have to check for the desired result multiple
	// times. Eventually the current state of the system should be reconciled
	// towards the desired state of deleted resources.

	{
		o := func() error {
			i := &message.SearchI{
				Obj: []*message.SearchI_Obj{
					{
						Metadata: map[string]string{
							"timeline.venturemark.co/id": ti1,
							"update.venturemark.co/id":   up1,
							"venture.venturemark.co/id":  vei,
						},
					},
				},
			}

			o, err := cli.Message().Search(context.Background(), i)
			if err != nil {
				t.Fatal(err)
			}

			if len(o.Obj) != 0 {
				return tracer.Mask(fmt.Errorf("there must be zero messages"))
			}

			return nil
		}

		err = b.Execute(o)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		o := func() error {
			i := &update.SearchI{
				Obj: []*update.SearchI_Obj{
					{
						Metadata: map[string]string{
							"timeline.venturemark.co/id": ti1,
							"venture.venturemark.co/id":  vei,
						},
					},
				},
			}

			o, err := cli.Update().Search(context.Background(), i)
			if err != nil {
				return tracer.Mask(err)
			}

			if len(o.Obj) != 0 {
				return tracer.Mask(fmt.Errorf("there must be zero updates"))
			}

			return nil
		}

		err = b.Execute(o)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		i := &message.SearchI{
			Obj: []*message.SearchI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti2,
						"update.venturemark.co/id":   up2,
						"venture.venturemark.co/id":  vei,
					},
				},
			},
		}

		o, err := cli.Message().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one message")
		}
		s, ok := o.Obj[0].Metadata["message.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}
		if s != me2 {
			t.Fatal("id must match")
		}
	}

	{
		i := &update.SearchI{
			Obj: []*update.SearchI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti2,
						"venture.venturemark.co/id":  vei,
					},
				},
			},
		}

		o, err := cli.Update().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one update")
		}
		s, ok := o.Obj[0].Metadata["update.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}
		if s != up2 {
			t.Fatal("id must match")
		}
	}

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
				},
			},
		}

		o, err := cli.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 2 {
			t.Fatal("there must be two timelines")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: []*timeline.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti1,
						"venture.venturemark.co/id":  vei,
					},
					Jsnpatch: []*timeline.UpdateI_Obj_Jsnpatch{
						{
							Ope: "replace",
							Pat: "/obj/property/stat",
							Val: to.StringP("archived"),
						},
					},
				},
			},
		}

		o, err := cli.Timeline().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "updated" {
			t.Fatal("status must be updated")
		}
	}

	{
		i := &timeline.DeleteI{
			Obj: []*timeline.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti1,
						"venture.venturemark.co/id":  vei,
					},
				},
			},
		}

		o, err := cli.Timeline().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: []*timeline.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti2,
						"venture.venturemark.co/id":  vei,
					},
					Jsnpatch: []*timeline.UpdateI_Obj_Jsnpatch{
						{
							Ope: "replace",
							Pat: "/obj/property/stat",
							Val: to.StringP("archived"),
						},
					},
				},
			},
		}

		o, err := cli.Timeline().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "updated" {
			t.Fatal("status must be updated")
		}
	}

	{
		i := &timeline.DeleteI{
			Obj: []*timeline.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti2,
						"venture.venturemark.co/id":  vei,
					},
				},
			},
		}

		o, err := cli.Timeline().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["timeline.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		o := func() error {
			i := &message.SearchI{
				Obj: []*message.SearchI_Obj{
					{
						Metadata: map[string]string{
							"timeline.venturemark.co/id": ti2,
							"update.venturemark.co/id":   up2,
							"venture.venturemark.co/id":  vei,
						},
					},
				},
			}

			o, err := cli.Message().Search(context.Background(), i)
			if err != nil {
				t.Fatal(err)
			}

			if len(o.Obj) != 0 {
				return tracer.Mask(fmt.Errorf("there must be zero messages"))
			}

			return nil
		}

		err = b.Execute(o)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		o := func() error {
			i := &update.SearchI{
				Obj: []*update.SearchI_Obj{
					{
						Metadata: map[string]string{
							"timeline.venturemark.co/id": ti2,
							"venture.venturemark.co/id":  vei,
						},
					},
				},
			}

			o, err := cli.Update().Search(context.Background(), i)
			if err != nil {
				return tracer.Mask(err)
			}

			if len(o.Obj) != 0 {
				return tracer.Mask(fmt.Errorf("there must be zero updates"))
			}

			return nil
		}

		err = b.Execute(o)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
				},
			},
		}

		o, err := cli.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 0 {
			t.Fatal("there must be zero timelines")
		}
	}
}

// Test_TexUpd_004 ensures that deleting update resources which do not exist
// returns an error.
func Test_TexUpd_004(t *testing.T) {
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
		i := &texupd.DeleteI{
			Obj: []*texupd.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": "1",
						"update.venturemark.co/id":   "1",
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		_, err := cli.TexUpd().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}
}

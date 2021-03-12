// +build conformance

package tst

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/venturemark/apigengo/pkg/pbf/audience"
	"github.com/venturemark/apigengo/pkg/pbf/message"
	"github.com/venturemark/apigengo/pkg/pbf/texupd"
	"github.com/venturemark/apigengo/pkg/pbf/timeline"
	"github.com/venturemark/apigengo/pkg/pbf/update"
	"github.com/xh3b4sd/budget"
	"github.com/xh3b4sd/tracer"

	"github.com/venturemark/cfm/pkg/client"
	"github.com/venturemark/cfm/pkg/oauth"
	"github.com/venturemark/cfm/pkg/to"
)

// Test_Timeline_001 ensures that the lifecycle of timelines is covered from
// creation to deletion.
func Test_Timeline_001(t *testing.T) {
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

	var ti1 string
	{
		i := &timeline.CreateI{
			Obj: []*timeline.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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
						"venture.venturemark.co/id": "1",
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

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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

		{
			s, ok := o.Obj[0].Metadata["timeline.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ti2 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[0].Property.Name != "Internal Project" {
				t.Fatal("name must be Internal Project")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["timeline.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != ti1 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[1].Property.Name != "Marketing Campaign" {
				t.Fatal("name must be Marketing Campaign")
			}
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: []*timeline.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti1,
						"venture.venturemark.co/id":  "1",
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
						"venture.venturemark.co/id":  "1",
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
						"venture.venturemark.co/id":  "1",
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
						"venture.venturemark.co/id":  "1",
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
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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

// Test_Timeline_002 ensures that timeline names are unique.
func Test_Timeline_002(t *testing.T) {
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

	var tii string
	{
		i := &timeline.CreateI{
			Obj: []*timeline.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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

		tii = s
	}

	{
		i := &timeline.CreateI{
			Obj: []*timeline.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
					Property: &timeline.CreateI_Obj_Property{
						Name: "Marketing Campaign",
					},
				},
			},
		}

		_, err := cli.Timeline().Create(context.Background(), i)
		if err == nil {
			t.Fatal("name must be unique")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: []*timeline.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"venture.venturemark.co/id":  "1",
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
						"timeline.venturemark.co/id": tii,
						"venture.venturemark.co/id":  "1",
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
}

// Test_Timeline_003 ensures that the timeline state can be updated.
func Test_Timeline_003(t *testing.T) {
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

	var tii string
	{
		i := &timeline.CreateI{
			Obj: []*timeline.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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

		tii = s
	}

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
				},
			},
		}

		o, err := cli.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one timeline")
		}

		if o.Obj[0].Property.Desc != "" {
			t.Fatal("desc must be empty")
		}
		if o.Obj[0].Property.Name != "Marketing Campaign" {
			t.Fatal("name must be Internal Project")
		}
		if o.Obj[0].Property.Stat != "active" {
			t.Fatal("stat must be active")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: []*timeline.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"venture.venturemark.co/id":  "1",
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
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
				},
			},
		}

		o, err := cli.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one timeline")
		}

		if o.Obj[0].Property.Desc != "" {
			t.Fatal("desc must be empty")
		}
		if o.Obj[0].Property.Name != "Marketing Campaign" {
			t.Fatal("name must be Internal Project")
		}
		if o.Obj[0].Property.Stat != "archived" {
			t.Fatal("stat must be archived")
		}
	}

	{
		i := &timeline.DeleteI{
			Obj: []*timeline.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"venture.venturemark.co/id":  "1",
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
}

// Test_Timeline_004 ensures that timelines not having the archived state cannot
// be deleted.
func Test_Timeline_004(t *testing.T) {
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

	var tii string
	{
		i := &timeline.CreateI{
			Obj: []*timeline.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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

		tii = s
	}

	{
		i := &timeline.DeleteI{
			Obj: []*timeline.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		_, err := cli.Timeline().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("state must be archived for deletion")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: []*timeline.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": tii,
						"venture.venturemark.co/id":  "1",
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
						"timeline.venturemark.co/id": tii,
						"venture.venturemark.co/id":  "1",
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
}

// Test_Timeline_005 ensures that all timelines can be shown to users with
// permissions disabled.
func Test_Timeline_005(t *testing.T) {
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

	var ti1 string
	{
		i := &timeline.CreateI{
			Obj: []*timeline.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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

		ti1 = s
	}

	var ti2 string
	{
		i := &timeline.CreateI{
			Obj: []*timeline.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
					Property: &timeline.CreateI_Obj_Property{
						Name: "Internal Project",
					},
				},
			},
		}

		o, err := cl2.Timeline().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["timeline.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ti2 = s
	}

	var au1 string
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
							ti1,
						},
						User: []string{
							cr1.User(),
						},
					},
				},
			},
		}

		o, err := cl1.Audience().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/id"]
		if !ok {
			t.Fatal("audience ID must not be empty")
		}

		au1 = s
	}

	var au2 string
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
							ti2,
						},
						User: []string{
							cr2.User(),
						},
					},
				},
			},
		}

		o, err := cl2.Audience().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/id"]
		if !ok {
			t.Fatal("audience ID must not be empty")
		}

		au2 = s
	}

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
				},
			},
		}

		o, err := cl1.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 2 {
			t.Fatal("there must be two timelines")
		}
	}

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
				},
			},
		}

		o, err := cl2.Timeline().Search(context.Background(), i)
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
						"venture.venturemark.co/id":  "1",
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

		o, err := cl1.Timeline().Update(context.Background(), i)
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
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cl1.Timeline().Delete(context.Background(), i)
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
						"venture.venturemark.co/id":  "1",
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

		o, err := cl2.Timeline().Update(context.Background(), i)
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
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cl2.Timeline().Delete(context.Background(), i)
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
		i := &audience.DeleteI{
			Obj: []*audience.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"audience.venturemark.co/id": au1,
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cl1.Audience().Delete(context.Background(), i)
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
						"audience.venturemark.co/id": au2,
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cl2.Audience().Delete(context.Background(), i)
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

// Test_Timeline_006 ensures that only matching timelines can be shown to users with
// permissions enabled.
func Test_Timeline_006(t *testing.T) {
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

	var ti1 string
	{
		i := &timeline.CreateI{
			Obj: []*timeline.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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

		ti1 = s
	}

	var ti2 string
	{
		i := &timeline.CreateI{
			Obj: []*timeline.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
					},
					Property: &timeline.CreateI_Obj_Property{
						Name: "Internal Project",
					},
				},
			},
		}

		o, err := cl2.Timeline().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["timeline.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		ti2 = s
	}

	var au1 string
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
							ti1,
						},
						User: []string{
							cr1.User(),
						},
					},
				},
			},
		}

		o, err := cl1.Audience().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/id"]
		if !ok {
			t.Fatal("audience ID must not be empty")
		}

		au1 = s
	}

	var au2 string
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
							ti2,
						},
						User: []string{
							cr2.User(),
						},
					},
				},
			},
		}

		o, err := cl2.Audience().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["audience.venturemark.co/id"]
		if !ok {
			t.Fatal("audience ID must not be empty")
		}

		au2 = s
	}

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"permission.venturemark.co/id":     "audience",
						"permission.venturemark.co/status": "enabled",
						"venture.venturemark.co/id":        "1",
					},
				},
			},
		}

		o, err := cl1.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one timeline")
		}
		if o.Obj[0].Property.Name != "Marketing Campaign" {
			t.Fatal("name must be Marketing Campaign")
		}
	}

	{
		i := &timeline.SearchI{
			Obj: []*timeline.SearchI_Obj{
				{
					Metadata: map[string]string{
						"permission.venturemark.co/id":     "audience",
						"permission.venturemark.co/status": "enabled",
						"venture.venturemark.co/id":        "1",
					},
				},
			},
		}

		o, err := cl2.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one timeline")
		}
		if o.Obj[0].Property.Name != "Internal Project" {
			t.Fatal("name must be Internal Project")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: []*timeline.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti1,
						"venture.venturemark.co/id":  "1",
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

		o, err := cl1.Timeline().Update(context.Background(), i)
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
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cl1.Timeline().Delete(context.Background(), i)
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
						"venture.venturemark.co/id":  "1",
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

		o, err := cl2.Timeline().Update(context.Background(), i)
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
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cl2.Timeline().Delete(context.Background(), i)
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
		i := &audience.DeleteI{
			Obj: []*audience.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"audience.venturemark.co/id": au1,
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cl1.Audience().Delete(context.Background(), i)
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
						"audience.venturemark.co/id": au2,
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		o, err := cl2.Audience().Delete(context.Background(), i)
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

// Test_Timeline_007 ensures that the cascaded deletion of timelines is working
// as ecpected.
func Test_Timeline_007(t *testing.T) {
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

	var ti1 string
	{
		i := &timeline.CreateI{
			Obj: []*timeline.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": "1",
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
						"venture.venturemark.co/id": "1",
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
						"venture.venturemark.co/id":  "1",
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
			t.Fatal("texupd ID must not be empty")
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
						"venture.venturemark.co/id":  "1",
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
			t.Fatal("texupd ID must not be empty")
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
						"venture.venturemark.co/id":  "1",
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
						"venture.venturemark.co/id":  "1",
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
		i := &timeline.UpdateI{
			Obj: []*timeline.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti1,
						"venture.venturemark.co/id":  "1",
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
						"venture.venturemark.co/id":  "1",
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
							"venture.venturemark.co/id":  "1",
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
							"venture.venturemark.co/id":  "1",
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
						"venture.venturemark.co/id":  "1",
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
						"venture.venturemark.co/id":  "1",
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
						"venture.venturemark.co/id": "1",
					},
				},
			},
		}

		o, err := cli.Timeline().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one timeline")
		}
	}

	{
		i := &timeline.UpdateI{
			Obj: []*timeline.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": ti2,
						"venture.venturemark.co/id":  "1",
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
						"venture.venturemark.co/id":  "1",
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
							"venture.venturemark.co/id":  "1",
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
							"venture.venturemark.co/id":  "1",
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
						"venture.venturemark.co/id": "1",
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

// Test_Timeline_008 ensures that deleting timeline resources which do not exist
// returns an error.
func Test_Timeline_008(t *testing.T) {
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
		i := &timeline.DeleteI{
			Obj: []*timeline.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"timeline.venturemark.co/id": "1",
						"venture.venturemark.co/id":  "1",
					},
				},
			},
		}

		_, err := cli.Timeline().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}
}

// +build conformance

package tst

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/venturemark/apigengo/pkg/pbf/invite"
	"github.com/venturemark/apigengo/pkg/pbf/role"
	"github.com/venturemark/apigengo/pkg/pbf/user"
	"github.com/venturemark/apigengo/pkg/pbf/venture"
	"github.com/xh3b4sd/budget"
	"github.com/xh3b4sd/tracer"

	"github.com/venturemark/cfm/pkg/client"
	"github.com/venturemark/cfm/pkg/oauth"
	"github.com/venturemark/cfm/pkg/to"
)

// Test_Invite_001 ensures that the lifecycle of invites is covered from
// creation to deletion.
func Test_Invite_001(t *testing.T) {
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

		o, err := cli.User().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		_, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}
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

	var in1 string
	{
		i := &invite.CreateI{
			Obj: []*invite.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
					Property: &invite.CreateI_Obj_Property{
						Mail: "user1@site.net",
					},
				},
			},
		}

		o, err := cli.Invite().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["invite.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		in1 = s
	}

	var in2 string
	{
		i := &invite.CreateI{
			Obj: []*invite.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
					Property: &invite.CreateI_Obj_Property{
						Mail: "user2@site.net",
					},
				},
			},
		}

		o, err := cli.Invite().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["invite.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}

		in2 = s
	}

	{
		i := &invite.SearchI{
			Obj: []*invite.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
				},
			},
		}

		o, err := cli.Invite().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 2 {
			t.Fatal("there must be two invites")
		}

		{
			s, ok := o.Obj[0].Metadata["invite.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != in2 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[0].Property.Mail != "user2@site.net" {
				t.Fatal("name must be user2@site.net")
			}
			if o.Obj[0].Property.Stat != "pending" {
				t.Fatal("name must be pending")
			}
		}

		{
			s, ok := o.Obj[1].Metadata["invite.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != in1 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[1].Property.Mail != "user1@site.net" {
				t.Fatal("name must be user1@site.net")
			}
			if o.Obj[1].Property.Stat != "pending" {
				t.Fatal("name must be pending")
			}
		}
	}

	{
		i := &invite.SearchI{
			Obj: []*invite.SearchI_Obj{
				{
					Metadata: map[string]string{
						"subject.venturemark.co/email": "user1@site.net",
						"venture.venturemark.co/id":    vei,
					},
				},
			},
		}

		o, err := cli.Invite().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one invite")
		}

		{
			s, ok := o.Obj[0].Metadata["invite.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != in1 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[0].Property.Mail != "user1@site.net" {
				t.Fatal("name must be user1@site.net")
			}
			if o.Obj[0].Property.Stat != "pending" {
				t.Fatal("name must be pending")
			}
		}
	}

	{
		i := &invite.UpdateI{
			Obj: []*invite.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"invite.venturemark.co/id":  in1,
						"venture.venturemark.co/id": vei,
					},
					Jsnpatch: []*invite.UpdateI_Obj_Jsnpatch{
						{
							Ope: "replace",
							Pat: "/obj/property/stat",
							Val: to.StringP("accepted"),
						},
					},
				},
			},
		}

		o, err := cli.Invite().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["invite.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "updated" {
			t.Fatal("status must be updated")
		}
	}

	{
		i := &invite.SearchI{
			Obj: []*invite.SearchI_Obj{
				{
					Metadata: map[string]string{
						"subject.venturemark.co/email": "user1@site.net",
						"venture.venturemark.co/id":    vei,
					},
				},
			},
		}

		o, err := cli.Invite().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one invite")
		}

		{
			s, ok := o.Obj[0].Metadata["invite.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != in1 {
				t.Fatal("id must match across actions")
			}
			if o.Obj[0].Property.Mail != "user1@site.net" {
				t.Fatal("name must be user1@site.net")
			}
			if o.Obj[0].Property.Stat != "accepted" {
				t.Fatal("name must be accepted")
			}
		}
	}

	{
		i := &invite.DeleteI{
			Obj: []*invite.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"invite.venturemark.co/id":  in1,
						"venture.venturemark.co/id": vei,
					},
				},
			},
		}

		o, err := cli.Invite().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["invite.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &invite.UpdateI{
			Obj: []*invite.UpdateI_Obj{
				{
					Metadata: map[string]string{
						"invite.venturemark.co/id":  in2,
						"venture.venturemark.co/id": vei,
					},
					Jsnpatch: []*invite.UpdateI_Obj_Jsnpatch{
						{
							Ope: "replace",
							Pat: "/obj/property/stat",
							Val: to.StringP("rejected"),
						},
					},
				},
			},
		}

		o, err := cli.Invite().Update(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["invite.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "updated" {
			t.Fatal("status must be updated")
		}
	}

	{
		i := &invite.DeleteI{
			Obj: []*invite.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"invite.venturemark.co/id":  in2,
						"venture.venturemark.co/id": vei,
					},
				},
			},
		}

		o, err := cli.Invite().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["invite.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		o := func() error {
			i := &role.SearchI{
				Obj: []*role.SearchI_Obj{
					{
						Metadata: map[string]string{
							"resource.venturemark.co/kind": "invite",
							"invite.venturemark.co/id":     in2,
							"venture.venturemark.co/id":    vei,
						},
					},
				},
			}

			o, err := cli.Role().Search(context.Background(), i)
			if err != nil {
				t.Fatal(err)
			}

			if len(o.Obj) != 0 {
				return tracer.Mask(fmt.Errorf("there must be zero roles"))
			}

			return nil
		}

		err = b.Execute(o)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		i := &invite.SearchI{
			Obj: []*invite.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
				},
			},
		}

		o, err := cli.Invite().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 0 {
			t.Fatal("there must be zero invites")
		}
	}

	{
		i := &venture.DeleteI{
			Obj: []*venture.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
				},
			},
		}

		_, err := cli.Venture().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		i := &user.DeleteI{}

		_, err := cli.User().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}
	}

	{
		o := func() error {
			emp, err := cli.Redigo().Empty()
			if err != nil {
				t.Fatal(err)
			}

			if !emp {
				return tracer.Mask(fmt.Errorf("storage must be empty"))
			}

			return nil
		}

		err = b.Execute(o)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// Test_Invite_002 ensures that emails are unique.
func Test_Invite_002(t *testing.T) {
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
		i := &user.CreateI{
			Obj: []*user.CreateI_Obj{
				{
					Property: &user.CreateI_Obj_Property{
						Name: "marcojelli",
					},
				},
			},
		}

		o, err := cli.User().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		_, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}
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

	{
		i := &invite.CreateI{
			Obj: []*invite.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
					Property: &invite.CreateI_Obj_Property{
						Mail: "user1@site.net",
					},
				},
			},
		}

		o, err := cli.Invite().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		_, ok := o.Obj[0].Metadata["invite.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}
	}

	{
		i := &invite.CreateI{
			Obj: []*invite.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
					Property: &invite.CreateI_Obj_Property{
						Mail: "user1@site.net",
					},
				},
			},
		}

		_, err := cli.Invite().Create(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}
}

// Test_Invite_003 ensures that invites can only be created by users who are
// owners of a venture. Additionally the test verifies that a legitimate email
// address must be specified.
func Test_Invite_003(t *testing.T) {
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

		_, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}
	}

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

		_, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}
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

	{
		i := &invite.CreateI{
			Obj: []*invite.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
					Property: &invite.CreateI_Obj_Property{
						Mail: "user1@site.net",
					},
				},
			},
		}

		o, err := cl1.Invite().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		_, ok := o.Obj[0].Metadata["invite.venturemark.co/id"]
		if !ok {
			t.Fatal("id must not be empty")
		}
	}

	{
		i := &invite.CreateI{
			Obj: []*invite.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
					Property: &invite.CreateI_Obj_Property{
						Mail: "user2@site.net",
					},
				},
			},
		}

		_, err := cl2.Invite().Create(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}

	{
		i := &invite.CreateI{
			Obj: []*invite.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
					Property: &invite.CreateI_Obj_Property{
						Mail: "",
					},
				},
			},
		}

		_, err := cl1.Invite().Create(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}

	{
		i := &invite.CreateI{
			Obj: []*invite.CreateI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": vei,
					},
					Property: &invite.CreateI_Obj_Property{
						Mail: "garbage",
					},
				},
			},
		}

		_, err := cl1.Invite().Create(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}
}

// Test_Invite_004 ensures that deleting invite resources which do not exist
// returns an error.
func Test_Invite_004(t *testing.T) {
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
		i := &invite.DeleteI{
			Obj: []*invite.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"invite.venturemark.co/id":  "1",
						"venture.venturemark.co/id": "1",
					},
				},
			},
		}

		_, err := cli.Invite().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}
}

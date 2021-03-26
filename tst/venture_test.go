// +build conformance

package tst

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/venturemark/apigengo/pkg/pbf/role"
	"github.com/venturemark/apigengo/pkg/pbf/timeline"
	"github.com/venturemark/apigengo/pkg/pbf/user"
	"github.com/venturemark/apigengo/pkg/pbf/venture"
	"github.com/xh3b4sd/budget"
	"github.com/xh3b4sd/tracer"

	"github.com/venturemark/cfm/pkg/client"
	"github.com/venturemark/cfm/pkg/oauth"
)

// Test_Venture_001 ensures that the lifecycle of ventures is covered from
// creation to deletion.
func Test_Venture_001(t *testing.T) {
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
		i := &role.SearchI{
			Obj: []*role.SearchI_Obj{
				{
					Metadata: map[string]string{
						"resource.venturemark.co/kind": "venture",
						"venture.venturemark.co/id":    ve1,
					},
				},
			},
		}

		o, err := cl1.Role().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one role")
		}
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
		i := &role.SearchI{
			Obj: []*role.SearchI_Obj{
				{
					Metadata: map[string]string{
						"resource.venturemark.co/kind": "venture",
						"venture.venturemark.co/id":    ve2,
					},
				},
			},
		}

		o, err := cl1.Role().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one role")
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

		_, err := cl2.Venture().Search(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}

	}

	{
		i := &role.CreateI{
			Obj: []*role.CreateI_Obj{
				{
					Metadata: map[string]string{
						"resource.venturemark.co/kind": "venture",
						"role.venturemark.co/kind":     "member",
						"subject.venturemark.co/id":    us2,
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
		i := &role.CreateI{
			Obj: []*role.CreateI_Obj{
				{
					Metadata: map[string]string{
						"resource.venturemark.co/kind": "timeline",
						"role.venturemark.co/kind":     "member",
						"subject.venturemark.co/id":    us2,
						"timeline.venturemark.co/id":   tii,
						"venture.venturemark.co/id":    ve1,
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
						"subject.venturemark.co/id": us1,
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
						"subject.venturemark.co/id": us2,
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
		i := &venture.SearchI{
			Obj: []*venture.SearchI_Obj{
				{
					Metadata: map[string]string{
						"venture.venturemark.co/id": ve1,
					},
				},
			},
		}

		o, err := cl2.Venture().Search(context.Background(), i)
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
			s, ok := o.Obj[0].Metadata["venture.venturemark.co/status"]
			if !ok {
				t.Fatal("status must not be empty")
			}

			if s != "deleted" {
				t.Fatal("status must be deleted")
			}
		}
	}

	{
		o := func() error {
			i := &role.SearchI{
				Obj: []*role.SearchI_Obj{
					{
						Metadata: map[string]string{
							"resource.venturemark.co/kind": "venture",
							"venture.venturemark.co/id":    ve1,
						},
					},
				},
			}

			o, err := cl1.Role().Search(context.Background(), i)
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
		o := func() error {
			i := &role.SearchI{
				Obj: []*role.SearchI_Obj{
					{
						Metadata: map[string]string{
							"resource.venturemark.co/kind": "venture",
							"venture.venturemark.co/id":    ve2,
						},
					},
				},
			}

			o, err := cl1.Role().Search(context.Background(), i)
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
		i := &venture.SearchI{
			Obj: []*venture.SearchI_Obj{
				{
					Metadata: map[string]string{
						"subject.venturemark.co/id": us1,
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
						"subject.venturemark.co/id": us2,
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
		i := &user.DeleteI{
			Obj: []*user.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"user.venturemark.co/id": us1,
					},
				},
			},
		}

		o, err := cl1.User().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["user.venturemark.co/status"]
		if !ok {
			t.Fatal("status must not be empty")
		}

		if s != "deleted" {
			t.Fatal("status must be deleted")
		}
	}

	{
		i := &user.DeleteI{
			Obj: []*user.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"user.venturemark.co/id": us2,
					},
				},
			},
		}

		o, err := cl2.User().Delete(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		s, ok := o.Obj[0].Metadata["user.venturemark.co/status"]
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
							"resource.venturemark.co/kind": "user",
							"user.venturemark.co/id":       us1,
						},
					},
				},
			}

			o, err := cl1.Role().Search(context.Background(), i)
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
		i := &user.SearchI{
			Obj: []*user.SearchI_Obj{
				{
					Metadata: map[string]string{
						"subject.venturemark.co/id": us1,
					},
				},
			},
		}

		o, err := cl1.User().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 0 {
			t.Fatal("there must be zero users")
		}
	}

	{
		o := func() error {
			i := &role.SearchI{
				Obj: []*role.SearchI_Obj{
					{
						Metadata: map[string]string{
							"resource.venturemark.co/kind": "user",
							"user.venturemark.co/id":       us2,
						},
					},
				},
			}

			o, err := cl2.Role().Search(context.Background(), i)
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
		i := &user.SearchI{
			Obj: []*user.SearchI_Obj{
				{
					Metadata: map[string]string{
						"subject.venturemark.co/id": us2,
					},
				},
			},
		}

		o, err := cl2.User().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 0 {
			t.Fatal("there must be zero users")
		}
	}

	{
		o := func() error {
			emp, err := cl1.Redigo().Empty()
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

// +build conformance

package tst

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/venturemark/apigengo/pkg/pbf/role"
	"github.com/venturemark/apigengo/pkg/pbf/user"
	"github.com/xh3b4sd/budget"
	"github.com/xh3b4sd/tracer"

	"github.com/venturemark/cfm/pkg/client"
	"github.com/venturemark/cfm/pkg/oauth"
)

// Test_User_001 ensures that the lifecycle of users is covered from
// creation to deletion.
func Test_User_001(t *testing.T) {
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

	{
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

		if len(o.Obj) != 1 {
			t.Fatal("there must be one role")
		}
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

	{
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

		if len(o.Obj) != 1 {
			t.Fatal("there must be one role")
		}
	}

	{
		i := &user.SearchI{}

		o, err := cl1.User().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one user")
		}

		{
			s, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != us1 {
				t.Fatal("id must match across actions")
			}
		}

		{
			if o.Obj[0].Property.Name != "marcojelli" {
				t.Fatal("name must be marcojelli")
			}
		}
	}

	{
		i := &user.SearchI{
			Obj: []*user.SearchI_Obj{
				{
					Metadata: map[string]string{
						"user.venturemark.co/id": us1,
					},
				},
			},
		}

		o, err := cl1.User().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one user")
		}

		{
			s, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != us1 {
				t.Fatal("id must match across actions")
			}
		}

		{
			if o.Obj[0].Property.Name != "marcojelli" {
				t.Fatal("name must be marcojelli")
			}
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

		if len(o.Obj) != 1 {
			t.Fatal("there must be one user")
		}

		{
			s, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != us1 {
				t.Fatal("id must match across actions")
			}
		}

		{
			if o.Obj[0].Property.Name != "marcojelli" {
				t.Fatal("name must be marcojelli")
			}
		}
	}

	{
		i := &user.SearchI{
			Obj: []*user.SearchI_Obj{
				{
					Metadata: map[string]string{
						"user.venturemark.co/id": us2,
					},
				},
			},
		}

		_, err := cl1.User().Search(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
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

		o, err := cl1.User().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one user")
		}

		{
			s, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != us1 {
				t.Fatal("id must match across actions")
			}
		}

		{
			if o.Obj[0].Property.Name != "marcojelli" {
				t.Fatal("name must be marcojelli")
			}
		}
	}

	{
		i := &user.SearchI{}

		o, err := cl2.User().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one user")
		}

		{
			s, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != us2 {
				t.Fatal("id must match across actions")
			}
		}

		{
			if o.Obj[0].Property.Name != "disreszi" {
				t.Fatal("name must be disreszi")
			}
		}
	}

	{
		i := &user.SearchI{
			Obj: []*user.SearchI_Obj{
				{
					Metadata: map[string]string{
						"user.venturemark.co/id": us2,
					},
				},
			},
		}

		o, err := cl2.User().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one user")
		}

		{
			s, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != us2 {
				t.Fatal("id must match across actions")
			}
		}

		{
			if o.Obj[0].Property.Name != "disreszi" {
				t.Fatal("name must be disreszi")
			}
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

		if len(o.Obj) != 1 {
			t.Fatal("there must be one user")
		}

		{
			s, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != us2 {
				t.Fatal("id must match across actions")
			}
		}

		{
			if o.Obj[0].Property.Name != "disreszi" {
				t.Fatal("name must be disreszi")
			}
		}
	}

	{
		i := &user.SearchI{
			Obj: []*user.SearchI_Obj{
				{
					Metadata: map[string]string{
						"user.venturemark.co/id": us1,
					},
				},
			},
		}

		_, err := cl2.User().Search(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
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

		o, err := cl2.User().Search(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		if len(o.Obj) != 1 {
			t.Fatal("there must be one user")
		}

		{
			s, ok := o.Obj[0].Metadata["user.venturemark.co/id"]
			if !ok {
				t.Fatal("id must not be empty")
			}
			if s != us2 {
				t.Fatal("id must match across actions")
			}
		}

		{
			if o.Obj[0].Property.Name != "disreszi" {
				t.Fatal("name must be disreszi")
			}
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

		_, err := cl2.User().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
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

		_, err := cl1.User().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}

	{
		i := &user.DeleteI{}

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
}

// Test_User_002 ensures that deleting user resources which do not exist
// returns an error.
func Test_User_002(t *testing.T) {
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
		i := &user.DeleteI{
			Obj: []*user.DeleteI_Obj{
				{
					Metadata: map[string]string{
						"user.venturemark.co/id": "1",
					},
				},
			},
		}

		_, err := cli.User().Delete(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}
}

// Test_User_003 ensures that the users can only create one user object for
// themselves.
func Test_User_003(t *testing.T) {
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
						Name: "one",
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

	{
		i := &user.CreateI{
			Obj: []*user.CreateI_Obj{
				{
					Property: &user.CreateI_Obj_Property{
						Name: "two",
					},
				},
			},
		}

		_, err := cli.User().Create(context.Background(), i)
		if err == nil {
			t.Fatal("error must not be empty")
		}
	}
}

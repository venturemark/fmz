// +build conformance

package tst

import (
	"context"
	"testing"

	"github.com/venturemark/apigengo/pkg/pbf/audience"

	"github.com/venturemark/fmz/pkg/client"
)

func Test_Audience_001(t *testing.T) {
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

	{
		i := &audience.CreateI{
			Obj: &audience.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "org",
					"user.venturemark.co/id":         "usr",
				},
				Property: &audience.CreateI_Obj_Property{
					Name: "Employees",
					User: []string{
						"xh3b4sd",
						"marcoelli",
					},
				},
			},
		}

		o, err := cli.Audience().Create(context.Background(), i)
		if err != nil {
			t.Fatal(err)
		}

		_, ok := o.Obj.Metadata["audience.venturemark.co/id"]
		if !ok {
			t.Fatal("audience ID must not be empty")
		}
	}

	{
		i := &audience.CreateI{
			Obj: &audience.CreateI_Obj{
				Metadata: map[string]string{
					"organization.venturemark.co/id": "org",
					"user.venturemark.co/id":         "usr",
				},
				Property: &audience.CreateI_Obj_Property{
					Name: "Employees",
					User: []string{
						"xh3b4sd",
						"marcoelli",
					},
				},
			},
		}

		_, err := cli.Audience().Create(context.Background(), i)
		if err == nil {
			t.Fatal("audience name must be unique")
		}
	}
}

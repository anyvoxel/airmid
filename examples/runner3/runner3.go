package runner3

import (
	"context"
	"fmt"
	"reflect"

	"github.com/anyvoxel/airmid/anvil"
	"github.com/anyvoxel/airmid/app"
	"github.com/anyvoxel/airmid/ioc"
)

type selfString string

type runner3 struct {
	i selfString `airmid:"value:${test}"`
}

func (r *runner3) Run(ctx context.Context) {
	fmt.Println(r.i)
}

func (r *runner3) Stop(ctx context.Context) {

}

func init() {
	anvil.Must(
		app.RegisterBeanDefinition("runner3", ioc.MustNewBeanDefinition(reflect.TypeOf((*runner3)(nil)))),
	)
}

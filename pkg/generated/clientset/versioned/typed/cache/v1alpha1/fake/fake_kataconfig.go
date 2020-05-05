/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/harche/kata-operator/pkg/apis/cache/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeKataconfigs implements KataconfigInterface
type FakeKataconfigs struct {
	Fake *FakeCacheV1alpha1
	ns   string
}

var kataconfigsResource = schema.GroupVersionResource{Group: "cache.example.com", Version: "v1alpha1", Resource: "kataconfigs"}

var kataconfigsKind = schema.GroupVersionKind{Group: "cache.example.com", Version: "v1alpha1", Kind: "Kataconfig"}

// Get takes name of the kataconfig, and returns the corresponding kataconfig object, and an error if there is any.
func (c *FakeKataconfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Kataconfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(kataconfigsResource, c.ns, name), &v1alpha1.Kataconfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kataconfig), err
}

// List takes label and field selectors, and returns the list of Kataconfigs that match those selectors.
func (c *FakeKataconfigs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.KataconfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(kataconfigsResource, kataconfigsKind, c.ns, opts), &v1alpha1.KataconfigList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.KataconfigList{ListMeta: obj.(*v1alpha1.KataconfigList).ListMeta}
	for _, item := range obj.(*v1alpha1.KataconfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested kataconfigs.
func (c *FakeKataconfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(kataconfigsResource, c.ns, opts))

}

// Create takes the representation of a kataconfig and creates it.  Returns the server's representation of the kataconfig, and an error, if there is any.
func (c *FakeKataconfigs) Create(ctx context.Context, kataconfig *v1alpha1.Kataconfig, opts v1.CreateOptions) (result *v1alpha1.Kataconfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(kataconfigsResource, c.ns, kataconfig), &v1alpha1.Kataconfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kataconfig), err
}

// Update takes the representation of a kataconfig and updates it. Returns the server's representation of the kataconfig, and an error, if there is any.
func (c *FakeKataconfigs) Update(ctx context.Context, kataconfig *v1alpha1.Kataconfig, opts v1.UpdateOptions) (result *v1alpha1.Kataconfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(kataconfigsResource, c.ns, kataconfig), &v1alpha1.Kataconfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kataconfig), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeKataconfigs) UpdateStatus(ctx context.Context, kataconfig *v1alpha1.Kataconfig, opts v1.UpdateOptions) (*v1alpha1.Kataconfig, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(kataconfigsResource, "status", c.ns, kataconfig), &v1alpha1.Kataconfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kataconfig), err
}

// Delete takes name of the kataconfig and deletes it. Returns an error if one occurs.
func (c *FakeKataconfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(kataconfigsResource, c.ns, name), &v1alpha1.Kataconfig{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeKataconfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(kataconfigsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.KataconfigList{})
	return err
}

// Patch applies the patch and returns the patched kataconfig.
func (c *FakeKataconfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Kataconfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(kataconfigsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Kataconfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Kataconfig), err
}

/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	scheme "github.com/networkmachinery/networkmachinery-operators/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// NetworkDelayTestsGetter has a method to return a NetworkDelayTestInterface.
// A group's client should implement this interface.
type NetworkDelayTestsGetter interface {
	NetworkDelayTests() NetworkDelayTestInterface
}

// NetworkDelayTestInterface has methods to work with NetworkDelayTest resources.
type NetworkDelayTestInterface interface {
	Create(*v1alpha1.NetworkDelayTest) (*v1alpha1.NetworkDelayTest, error)
	Update(*v1alpha1.NetworkDelayTest) (*v1alpha1.NetworkDelayTest, error)
	UpdateStatus(*v1alpha1.NetworkDelayTest) (*v1alpha1.NetworkDelayTest, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.NetworkDelayTest, error)
	List(opts v1.ListOptions) (*v1alpha1.NetworkDelayTestList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.NetworkDelayTest, err error)
	NetworkDelayTestExpansion
}

// networkDelayTests implements NetworkDelayTestInterface
type networkDelayTests struct {
	client rest.Interface
}

// newNetworkDelayTests returns a NetworkDelayTests
func newNetworkDelayTests(c *NetworkmachineryV1alpha1Client) *networkDelayTests {
	return &networkDelayTests{
		client: c.RESTClient(),
	}
}

// Get takes name of the networkDelayTest, and returns the corresponding networkDelayTest object, and an error if there is any.
func (c *networkDelayTests) Get(name string, options v1.GetOptions) (result *v1alpha1.NetworkDelayTest, err error) {
	result = &v1alpha1.NetworkDelayTest{}
	err = c.client.Get().
		Resource("networkdelaytests").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of NetworkDelayTests that match those selectors.
func (c *networkDelayTests) List(opts v1.ListOptions) (result *v1alpha1.NetworkDelayTestList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.NetworkDelayTestList{}
	err = c.client.Get().
		Resource("networkdelaytests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested networkDelayTests.
func (c *networkDelayTests) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("networkdelaytests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a networkDelayTest and creates it.  Returns the server's representation of the networkDelayTest, and an error, if there is any.
func (c *networkDelayTests) Create(networkDelayTest *v1alpha1.NetworkDelayTest) (result *v1alpha1.NetworkDelayTest, err error) {
	result = &v1alpha1.NetworkDelayTest{}
	err = c.client.Post().
		Resource("networkdelaytests").
		Body(networkDelayTest).
		Do().
		Into(result)
	return
}

// Update takes the representation of a networkDelayTest and updates it. Returns the server's representation of the networkDelayTest, and an error, if there is any.
func (c *networkDelayTests) Update(networkDelayTest *v1alpha1.NetworkDelayTest) (result *v1alpha1.NetworkDelayTest, err error) {
	result = &v1alpha1.NetworkDelayTest{}
	err = c.client.Put().
		Resource("networkdelaytests").
		Name(networkDelayTest.Name).
		Body(networkDelayTest).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *networkDelayTests) UpdateStatus(networkDelayTest *v1alpha1.NetworkDelayTest) (result *v1alpha1.NetworkDelayTest, err error) {
	result = &v1alpha1.NetworkDelayTest{}
	err = c.client.Put().
		Resource("networkdelaytests").
		Name(networkDelayTest.Name).
		SubResource("status").
		Body(networkDelayTest).
		Do().
		Into(result)
	return
}

// Delete takes name of the networkDelayTest and deletes it. Returns an error if one occurs.
func (c *networkDelayTests) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("networkdelaytests").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *networkDelayTests) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("networkdelaytests").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched networkDelayTest.
func (c *networkDelayTests) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.NetworkDelayTest, err error) {
	result = &v1alpha1.NetworkDelayTest{}
	err = c.client.Patch(pt).
		Resource("networkdelaytests").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}

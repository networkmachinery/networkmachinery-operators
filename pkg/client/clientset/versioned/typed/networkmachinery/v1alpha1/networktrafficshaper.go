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

// NetworkTrafficShapersGetter has a method to return a NetworkTrafficShaperInterface.
// A group's client should implement this interface.
type NetworkTrafficShapersGetter interface {
	NetworkTrafficShapers() NetworkTrafficShaperInterface
}

// NetworkTrafficShaperInterface has methods to work with NetworkTrafficShaper resources.
type NetworkTrafficShaperInterface interface {
	Create(*v1alpha1.NetworkTrafficShaper) (*v1alpha1.NetworkTrafficShaper, error)
	Update(*v1alpha1.NetworkTrafficShaper) (*v1alpha1.NetworkTrafficShaper, error)
	UpdateStatus(*v1alpha1.NetworkTrafficShaper) (*v1alpha1.NetworkTrafficShaper, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.NetworkTrafficShaper, error)
	List(opts v1.ListOptions) (*v1alpha1.NetworkTrafficShaperList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.NetworkTrafficShaper, err error)
	NetworkTrafficShaperExpansion
}

// networkTrafficShapers implements NetworkTrafficShaperInterface
type networkTrafficShapers struct {
	client rest.Interface
}

// newNetworkTrafficShapers returns a NetworkTrafficShapers
func newNetworkTrafficShapers(c *NetworkmachineryV1alpha1Client) *networkTrafficShapers {
	return &networkTrafficShapers{
		client: c.RESTClient(),
	}
}

// Get takes name of the networkTrafficShaper, and returns the corresponding networkTrafficShaper object, and an error if there is any.
func (c *networkTrafficShapers) Get(name string, options v1.GetOptions) (result *v1alpha1.NetworkTrafficShaper, err error) {
	result = &v1alpha1.NetworkTrafficShaper{}
	err = c.client.Get().
		Resource("networktrafficshapers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of NetworkTrafficShapers that match those selectors.
func (c *networkTrafficShapers) List(opts v1.ListOptions) (result *v1alpha1.NetworkTrafficShaperList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.NetworkTrafficShaperList{}
	err = c.client.Get().
		Resource("networktrafficshapers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested networkTrafficShapers.
func (c *networkTrafficShapers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("networktrafficshapers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a networkTrafficShaper and creates it.  Returns the server's representation of the networkTrafficShaper, and an error, if there is any.
func (c *networkTrafficShapers) Create(networkTrafficShaper *v1alpha1.NetworkTrafficShaper) (result *v1alpha1.NetworkTrafficShaper, err error) {
	result = &v1alpha1.NetworkTrafficShaper{}
	err = c.client.Post().
		Resource("networktrafficshapers").
		Body(networkTrafficShaper).
		Do().
		Into(result)
	return
}

// Update takes the representation of a networkTrafficShaper and updates it. Returns the server's representation of the networkTrafficShaper, and an error, if there is any.
func (c *networkTrafficShapers) Update(networkTrafficShaper *v1alpha1.NetworkTrafficShaper) (result *v1alpha1.NetworkTrafficShaper, err error) {
	result = &v1alpha1.NetworkTrafficShaper{}
	err = c.client.Put().
		Resource("networktrafficshapers").
		Name(networkTrafficShaper.Name).
		Body(networkTrafficShaper).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *networkTrafficShapers) UpdateStatus(networkTrafficShaper *v1alpha1.NetworkTrafficShaper) (result *v1alpha1.NetworkTrafficShaper, err error) {
	result = &v1alpha1.NetworkTrafficShaper{}
	err = c.client.Put().
		Resource("networktrafficshapers").
		Name(networkTrafficShaper.Name).
		SubResource("status").
		Body(networkTrafficShaper).
		Do().
		Into(result)
	return
}

// Delete takes name of the networkTrafficShaper and deletes it. Returns an error if one occurs.
func (c *networkTrafficShapers) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("networktrafficshapers").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *networkTrafficShapers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("networktrafficshapers").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched networkTrafficShaper.
func (c *networkTrafficShapers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.NetworkTrafficShaper, err error) {
	result = &v1alpha1.NetworkTrafficShaper{}
	err = c.client.Patch(pt).
		Resource("networktrafficshapers").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}

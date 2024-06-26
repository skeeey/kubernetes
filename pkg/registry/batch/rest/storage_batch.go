/*
Copyright 2016 The Kubernetes Authors.

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

package rest

import (
	batchapiv1 "k8s.io/api/batch/v1"
	batchapiv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/apis/batch"
	cronjobstore "k8s.io/kubernetes/pkg/registry/batch/cronjob/storage"
	jobstore "k8s.io/kubernetes/pkg/registry/batch/job/storage"
)

// skeeey: [kube-apiserver] install default rest apis (storage interface) (3-2) (batch)
type RESTStorageProvider struct{}

func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(batch.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)
	// If you add a version here, be sure to add an entry in `k8s.io/kubernetes/cmd/kube-apiserver/app/aggregator.go with specific priorities.
	// TODO refactor the plumbing to provide the information in the APIGroupInfo

	if storageMap, err := p.v1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
		return genericapiserver.APIGroupInfo{}, err
	} else if len(storageMap) > 0 {
		apiGroupInfo.VersionedResourcesStorageMap[batchapiv1.SchemeGroupVersion.Version] = storageMap
	}

	if storageMap, err := p.v1beta1Storage(apiResourceConfigSource, restOptionsGetter); err != nil {
		return genericapiserver.APIGroupInfo{}, err
	} else if len(storageMap) > 0 {
		apiGroupInfo.VersionedResourcesStorageMap[batchapiv1beta1.SchemeGroupVersion.Version] = storageMap
	}

	return apiGroupInfo, nil
}

func (p RESTStorageProvider) v1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}

	// jobs
	if resource := "jobs"; apiResourceConfigSource.ResourceEnabled(batchapiv1.SchemeGroupVersion.WithResource(resource)) {
		// skeeey: [kube-apiserver] install default rest apis (storage interface) (3-3) (batch/job)
		jobsStorage, jobsStatusStorage, err := jobstore.NewREST(restOptionsGetter)
		if err != nil {
			return storage, err
		}
		storage[resource] = jobsStorage
		storage[resource+"/status"] = jobsStatusStorage
	}

	// cronjobs
	if resource := "cronjobs"; apiResourceConfigSource.ResourceEnabled(batchapiv1.SchemeGroupVersion.WithResource(resource)) {
		cronJobsStorage, cronJobsStatusStorage, err := cronjobstore.NewREST(restOptionsGetter)
		if err != nil {
			return storage, err
		}
		storage[resource] = cronJobsStorage
		storage[resource+"/status"] = cronJobsStatusStorage
	}
	return storage, nil
}

func (p RESTStorageProvider) v1beta1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}

	// cronjobs
	if resource := "cronjobs"; apiResourceConfigSource.ResourceEnabled(batchapiv1beta1.SchemeGroupVersion.WithResource(resource)) {
		cronJobsStorage, cronJobsStatusStorage, err := cronjobstore.NewREST(restOptionsGetter)
		if err != nil {
			return storage, err
		}
		storage[resource] = cronJobsStorage
		storage[resource+"/status"] = cronJobsStatusStorage
	}

	return storage, nil
}

func (p RESTStorageProvider) GroupName() string {
	return batch.GroupName
}

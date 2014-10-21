package authorizetoken

import (
	"fmt"

	kubeapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/apiserver"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"

	"github.com/openshift/origin/pkg/oauth/api"
	//"github.com/openshift/origin/pkg/oauth/api/validation"
)

// REST implements the RESTStorage interface in terms of an Registry.
type REST struct {
	registry Registry
}

// NewStorage returns a new REST.
func NewREST(registry Registry) apiserver.RESTStorage {
	return &REST{registry}
}

// New returns a new AuthorizeToken for use with Create and Update.
func (s *REST) New() runtime.Object {
	return &api.AuthorizeToken{}
}

// Get retrieves an AuthorizeToken by id.
func (s *REST) Get(ctx kubeapi.Context, id string) (runtime.Object, error) {
	token, err := s.registry.GetAuthorizeToken(id)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// List retrieves a list of AuthorizeTokens that match selector.
func (s *REST) List(ctx kubeapi.Context, selector, fields labels.Selector) (runtime.Object, error) {
	tokens, err := s.registry.ListAuthorizeTokens(selector)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Create registers the given AuthorizeToken.
func (s *REST) Create(ctx kubeapi.Context, obj runtime.Object) (<-chan runtime.Object, error) {
	token, ok := obj.(*api.AuthorizeToken)
	if !ok {
		return nil, fmt.Errorf("not an token: %#v", obj)
	}

	token.CreationTimestamp = util.Now()

	// if errs := validation.ValidateAuthorizeToken(token); len(errs) > 0 {
	// 	return nil, errors.NewInvalid("token", token.Name, errs)
	// }

	return apiserver.MakeAsync(func() (runtime.Object, error) {
		if err := s.registry.CreateAuthorizeToken(token); err != nil {
			return nil, err
		}
		return s.Get(ctx, token.Name)
	}), nil
}

// Update is not supported for AuthorizeTokens, as they are immutable.
func (s *REST) Update(ctx kubeapi.Context, obj runtime.Object) (<-chan runtime.Object, error) {
	return nil, fmt.Errorf("AuthorizeTokens may not be changed.")
}

// Delete asynchronously deletes an AuthorizeToken specified by its id.
func (s *REST) Delete(ctx kubeapi.Context, id string) (<-chan runtime.Object, error) {
	return apiserver.MakeAsync(func() (runtime.Object, error) {
		return &kubeapi.Status{Status: kubeapi.StatusSuccess}, s.registry.DeleteAuthorizeToken(id)
	}), nil
}
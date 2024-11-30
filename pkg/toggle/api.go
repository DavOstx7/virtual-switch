package toggle

import "context"

type TogglerAPI struct {
	toggler Toggler
}

func NewTogglerAPI(toggler Toggler) *TogglerAPI {
	return &TogglerAPI{
		toggler: toggler,
	}
}

func (api *TogglerAPI) On(ctx context.Context) error {
	return api.toggler.On(ctx)
}

func (api *TogglerAPI) Off() error {
	return api.toggler.Off()
}

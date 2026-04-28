package dto

type SiteDesign struct {
	Pages       []SitePage     `json:"pages"       validate:"required,min=1,dive"`
	APIRoutes   []SiteAPIRoute `json:"api_routes"  validate:"omitempty,dive"`
	KVSchema    []SiteKVEntry  `json:"kv_schema"   validate:"omitempty,dive"`
	Files       []SiteFilePlan `json:"files"       validate:"required,min=1,dive"`
	Constraints []string       `json:"constraints" validate:"omitempty,dive,required,max=500"`
}

type SitePage struct {
	Path        string `json:"path"        validate:"required,max=128,startswith=/"`
	Title       string `json:"title"       validate:"required,min=1,max=120"`
	Description string `json:"description" validate:"required,min=1,max=500"`
}

type SiteAPIRoute struct {
	Method      string `json:"method"      validate:"required"`
	Path        string `json:"path"        validate:"required,max=128,startswith=/"`
	Description string `json:"description" validate:"required,min=1,max=500"`
	Request     string `json:"request"     validate:"required,min=1,max=2000"`
	Response    string `json:"response"    validate:"required,min=1,max=2000"`
	UsesKV      bool   `json:"uses_kv"`
}

type SiteKVEntry struct {
	Key         string `json:"key"          validate:"required,min=1,max=128"`
	Description string `json:"description"  validate:"required,min=1,max=500"`
	ValueSchema string `json:"value_schema" validate:"required,min=1,max=2000"`
}

type SiteFilePlan struct {
	Path        string `json:"path"        validate:"required,min=1,max=256"`
	Description string `json:"description" validate:"required,min=1,max=500"`
}

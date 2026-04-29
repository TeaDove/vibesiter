package dto

type File struct {
	Path    string `json:"path"    validate:"required,max=128"`
	Content string `json:"content" validate:"required,max=1000000"`
}

type Files struct {
	Files []File `json:"files" validate:"required,min=1,dive"`
}

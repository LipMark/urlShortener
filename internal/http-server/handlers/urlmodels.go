package handlers

// Request describes request to the server. validate is used to check if url is present
type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

// Response describes server response
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

// URLSaver is an interface for saving url by original.
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

// URLGetter is an interface for getting url by alias.
type URLGetter interface {
	GetURL(alias string) (string, error)
}

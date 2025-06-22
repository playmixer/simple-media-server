package rest

type element struct {
	Name        string
	Path        string
	IsDir       bool
	Accessible  bool
	IsAviConver bool
	Size        int64
	Converting  bool
}

type tExplorerResponse struct {
	Current string    `json:"current"`
	List    []element `json:"list"`
}

type tConvertRequest struct {
	From string `json:"from"`
}

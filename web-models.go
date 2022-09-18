package eadmin

type TableResponse struct {
	Rooms []string `json:"rooms"`
	Data  []Entry  `json:"data"`
}

type Entry map[string]string

type ThreeD struct {
	Nodes []ThreeDNode `json:"nodes"`
	Links []ThreeDLink `json:"links"`
}

type ThreeDNode struct {
	ID    string `json:"id"`
	Group int    `json:"group"`
}

type ThreeDLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}

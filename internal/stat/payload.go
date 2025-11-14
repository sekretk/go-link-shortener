package stat

type GetStatisticResponse struct {
	Period string `json:"period"`
	Sum    int    `json:"sum"`
}

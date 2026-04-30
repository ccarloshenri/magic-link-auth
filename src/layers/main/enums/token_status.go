package enums

type TokenStatus string

const (
	Pending TokenStatus = "PENDING"
	Used    TokenStatus = "USED"
	Expired TokenStatus = "EXPIRED"
)

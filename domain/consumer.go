package domain

// Consumer - Consumer domain
type Consumer struct {
	Name      string `bson:"name,omitempty" json:"name"`
	Phone     string `bson:"phone,omitempty" json:"phone"`
	Accesskey string `bson:"accessKey,omitempty" json:"accessKey"`
	Status    string `bson:"status,omitempty" json:"status"`
}

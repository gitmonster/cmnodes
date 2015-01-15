package nodes

type NodeContent struct {
	Content     interface{}      `bson:"c"`

	DisplayReps []Representation `bson:"dr"`
}

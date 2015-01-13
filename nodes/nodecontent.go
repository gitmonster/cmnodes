package nodes

type NodeContent struct {
	Content     interface{}      `bson:"c"`
	EditRep     Representation   `bson:"er"`
	DisplayReps []Representation `bson:"dr"`
}

package enums

type Network string

const (
	NetworkMAIN   Network = "https://api.trongrid.io"
	NetworkSHASTA Network = "https://api.shasta.trongrid.io"
	NetworkNILE   Network = "https://nile.trongrid.io"
)

type Node string

const (
	MAIN_NODE   Node = "grpc.trongrid.io:50051"
	SHASTA_NODE Node = "grpc.shasta.trongrid.io:50051"
)

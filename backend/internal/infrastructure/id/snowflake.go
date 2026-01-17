package id

import "github.com/bwmarrin/snowflake"

type SnowflakeGenerator struct {
	node *snowflake.Node
}

func NewSnowFlakeGenerator(nodeID int64) (*SnowflakeGenerator, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}

	return &SnowflakeGenerator{node: node}, nil
}

func (gen *SnowflakeGenerator) NextID() (int64, error) {
	return gen.node.Generate().Int64(), nil
}

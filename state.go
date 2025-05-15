package main

import (
	"github.com/VuTLy/blogAggregator/internal/config"
	"github.com/VuTLy/blogAggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

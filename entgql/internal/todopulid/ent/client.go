// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"log"

	"github.com/facebookincubator/ent-contrib/entgql/internal/todopulid/ent/migrate"
	"github.com/facebookincubator/ent-contrib/entgql/internal/todopulid/ent/schema/pulid"

	"github.com/facebookincubator/ent-contrib/entgql/internal/todopulid/ent/todo"

	"github.com/facebook/ent/dialect"
	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Todo is the client for interacting with the Todo builders.
	Todo *TodoClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Todo = NewTodoClient(c.config)
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %v", err)
	}
	cfg := config{driver: tx, log: c.log, debug: c.debug, hooks: c.hooks}
	return &Tx{
		ctx:    ctx,
		config: cfg,
		Todo:   NewTodoClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(*sql.Driver).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %v", err)
	}
	cfg := config{driver: &txDriver{tx: tx, drv: c.driver}, log: c.log, debug: c.debug, hooks: c.hooks}
	return &Tx{
		config: cfg,
		Todo:   NewTodoClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Todo.
//		Query().
//		Count(ctx)
//
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := config{driver: dialect.Debug(c.driver, c.log), log: c.log, debug: true, hooks: c.hooks}
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Todo.Use(hooks...)
}

// TodoClient is a client for the Todo schema.
type TodoClient struct {
	config
}

// NewTodoClient returns a client for the Todo from the given config.
func NewTodoClient(c config) *TodoClient {
	return &TodoClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `todo.Hooks(f(g(h())))`.
func (c *TodoClient) Use(hooks ...Hook) {
	c.hooks.Todo = append(c.hooks.Todo, hooks...)
}

// Create returns a create builder for Todo.
func (c *TodoClient) Create() *TodoCreate {
	mutation := newTodoMutation(c.config, OpCreate)
	return &TodoCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Todo entities.
func (c *TodoClient) CreateBulk(builders ...*TodoCreate) *TodoCreateBulk {
	return &TodoCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Todo.
func (c *TodoClient) Update() *TodoUpdate {
	mutation := newTodoMutation(c.config, OpUpdate)
	return &TodoUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *TodoClient) UpdateOne(t *Todo) *TodoUpdateOne {
	mutation := newTodoMutation(c.config, OpUpdateOne, withTodo(t))
	return &TodoUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *TodoClient) UpdateOneID(id pulid.ID) *TodoUpdateOne {
	mutation := newTodoMutation(c.config, OpUpdateOne, withTodoID(id))
	return &TodoUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Todo.
func (c *TodoClient) Delete() *TodoDelete {
	mutation := newTodoMutation(c.config, OpDelete)
	return &TodoDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *TodoClient) DeleteOne(t *Todo) *TodoDeleteOne {
	return c.DeleteOneID(t.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *TodoClient) DeleteOneID(id pulid.ID) *TodoDeleteOne {
	builder := c.Delete().Where(todo.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &TodoDeleteOne{builder}
}

// Query returns a query builder for Todo.
func (c *TodoClient) Query() *TodoQuery {
	return &TodoQuery{config: c.config}
}

// Get returns a Todo entity by its id.
func (c *TodoClient) Get(ctx context.Context, id pulid.ID) (*Todo, error) {
	return c.Query().Where(todo.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TodoClient) GetX(ctx context.Context, id pulid.ID) *Todo {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryParent queries the parent edge of a Todo.
func (c *TodoClient) QueryParent(t *Todo) *TodoQuery {
	query := &TodoQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := t.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(todo.Table, todo.FieldID, id),
			sqlgraph.To(todo.Table, todo.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, todo.ParentTable, todo.ParentColumn),
		)
		fromV = sqlgraph.Neighbors(t.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryChildren queries the children edge of a Todo.
func (c *TodoClient) QueryChildren(t *Todo) *TodoQuery {
	query := &TodoQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := t.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(todo.Table, todo.FieldID, id),
			sqlgraph.To(todo.Table, todo.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, todo.ChildrenTable, todo.ChildrenColumn),
		)
		fromV = sqlgraph.Neighbors(t.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *TodoClient) Hooks() []Hook {
	return c.hooks.Todo
}

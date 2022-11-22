//go:build wasm
// +build wasm

package main

import (
	"context"
	"encoding/json"
	"syscall/js"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/Permify/permify/internal/commands"
	"github.com/Permify/permify/pkg/database"
	"github.com/Permify/permify/pkg/development"
	"github.com/Permify/permify/pkg/development/graph"
	v1 "github.com/Permify/permify/pkg/pb/base/v1"
)

// Requests for Permify Playground

var dev *development.Container

// check - Permission check request
func check() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		params := &v1.PermissionCheckRequest{}
		err := protojson.Unmarshal([]byte(string(args[0].String())), params)
		if err != nil {
			return js.ValueOf([]interface{}{false, err.Error()})
		}
		var result commands.CheckResponse
		result, err = development.Check(context.Background(), dev.P, params.Subject, params.Action, params.Entity, string(args[1].String()), "")
		if err != nil {
			return js.ValueOf([]interface{}{false, err.Error()})
		}
		return js.ValueOf([]interface{}{result.Can, nil})
	})
}

// lookupQuery - Schema lookup request
func lookupQuery() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		params := &v1.PermissionLookupQueryRequest{}
		err := protojson.Unmarshal([]byte(string(args[0].String())), params)
		if err != nil {
			return js.ValueOf([]interface{}{"", err.Error()})
		}
		var result commands.LookupQueryResponse
		result, err = development.LookupQuery(context.Background(), dev.P, params.EntityType, params.Action, params.Subject, string(args[1].String()))
		if err != nil {
			return js.ValueOf([]interface{}{"", err.Error()})
		}
		return js.ValueOf([]interface{}{result.Query, nil})
	})
}

// writeSchema - Writes schema
func writeSchema() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		version, err := development.WriteSchema(context.Background(), dev.S, string(args[0].String()))
		if err != nil {
			return js.ValueOf([]interface{}{"", err.Error()})
		}
		return js.ValueOf([]interface{}{version, nil})
	})
}

// writeTuple - Writes relation tuples
func writeTuple() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		t := &v1.Tuple{}
		err := protojson.Unmarshal([]byte(string(args[0].String())), t)
		if err != nil {
			return js.ValueOf([]interface{}{err.Error()})
		}
		_, err = development.WriteTuple(context.Background(), dev.R, []*v1.Tuple{t}, string(args[1].String()))
		if err != nil {
			return js.ValueOf([]interface{}{err.Error()})
		}
		return js.ValueOf([]interface{}{nil})
	})
}

// deleteTuple - Delete relation tuple
func deleteTuple() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		t := &v1.Tuple{}
		err := protojson.Unmarshal([]byte(string(args[0].String())), t)
		if err != nil {
			return js.ValueOf([]interface{}{err.Error()})
		}
		_, err = development.DeleteTuple(context.Background(), dev.R, &v1.TupleFilter{
			Entity: &v1.EntityFilter{
				Type: t.GetEntity().GetType(),
				Ids:  []string{t.GetEntity().GetId()},
			},
			Relation: t.GetRelation(),
			Subject: &v1.SubjectFilter{
				Type:     t.GetSubject().GetType(),
				Ids:      []string{t.GetSubject().GetId()},
				Relation: t.GetSubject().GetRelation(),
			},
		})
		if err != nil {
			return js.ValueOf([]interface{}{err.Error()})
		}
		return js.ValueOf([]interface{}{nil})
	})
}

// readSchema - Read Permify Schema
func readSchema() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		sch, err := development.ReadSchema(context.Background(), dev.S, string(args[0].String()))
		if err != nil {
			return js.ValueOf([]interface{}{nil, err.Error()})
		}
		result, err := protojson.Marshal(sch)
		if err != nil {
			return js.ValueOf([]interface{}{nil, err.Error()})
		}
		return js.ValueOf([]interface{}{string(result), nil})
	})
}

// readTuple - Read, filter relation tuples
func readTuple() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		params := &v1.TupleFilter{}
		err := protojson.Unmarshal([]byte(string(args[0].String())), params)
		if err != nil {
			return js.ValueOf([]interface{}{nil, err.Error()})
		}
		var collection database.ITupleCollection
		collection, err = development.ReadTuple(context.Background(), dev.R, params, "")
		if err != nil {
			return js.ValueOf([]interface{}{nil, err.Error()})
		}
		var result []byte
		t := &v1.Tuples{
			Tuples: collection.GetTuples(),
		}
		result, err = protojson.Marshal(t)
		if err != nil {
			return js.ValueOf([]interface{}{nil, err.Error()})
		}
		return js.ValueOf([]interface{}{string(result), nil})
	})
}

// readSchemaGraph - read schema graph
func readSchemaGraph() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		sch, err := development.ReadSchema(context.Background(), dev.S, string(args[0].String()))
		if err != nil {
			return js.ValueOf([]interface{}{nil, err.Error()})
		}
		r, err := graph.SchemaToGraph(sch)
		if err != nil {
			return js.ValueOf([]interface{}{nil, err.Error()})
		}
		result, err := json.Marshal(struct {
			Nodes []*graph.Node `json:"nodes"`
			Edges []*graph.Edge `json:"edges"`
		}{Nodes: r.Nodes(), Edges: r.Edges()})
		if err != nil {
			return js.ValueOf([]interface{}{nil, err.Error()})
		}
		return js.ValueOf([]interface{}{string(result), nil})
	})
}

func main() {
	ch := make(chan struct{}, 0)
	dev = development.NewContainer()
	js.Global().Set("check", check())
	js.Global().Set("writeSchema", writeSchema())
	js.Global().Set("writeTuple", writeTuple())
	js.Global().Set("readSchema", readSchema())
	js.Global().Set("readTuple", readTuple())
	js.Global().Set("deleteTuple", deleteTuple())
	js.Global().Set("readSchemaGraph", readSchemaGraph())
	js.Global().Set("lookupQuery", lookupQuery())
	<-ch
}

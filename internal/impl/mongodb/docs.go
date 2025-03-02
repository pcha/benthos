package mongodb

import (
	"github.com/Jeffail/benthos/v3/internal/docs"
	"github.com/Jeffail/benthos/v3/internal/impl/mongodb/client"
	"github.com/Jeffail/benthos/v3/public/service"
)

func processorOperationDocs(defaultOperation client.Operation) docs.FieldSpec {
	fs := outputOperationDocs(defaultOperation)
	return fs.HasOptions(append(fs.Options, string(client.OperationFindOne))...)
}

func outputOperationDocs(defaultOperation client.Operation) docs.FieldSpec {
	return docs.FieldCommon(
		"operation",
		"The mongodb operation to perform.",
	).HasOptions(
		string(client.OperationInsertOne),
		string(client.OperationDeleteOne),
		string(client.OperationDeleteMany),
		string(client.OperationReplaceOne),
		string(client.OperationUpdateOne),
	).HasDefault(defaultOperation)
}

func writeConcernDocs() docs.FieldSpecs {
	return docs.FieldSpecs{
		docs.FieldCommon("w", "W requests acknowledgement that write operations propagate to the specified number of mongodb instances."),
		docs.FieldCommon("j", "J requests acknowledgement from MongoDB that write operations are written to the journal."),
		docs.FieldCommon("w_timeout", "The write concern timeout."),
	}
}

func mapExamples() []interface{} {
	examples := []interface{}{"root.a = this.foo\nroot.b = this.bar"}
	return examples
}

var urlField = service.NewStringField("url").
	Description("The URL of the target MongoDB DB.").
	Example("mongodb://localhost:27017")

var queryField = service.NewBloblangField("query").Description("Bloblang expression describing MongoDB query.").Example(`
      root.from = {"$lte": timestamp_unix()}
      root.to = {"$gte": timestamp_unix()}
`)

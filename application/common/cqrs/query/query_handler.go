/*
Package query provides a generic interface for handling queries.

It includes the `IQueryHandler` interface for processing queries of any type and returning results of any type.
*/
package query

// IQueryHandler defines a generic interface for handling queries.
//
// Type Parameters:
// - Query: The type of the query to be handled.
// - Result: The type of the result returned after handling the query.
type IQueryHandler[Query any, Result any] interface {

	// Handle processes the provided query and returns the result or an error.
	Handle(query Query) (Result, error)
}

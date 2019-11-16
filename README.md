# GraphQL Client

This package provides a simple GraphQL client with automatic mapping of Go structs to GraphQL queries and mutations.


## Overview

This is initial version of the library and it does not support all cool GraphQL features. 

It allows for both passing queries as string or using automatic mapping from Go variables.
Bare in mind that automatic mapping impacts the performance as it uses reflection heavily.


## Usage

### Queries and Mutations with automatic mapping

The automatic mapping is done via reflection. It uses `json` tags for fields names.

For the following schema:
```graphql
type Dog {
    id: ID!
    name: String!
}

input DogInput {
    name: String!
}

type Query {
    dogs: [Dog]
    dog(id: ID!): Dog!
}

type Mutation {
    createDog(in: DogInput!): Dog
}
```

The client can be used with automatic mapping as presented:

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/szymongib/graphql-client/graphql"
)

type Dog struct {
    Id   string `json:"id"`
    Name string `json:"name"`
}

type DogInput struct {
    Name string `json:"name"`
}

func main(){
    apiAddress := "http://localhost:8000/graphql"

    gqlClient := graphql.NewClient(apiAddress)
    
    dogInput := DogInput{
        Name: "Rex",
    }
    
    var dog Dog
    err := gqlClient.Mutate(context.Background(), "createDog", graphql.OperationInput{"in": dogInput}, &dog)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    var allDogs []Dog
    err = gqlClient.Query(context.Background(), "dogs", nil, &allDogs)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    // ...
}
```


### Mapping structs to Graphql

Mapping structs can be used without the client.

The following code:
```go
type MyStruct struct {
    StringData string          `json:"stringData"`
    Inner      MyInnerStruct   `json:"inner"`
    InnerSlice []MyInnerStruct `json:"innerSlice"`
}

type MyInnerStruct struct {
    InnerStringData string `json:"innerStringData"`
    IntData         int    `json:"intData"`
}

func main() {
    query := graphql.ParseToGQLQuery(MyStruct{})
    fmt.Println(query)
}
```

Will result in query as such:
```
{
	stringData 
	inner {
		innerStringData  
		intData  
	}
	innerSlice {
		innerStringData  
        intData 
	}
}
```


## Summary

Any feedback highly appreciated!

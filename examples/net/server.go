// Copyright 2015 LinkedIn Corp. Licensed under the Apache License,
// Version 2.0 (the "License"); you may not use this file except in
// compliance with the License.  You may obtain a copy of the License
// at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.Copyright [201X] LinkedIn Corp. Licensed under the Apache
// License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License.  You may obtain a copy of
// the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.

package main

import (
	"github.com/linkedin/goavro"
	"log"
	"net"
)

const recordSchema = `
{
  "type": "record",
  "name": "comments",
  "doc:": "A basic schema for storing blog comments",
  "namespace": "com.example",
  "fields": [
    {
      "doc": "Name of user",
      "type": "string",
      "name": "username"
    },
    {
      "doc": "The content of the user's message",
      "type": "string",
      "name": "comment"
    },
    {
      "doc": "Unix epoch time in milliseconds",
      "type": "long",
      "name": "timestamp"
    }
  ]
}
`

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serveClient(conn)
	}
}

func serveClient(conn net.Conn) {
	fw, err := goavro.NewWriter(
		goavro.Compression(goavro.CompressionSnappy),
		goavro.WriterSchema(recordSchema),
		goavro.ToWriter(conn))
	if err != nil {
		log.Fatal(err)
	}
	defer fw.Close()

	// create a record that matches the schema we want to encode
	someRecord, err := goavro.NewRecord(goavro.RecordSchema(recordSchema))
	if err != nil {
		log.Fatal(err)
	}
	// identify field name to set datum for
	someRecord.Set("username", "Aquaman")
	someRecord.Set("comment", "The Atlantic is oddly cold this morning!")
	// you can fully qualify the field name
	someRecord.Set("com.example.timestamp", int64(1082196484))
	fw.Write(someRecord)

	// create another record
	if someRecord, err = goavro.NewRecord(goavro.RecordSchema(recordSchema)); err != nil {
		log.Fatal(err)
	}
	someRecord.Set("username", "Batman")
	someRecord.Set("comment", "Who are all of these crazies?")
	someRecord.Set("com.example.timestamp", int64(1427383430))
	fw.Write(someRecord)
}

package proxy

// pg_parser.go would contain a more robust implementation of the PostgreSQL wire protocol
// For this MVP starter, the parsing logic is currently embedded within the Sniffer in server.go
// to keep the code footprint minimal. A complete version would use pgproto3 or custom state machines.

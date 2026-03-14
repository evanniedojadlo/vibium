package bidi

import "net/http"

// ConnectRemote connects to a remote BiDi endpoint, creates a client,
// and establishes a session. Returns the connection, client, and session ID.
func ConnectRemote(url string, headers http.Header) (*Connection, *Client, string, error) {
	conn, err := ConnectWithHeaders(url, headers)
	if err != nil {
		return nil, nil, "", err
	}

	client := NewClient(conn)

	result, err := client.SessionNew(map[string]interface{}{})
	if err != nil {
		conn.Close()
		return nil, nil, "", err
	}

	return conn, client, result.SessionID, nil
}

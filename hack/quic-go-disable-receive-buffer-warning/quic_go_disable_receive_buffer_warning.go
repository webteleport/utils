package hack

import "os"

func init() {
	os.Setenv("QUIC_GO_DISABLE_RECEIVE_BUFFER_WARNING", "true")
}

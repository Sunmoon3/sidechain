package musig

import "testing"

/**
 * @description TODO
 * @author qfl
 * @date 2023/6/14 19:17
 * @version 1.0
 */

func TestEach(t *testing.T) {
	bytes := [][]byte{
		[]byte{1, 2, 3},
		[]byte{4, 5, 6},
	}
	Each(bytes)
}

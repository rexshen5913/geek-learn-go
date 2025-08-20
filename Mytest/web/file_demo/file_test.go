package filedemo

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFile(t *testing.T) {

	fmt.Println(os.Getwd())

	// 這樣打開的話，只有讀取的權限
	f, err := os.Open("testdata/my_file.txt")
	require.NoError(t, err)

	bs := make([]byte, 64)
	n, err := f.Read(bs)
	require.NoError(t, err)
	fmt.Println(n)
	fmt.Println(string(bs))
	f.Close()

	// 這樣打開的話，有讀取和寫入的權限
	f, err = os.OpenFile("testdata/my_file.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	require.NoError(t, err)
	f.WriteString("Hello, World!")
	f.Close()

	f, err = os.Create("testdata/my_file_create.txt")
	require.NoError(t, err)
	f.WriteString("Hello, World!")
	f.Close()

}

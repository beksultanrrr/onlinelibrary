package data
import ( 
	"fmt"
	"strconv"
)
// Declare a custom Runtime type, which has the underlying type int32 (the same as our // book struct field).
type Readtime int32
// Implement a MarshalJSON() method on the Runtime type so that it satisfies the
// json.Marshaler interface. This should return the JSON-encoded value for the book // runtime (in our case, it will return a string in the format "<runtime> mins").
 func (r Readtime) MarshalJSON() ([]byte, error) {
// Generate a string containing the book runtime in the required format.
	jsonValue := fmt.Sprintf("%d mins", r)
	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil 

}
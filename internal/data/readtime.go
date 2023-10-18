package data
import ( 
	"errors"
	"fmt"
	"strconv"
	"strings"

)
var ErrInvalidReadtimeFormat = errors.New("invalid readtime format")
type Readtime int32
// Implement a MarshalJSON() method on the Readtime type so that it satisfies the
// json.Marshaler interface. This should return the JSON-encoded value for the book // Readtime (in our case, it will return a string in the format "<Readtime> mins").
 func (r Readtime) MarshalJSON() ([]byte, error) {
// Generate a string containing the book Readtime in the required format.
	jsonValue := fmt.Sprintf("%d mins", r)
	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil 

}

func (r *Readtime) UnmarshalJSON(jsonValue []byte) error {
	// We expect that the incoming JSON value will be a string in the format
	// "<Readtime> mins", and the first thing we need to do is remove the surrounding // double-quotes from this string. If we can't unquote it, then we return the
	// ErrInvalidReadtimeFormat error.
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
	return ErrInvalidReadtimeFormat }
	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedJSONValue, " ")
	// Sanity check the parts of the string to make sure it was in the expected format. // If it isn't, we return the ErrInvalidReadtimeFormat error again.
	if len(parts) != 2 || parts[1] != "mins" {
	return ErrInvalidReadtimeFormat }
	// Otherwise, parse the string containing the number into an int32. Again, if this // fails return the ErrInvalidReadtimeFormat error.
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
	return ErrInvalidReadtimeFormat }
	// Convert the int32 to a Readtime type and assign this to the receiver. Note that we // use the * operator to deference the receiver (which is a pointer to a Readtime
	// type) in order to set the underlying value of the pointer.
	*r = Readtime(i)
	return nil
	}
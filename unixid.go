package unixid

import (
	"reflect"
	"strconv"
	"unsafe"
)

func (id *UnixID) setValue(rv *reflect.Value, valueOut *string, sizeOut []byte) error {

	*valueOut = id.unixIdNano()

	size := uint8(len(*valueOut))

	sizeOut = append(sizeOut, size)

	// agregamos el id al campo de la estructura origen
	rv.SetString(*valueOut)

	return nil

}

func (id *UnixID) unixIdNano() string {

	currentUnixNano := id.timeNano.UnixNano()

	if currentUnixNano == id.lastUnixNano {
		//mientras sean iguales sumar numero correlativo
		id.correlativeNumber++
	} else {
		id.correlativeNumber = 0
	}
	// actualizo la variable unix nano
	id.lastUnixNano = currentUnixNano

	currentUnixNano += id.correlativeNumber

	return strconv.FormatInt(currentUnixNano, 10)

}

func (id *UnixID) unixIdNanoLAB() string {

	currentUnixNano := id.timeNano.UnixNano()

	if currentUnixNano == id.lastUnixNano {
		//mientras sean iguales sumar numero correlativo
		id.correlativeNumber++
	} else {
		id.correlativeNumber = 0
	}
	// actualizo la variable unix nano
	id.lastUnixNano = currentUnixNano

	currentUnixNano += id.correlativeNumber

	id.buf = id.buf[:0]

	// fmt.Println("size buffer:", sizeBuf)

	id.buf = strconv.AppendInt(id.buf, currentUnixNano, 10)
	// id.buf = strconv.AppendUint(id.buf, currentUnixNano, 10)

	// fmt.Println("tmpBuf:", id.buf, "size id buffer:", len(id.buf))

	// return string(id.buf)
	return *(*string)(unsafe.Pointer(&id.buf))
	// return unsafe.String(unsafe.SliceData(id.buf), len(id.buf))
}

// v := int64(42)
// b := unsafe.Slice((*byte)(unsafe.Pointer(&v)), unsafe.Sizeof(v))
// fmt.Println(b, "id:", string(b))

// b := *(*[]byte)(unsafe.Pointer(&v)) // Cast directly to a []byte pointer
// fmt.Println(b, "id:", string(b))    // Output: [42 0 0 0 0 0 0 0] id: *

// buf := unsafe.Slice((*byte)(unsafe.Pointer(&currentUnixNano)), unsafe.Sizeof(currentUnixNano))

// fmt.Println("id:", string(buf))

// return unsafe.String(unsafe.SliceData(buf), len(buf))
// t := time.Now().UTC().UnixNano()
//     b := unsafe.Slice((*byte)(unsafe.Pointer(&t)), unsafe.Sizeof(t))
//     fmt.Println(b)

// https://stackoverflow.com/questions/76431857/what-is-the-fastest-way-to-convert-int64-to-byte-array

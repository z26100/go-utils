package time

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
	"time"
)

var (
	DateTimeFormat = time.RFC3339
)

type TimeCodec struct {
}

func (t TimeCodec) EncodeValue(context bsoncodec.EncodeContext, writer bsonrw.ValueWriter, value reflect.Value) error {
	dat := value.Interface().(primitive.DateTime).Time().UTC().Format(DateTimeFormat)
	writer.WriteString(string(dat))
	return nil
}

func (t TimeCodec) DecodeValue(context bsoncodec.DecodeContext, reader bsonrw.ValueReader, value reflect.Value) error {
	ts, _ := reader.ReadString()
	ti, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(ti))

	return nil
}

type ObjectIDCodec struct {
}

func (o ObjectIDCodec) EncodeValue(context bsoncodec.EncodeContext, writer bsonrw.ValueWriter, value reflect.Value) error {
	if value.Kind() != reflect.ValueOf(primitive.ObjectID{}).Kind() {
		return errors.New("invalid object id")
	}
	writer.WriteString(formatHex(value.Interface().(primitive.ObjectID).Hex()))
	return nil
}

func formatHex(hex string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		hex[0:6],
		hex[6:12],
		hex[12:18],
		hex[18:24])
}

func decodeHex(hex string) string {
	return strings.ReplaceAll(hex, "-", "")
}

func (o ObjectIDCodec) DecodeValue(context bsoncodec.DecodeContext, reader bsonrw.ValueReader, value reflect.Value) error {
	hex, err := reader.ReadString()
	if err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(decodeHex(hex))
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(id))
	return nil
}
func CreateCustomRegistry() *bsoncodec.RegistryBuilder {
	var primitiveCodecs bson.PrimitiveCodecs
	rb := bsoncodec.NewRegistryBuilder()
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	rb.RegisterCodec(reflect.TypeOf(primitive.ObjectID{}), ObjectIDCodec{})
	rb.RegisterCodec(reflect.TypeOf(primitive.DateTime(0)), TimeCodec{})
	primitiveCodecs.RegisterPrimitiveCodecs(rb)
	return rb
}

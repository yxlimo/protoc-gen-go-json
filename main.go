package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/yxlimo/protoc-gen-go-json/v2/gen"
)

var (
	multiline       = flag.Bool("multiline", false, "generate multiline json")
	useEnumNumbers  = flag.Bool("use_enum_numbers", false, "render enums as integers as opposed to strings")
	emitUnpopulated = flag.Bool("emit_unpopulated", false, "render fields with zero values")
	userProtoNames  = flag.Bool("use_proto_names", false, "use original (.proto) name for fields")
	allowPartial    = flag.Bool("allow_partial", false, "allow partial results")
	discardUnknown  = flag.Bool("discard_unknown", true, "allow messages to contain unknown fields when unmarshaling")

	sqlSupport = flag.Bool("sql_support", false, "generate sql.Scanner and driver.Valuer method")
)

func main() {
	flag.Parse()
	defer glog.Flush()

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gp *protogen.Plugin) error {

		opts := gen.Options{
			Multiline:       *multiline,
			UseEnumNumbers:  *useEnumNumbers,
			EmitUnpopulated: *emitUnpopulated,
			UseProtoNames:   *userProtoNames,
			AllowPartial:    *allowPartial,
			DiscardUnknown:  *discardUnknown,
			SqlSupport:      *sqlSupport,
		}

		for _, name := range gp.Request.FileToGenerate {
			f := gp.FilesByPath[name]

			if len(f.Messages) == 0 {
				glog.V(1).Infof("Skipping %s, no messages", name)
				continue
			}

			glog.V(1).Infof("Processing %s", name)
			glog.V(2).Infof("Generating %s\n", fmt.Sprintf("%s.pb.json.go", f.GeneratedFilenamePrefix))

			gf := gp.NewGeneratedFile(fmt.Sprintf("%s.pb.json.go", f.GeneratedFilenamePrefix), f.GoImportPath)

			err := gen.ApplyTemplate(gf, f, opts)
			if err != nil {
				gf.Skip()
				gp.Error(err)
				continue
			}
		}

		return nil
	})
}

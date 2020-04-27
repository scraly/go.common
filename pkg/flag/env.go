package flag

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/scraly/go.common/pkg/log"
	"github.com/fatih/structs"
)

// AsEnvVariables sets struct values from environment variables
func AsEnvVariables(o interface{}, prefix string, skipCommented bool) map[string]string {
	r := map[string]string{}
	prefix = strings.ToUpper(prefix)
	delim := "_"
	if prefix == "" {
		delim = ""
	}
	fields := structs.Fields(o)
	for _, f := range fields {
		if skipCommented {
			tag := f.Tag("commented")
			if tag != "" {
				commented, err := strconv.ParseBool(tag)
				log.CheckErr("Unable to parse tag value", err)
				if commented {
					continue
				}
			}
		}
		if structs.IsStruct(f.Value()) {
			rf := AsEnvVariables(f.Value(), prefix+delim+f.Name(), skipCommented)
			for k, v := range rf {
				r[k] = v
			}
		} else {
			r[prefix+"_"+strings.ToUpper(f.Name())] = fmt.Sprintf("%v", f.Value())
		}
	}
	return r
}

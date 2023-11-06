package xml

import (
	"github.com/viant/tagly/format"
	"strconv"
	"strings"
)

type (

	//ColumnHeader represents column header
	ColumnHeader struct {
		ID   string `json:",omitempty" xmlify:"path=@id"`
		Type string `json:",omitempty" xmlify:"path=@type"`
	}

	//ColumnHolder represents a column holder
	ColumnHolder struct {
		Columns []*ColumnHeader `xmlify:"name=column"`
	}

	//Column represents a column value
	Column struct {
		LongType   *string `json:",omitempty" xmlify:"omitempty,path=@lg"`
		IntType    *string `json:",omitempty" xmlify:"omitempty,path=@long"`
		DoubleType *string `json:",omitempty" xmlify:"omitempty,path=@db"`
		DateType   *string `json:",omitempty" xmlify:"omitempty,path=@ts"`
		Value      *string `json:",omitempty" xmlify:"omitempty,omittagname"`
		ValueAttr  *string `json:",omitempty" xmlify:"omitempty,path=@nil"`
	}

	//Row defines tabular row
	Row struct {
		Columns []*Column `xmlify:"name=c"`
	}

	//RowsHolder represents rows holder
	RowsHolder struct {
		Rows []*Row `xmlify:"name=r"`
	}

	//Tabular represents tabular holder
	Tabular struct {
		ColumnsWrapper ColumnHolder `xmlify:"name=columns"`
		RowsWrapper    RowsHolder   `xmlify:"name=rows"`
	}

	//Filter represents filter
	Filter struct {
		Name          string
		Tag           format.Tag
		IncludeString []string `json:",omitempty" xmlify:"omitempty"`
		ExcludeString []string `json:",omitempty" xmlify:"omitempty"`
		IncludeInt    []int    `json:",omitempty" xmlify:"omitempty"`
		ExcludeInt    []int    `json:",omitempty" xmlify:"omitempty"`
		IncludeBool   []bool   `json:",omitempty" xmlify:"omitempty"`
		ExcludeBool   []bool   `json:",omitempty" xmlify:"omitempty"`
	}

	//FilterHolder represents filter holder
	FilterHolder struct {
		Filters []*Filter
	}
)

// MarshalXML customizes XML marshaling
func (f *FilterHolder) MarshalXML() ([]byte, error) {

	var sb strings.Builder

	sb.WriteString("<filter>")

	for _, filter := range f.Filters {
		wasInclusion := false

		//TODO check if filter is empty
		sb.WriteString("\n")
		sb.WriteString("<")

		if filter.Tag.Name != "" {
			sb.WriteString(filter.Tag.Name)
		} else {
			sb.WriteString(filter.Name)
		}

		sb.WriteString(" ")

		switch {
		case filter.IncludeInt != nil:
			sb.WriteString(`include-ids="`)
			for i, value := range filter.IncludeInt {
				if i > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(strconv.Itoa(value))
			}
			sb.WriteString(`"`)
			wasInclusion = true
		case filter.IncludeString != nil:
			sb.WriteString(`include-ids="`)
			for i, value := range filter.IncludeString {
				if i > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(value)
			}
			sb.WriteString(`"`)
			wasInclusion = true
		case filter.IncludeBool != nil:
			sb.WriteString(`include-ids="`)
			for i, value := range filter.IncludeBool {
				if i > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(strconv.FormatBool(value))
			}
			sb.WriteString(`"`)
			wasInclusion = true
		}

		isExclusion := filter.ExcludeInt != nil || filter.ExcludeString != nil || filter.ExcludeBool != nil
		if wasInclusion && isExclusion {
			sb.WriteString(` `)
		}

		switch {
		case filter.ExcludeInt != nil:
			sb.WriteString(`exclude-ids="`)
			for i, value := range filter.ExcludeInt {
				if i > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(strconv.Itoa(value))
			}
			sb.WriteString(`"`)
		case filter.ExcludeString != nil:
			sb.WriteString(`exclude-ids="`)
			for i, value := range filter.ExcludeString {
				if i > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(value)
			}
			sb.WriteString(`"`)
		case filter.ExcludeBool != nil:
			sb.WriteString(`exclude-ids="`)
			for i, value := range filter.ExcludeBool {
				if i > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(strconv.FormatBool(value))
			}
			sb.WriteString(`"`)
		}

		sb.WriteString("/>")
	}

	sb.WriteString("\n")
	sb.WriteString("</filter>")
	//fmt.Println(sb.String())
	return []byte(sb.String()), nil
}

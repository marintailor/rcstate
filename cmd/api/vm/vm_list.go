package vm

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/marintailor/rcstate/cmd/api/gce"
)

// List returns a table formatted list of instances.
func (vm *VirtualMachine) List() string {
	if _, err := vm.Instances.GetInstancesList(); err != nil {
		fmt.Println("list: get instances list:", err)
	}

	columnNames := []string{"Name", "Status", "Internal", "External", "Type", "Preemptible"}
	columnWidth := getColumnsWidth(columnNames, vm.Instances.List)

	var out bytes.Buffer

	tableHeader(&out, columnWidth, vm.Project, vm.Zone)
	tableColumns(&out, columnNames, columnWidth)
	tableRows(&out, vm.Instances.List, columnNames, columnWidth)

	return out.String()
}

// tableHeader writes to a writer the header for the table.
func tableHeader(w io.Writer, columnWidths []int, project string, zone string) {
	line := ""
	for _, w := range columnWidths {
		line = line + strings.Repeat("=", w+2)
	}

	fmt.Fprintf(w, "\n%s\nPROJECT: %s    ZONE: %s\n%s\n\n", line, project, zone, line)
}

// tableColumns writes to a writer the column names of the table.
func tableColumns(w io.Writer, columnNames []string, columnWidths []int) {
	for i, c := range columnNames {
		fmt.Fprintf(w, "%s%s  ", c, strings.Repeat(" ", columnWidths[i]-len(c)))
	}
	fmt.Fprintln(w)

	for i := range columnNames {
		fmt.Fprintf(w, "%s  ", strings.Repeat("-", columnWidths[i]))
	}
	fmt.Fprintln(w)
}

// tableRows writes to a writer a row of details for each instance.
func tableRows(w io.Writer, list []gce.Instance, columnNames []string, columnWidths []int) {
	for _, inst := range list {
		row := ""
		v := reflect.ValueOf(inst)
		for j := 0; j < len(columnNames); j++ {
			switch v.FieldByName(columnNames[j]).Kind() {
			case reflect.Bool:
				f := v.FieldByName(columnNames[j]).Bool()
				row += fmt.Sprintf("%t  ", f)
			case reflect.String:
				f := v.FieldByName(columnNames[j]).String()
				row += fmt.Sprintf("%s%s  ", f, strings.Repeat(" ", columnWidths[j]-len(f)))
			}
		}

		fmt.Fprintf(w, "%s\n", row)
	}
}

// getColumnsWidth returns a slice of widths for each table column based on values of instance's details.
func getColumnsWidth(columnNames []string, list []gce.Instance) []int {
	columnWidths := make([]int, len(columnNames))

	// Initial column width based on column names
	for i, v := range columnNames {
		columnWidths[i] = len(v)
	}

	// Get the width of each columns based on field's value
	for _, instance := range list {
		v := reflect.ValueOf(instance)
		for i := 0; i < len(columnNames); i++ {
			f := v.FieldByName(columnNames[i]).String()
			if len(f) > columnWidths[i] {
				columnWidths[i] = len(f)
			}
		}
	}

	return columnWidths
}

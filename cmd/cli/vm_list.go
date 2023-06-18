package cli

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/marintailor/rcstate/cmd/api/vm"
)

// list returns a table formatted list of instances.
func vmList(args []string) int {
	v, err := vm.NewVirtualMachine(args)
	if err != nil {
		fmt.Println("list: new environment:", err)
		return 1
	}

	if err := v.Cfg.ParseFlags(args); err != nil {
		fmt.Println("list: parse flags:", err)
		return 1
	}

	if err := v.Instances.GetInstancesList(); err != nil {
		fmt.Println("list: get instances list:", err)
		return 1
	}

	fmt.Println("LIST")
	fmt.Println(v.Cfg)

	columnNames := []string{"Name", "Status", "Internal", "External", "Type", "Preemptible"}
	columnWidth := getColumnsWidth(columnNames, v)

	var out bytes.Buffer

	tableHeader(&out, v, columnWidth)
	tableColumns(&out, columnNames, columnWidth)
	tableRows(&out, v, columnNames, columnWidth)

	fmt.Println(out.String())

	return 0
}

// tableHeader writes to a writer the header for the table.
func tableHeader(w io.Writer, v *vm.VirtualMachine, columnWidths []int) {
	line := ""
	for _, w := range columnWidths {
		line = line + strings.Repeat("=", w+2)
	}

	fmt.Fprintf(w, "\n%s\nPROJECT: %s    ZONE: %s\n%s\n\n", line, v.Cfg.Project, v.Cfg.Zone, line)
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
func tableRows(w io.Writer, v *vm.VirtualMachine, columnNames []string, columnWidths []int) {
	fmt.Println("v.Instances.List=", v.Instances.List)
	for _, inst := range v.Instances.List {
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
func getColumnsWidth(columnNames []string, v *vm.VirtualMachine) []int {
	columnWidths := make([]int, len(columnNames))

	// Initial column width based on column names
	for i, v := range columnNames {
		columnWidths[i] = len(v)
	}

	// Get the width of each columns based on field's value
	for _, instance := range v.Instances.List {
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

package awsconfig

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws/endpoints"
)

func ParseRegion(i string) (interface{}, error) {
	i = strings.TrimSpace(i)
	if !IsValidRegion(i) {
		return i, fmt.Errorf("'%s' is not a valid region", i)
	}
	return i, nil
}

func WarningChangeRegion(i interface{}) {
	region := fmt.Sprint(i)
	fmt.Fprintf(os.Stderr, "Region updated to '%s'.\nYou might want to update your default AMI with `awless config set instance.image $(awless search images amazonlinux --id-only --silent)`\n", region)
}

func ParseInstanceType(i string) (interface{}, error) {
	if !isValidInstanceType(i) {
		return i, fmt.Errorf("'%s' is not a valid instance type", i)
	}
	return i, nil
}

func StdinRegionSelector() string {
	fmt.Println("Please choose one region:")
	var region string

	regions := allRegions()
	fmt.Println(strings.Join(regions, ", "))
	fmt.Println()
	fmt.Print("Value ? > ")
	fmt.Scan(&region)
	for !isValidRegion(region, regions) {
		fmt.Printf("'%s' is not a valid region\n", region)
		fmt.Print("Value ? > ")
		fmt.Scan(&region)
	}
	return region
}

func StdinInstanceTypeSelector() string {
	fmt.Println("Please choose one instance type")
	fmt.Println()
	fmt.Println("Here are few examples:")

	var instanceType string
	t := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(t, "\tinstance type\tvCPU\tMemory (GiB)")
	fmt.Fprintln(t, "\tt2.nano\t1\t0.5")
	fmt.Fprintln(t, "\tt2.micro\t1\t1")
	fmt.Fprintln(t, "\tt2.small\t1\t2")
	fmt.Fprintln(t, "\tt2.medium\t2\t4")
	fmt.Fprintln(t, "\tt2.large\t2\t8")
	fmt.Fprintln(t, "\tt2.xlarge\t4\t16")
	fmt.Fprintln(t, "\tt2.2xlarge\t8\t32")
	fmt.Fprintln(t, "\tm4.large\t2\t8")
	fmt.Fprintln(t, "\tm4.xlarge\t4\t16")
	fmt.Fprintln(t, "\tc4.large\t2\t3.75")
	fmt.Fprintln(t, "\tc4.xlarge\t4\t7.5")
	fmt.Fprintln(t, "\t...")
	t.Flush()

	fmt.Println()
	fmt.Print("Value ? > ")
	fmt.Scan(&instanceType)
	for !isValidInstanceType(instanceType) {
		fmt.Printf("'%s' is not a valid instance type\n", instanceType)
		fmt.Print("Value ? > ")
		fmt.Scan(&instanceType)
	}
	return instanceType
}

func IsValidRegion(given string) bool {
	return isValidRegion(given, allRegions())
}

func isValidRegion(given string, regions []string) bool {
	for _, r := range regions {
		if given == r {
			return true
		}
	}
	return false
}

func isValidInstanceType(given string) bool {
	return regexp.MustCompile("\\w+\\.\\w+").MatchString(given)
}

func allRegions() []string {
	var regions sort.StringSlice
	partitions := endpoints.DefaultResolver().(endpoints.EnumPartitions).Partitions()
	for _, p := range partitions {
		for id := range p.Regions() {
			regions = append(regions, id)
		}
	}
	sort.Sort(regions)
	return regions
}

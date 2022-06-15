package homework8

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	errorNoOpFlag          = errors.New("-operation flag has to be specified")
	errorNoFilenameFlag    = errors.New("-fileName flag has to be specified")
	errorNoItemFlag        = errors.New("-item flag has to be specified")
	errorNoIdFlag          = errors.New("-id flag has to be specified")
	errorThereAreNoCircles = errors.New("There are no circles in box")
)

var (
	fn = flag.String("fileName", "", "Filename to store")
	op = flag.String("operation", "", "Operation: add|remove|list|findById")
	it = flag.String("item", "", "Json string with item")
	id = flag.String("id", "", "Item id")
)

type Arguments map[string]string

type Item struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func parseArgs() Arguments {
	flag.Parse()

	return Arguments{
		"id":        *id,
		"operation": *op,
		"item":      *it,
		"fileName":  *fn,
	}
}

func readJsonItems(filename string) (items []Item, err error) {
	_, err1 := os.Stat(filename)

	// assume no stat - no file to read - nothing to return
	if err1 != nil {
		return
	}

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return items, err
	}

	err = json.Unmarshal(bytes, &items)
	if err != nil {
		return items, err
	}

	return
}

func writeJsonItems(filename string, items []Item) error {
	bytes, err := json.Marshal(items)

	if err != nil {
		return err
	}

	return os.WriteFile(filename, bytes, 0644)
}

func Perform(args Arguments, writer io.Writer) error {
	operation := args["operation"]

	if operation == "" {
		return errorNoOpFlag
	}

	filename := args["fileName"]

	if filename == "" {
		return errorNoFilenameFlag
	}

	switch operation {
	case "list":
		items, err := readJsonItems(filename)

		if err != nil {
			panic(err)
		}

		bytes, err := json.Marshal(items)

		if err != nil {
			panic(err)
		}

		written, err := writer.Write(bytes)
		if len(bytes) != written {
			panic(err)
		}
	case "add":
		if args["item"] == "" {
			return errorNoItemFlag
		}

		var item Item

		err := json.Unmarshal([]byte(args["item"]), &item)
		if err != nil {
			panic(err)
		}

		items, err := readJsonItems(filename)
		if err != nil {
			panic(err)
		}

		for _, value := range items {
			if item.Age == value.Age && item.Email == value.Email && item.Id == value.Id {
				fmt.Fprintf(writer, "Item with id %s already exists", item.Id)
				return nil
			}
		}

		err = writeJsonItems(filename, append(items, item))

		if err != nil {
			panic(err)
		}
	case "remove":
		if args["id"] == "" {
			return errorNoIdFlag
		}

		items, err := readJsonItems(filename)
		if err != nil {
			panic(err)
		}

		for i, value := range items {
			if args["id"] == value.Id {
				return writeJsonItems(filename, append(items[:i], items[i+1:]...))
			}
		}

		fmt.Fprintf(writer, "Item with id %s not found", args["id"])
	case "findById":
		if args["id"] == "" {
			return errorNoIdFlag
		}

		items, err := readJsonItems(filename)
		if err != nil {
			panic(err)
		}

		for _, value := range items {
			if args["id"] == value.Id {
				bytes, err := json.Marshal(value)

				if err != nil {
					panic(err)
				}

				written, err := writer.Write(bytes)
				if len(bytes) != written {
					panic(err)
				}
			}
		}
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}

	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

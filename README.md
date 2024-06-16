# jsonschematics

`jsonschematics` is a Go package designed to validate and manipulate JSON data structures using schematics.

## Features

- Validate JSON objects against defined schematics
- Convert schematics to JSON schemas
- Handle complex data validation scenarios
- Perform Operations on the Data

## Installation

To install the package, use the following command:

```sh
go get github.com/ashbeelghouri/jsonschematics
```

## Usage

### Validation

#### Validating JSON Data
You can validate JSON data against a defined schematic using the Validate function. Here's an example:

```golang
package main

import (
    "fmt"
    "github.com/ashbeelghouri/jsonschematics"
)

func main() {
    schema := jsonschematics.Schematics{
        // Define your schema here
    }

    data := map[string]interface{}{
        "Name": "John",
        "Age":  30,
    }

    err := schema.Validate(data)
    if err != nil {
        fmt.Println("Validation errors:", err)
    } else {
        fmt.Println("Validation successful")
    }
}
```

#### Loading Schematics From JSON file
Instead of defining the Schema directly, Load the schema from JSON file:

```golang
package main

import (
    "fmt"
    "github.com/ashbeelghouri/jsonschematics"
)

func main() {
    schematics, err := jsonschematics.LoadFromJsonFile("path-to-your-schema.json")
    if err != nil {
        fmt.Println("Unable to load the schema:", err)
    }else {
        fmt.Println("Schema Loaded Successfully")
    }
}
```
see the API Reference for json fields mapping.


#### Loading Schematics From map[string]interface{}
If you want to load the schema from map[string]interface, you can use the below example:

```golang
package main

import (
    "fmt"
    "github.com/ashbeelghouri/jsonschematics"
)

func main() {
    schema := map[string]interface{}{
        ... define your schema
    }

    schematics, err := jsonschematics.LoadFromMap(&schema)
    if err != nil {
        fmt.Println("Unable to load the schema:", err)
    }else {
        fmt.Println("Schema Loaded Successfully")
    }
}
```

#### Adding Custom Validation Functions
You can also add your functions to validate the data:

##### Example 1
```golang
package main

import (
    "fmt"
    "github.com/ashbeelghouri/jsonschematics"
)

func main() {
    schematics, err := jsonschematics.LoadFromJsonFile("path-to-your-schema.json")
    if err != nil {
        fmt.Println("Unable to load the schema:", err)
    }
    schema.Validators.RegisterValidator("StringIsInsideArr", StringInArr)
}

func StringInArr(i interface{}, attr map[string]interface{}) error {
	str := i.(string)
	strArr := attr["arr"].([]string)
	found := false
	if len(strArr) > 0 {
		for _, item := range strArr {
			if item == str {
				found = true
			}
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("string not found in array"))
	}
	return nil
}

```

##### Example 2
```golang
package main

import (
    "fmt"
    "github.com/ashbeelghouri/jsonschematics"
)

func main() {
    schematics, err := jsonschematics.LoadFromJsonFile("path-to-your-schema.json")
    if err != nil {
        fmt.Println("Unable to load the schema:", err)
    }
    schema.Validators.RegisterValidator("StringIsInsideArr", func(i interface{}, attr map[string]interface{}) error {
    	str := i.(string)
    	strArr := attr["arr"].([]string)
    	found := false
    	if len(strArr) > 0 {
    		for _, item := range strArr {
    			if item == str {
    				found = true
    			}
    		}
    	}
    	if !found {
    		return errors.New(fmt.Sprintf("string not found in array"))
    	}
    	return nil
    })
}

```

### Operations

#### Perform Operations on Object
```golang
package main

import (
    "fmt"
    "github.com/ashbeelghouri/jsonschematics"
)

func main() {
    schema := jsonschematics.Schematics{
        // Define your schema here
    }

    data := map[string]interface{}{
        "Name": "John",
        "Age":  30,
    }

    newData := schema.PerformOperations(data)
    fmt.Printf("Data after Operations: %v", newData)
}
```

#### Adding Custom Operator Functions
You can also add your functions to operate on the data:
```golang
package main

import (
    "fmt"
    "github.com/ashbeelghouri/jsonschematics"
)

func main() {
    schematics, err := jsonschematics.LoadFromJsonFile("path-to-your-schema.json")
    if err != nil {
        fmt.Println("Unable to load the schema:", err)
    }
    schema.Operators.RegisterOperation("CapitalizeString", Capitalize)
}

func Capitalize(i interface{}, attributes map[string]interface{}) *interface{} {
	str := i.(string)
	var opResult interface{} = strings.ToUpper(string(str[0])) + strings.ToLower(str[1:])
	return &opResult
}
```

## API Reference

### Example Files
- [Schema](https://github.com/ashbeelghouri/jsonschematics/blob/master/json/schema.json)
- [Data](https://github.com/ashbeelghouri/jsonschematics/blob/master/json/data.json)

### Structs

#### Schematics
```golang
- Schema                                       		Schema
- Validators                                   		validators.Validators
- Operators				       		operators.Operators
- Prefix                                       		string
- Separator                                    		string
- ArrayIdKey                                   		string
- LoadSchema(filePath string)                  		error
- Validate(data map[string]interface{})        		*ErrorMessages
- ValidateArray(data []map[string]interface{}) 		*[]ArrayOfErrors
- PerformOperations(data map[string]interface{})	*map[string]interface{}
- PerformArrOperations(data []map[string]interface{})	*[]map[string]interface{}
- MakeFlat(data map[string]interface)          		*map[string]interface{}
```

##### Schema
```golang
- Version string `json:"version"`
- Fields []Field `json:"fields"`
```

###### >Explanation
```sh
* Version is for the maintenance of the schema
* Fields contains the validation logic for all the keys
```

##### Field
```golang
- DependsOn   []string 			`json:"depends_on"`
- TargetKey   string 			`json:"target_key"`
- Description string 			`json:"description"`
- Validators  []string 			`json:"validators"`
- Constants   map[string]interface{} 	`json:"constants"`
- Operators   []string            	`json:"operators"`
```

###### >Explanation
```sh
* DependsOn will check if the keys in array exists in data
* TargetKey will target the value in the data throught the key
* Description can have anything to explain the data, this can also be empty
* Validators is an array of string "validation functions"
* Constants will have dynanmic constants for each validator
```

##### Constant
```golang
- Attributes map[string]interface{} 	`json:"attributes"`
- ErrMsg     string 			`json:"err"`
```
###### >Explanation
```sh
* Attributes are passed in to the validation function so, it can have any map string interface.
* ErrMsg is a string that is shown as an error when validation fails
```

#### Errors
- ArrayOfErrors
- ErrorMessages
- ErrorMessage

##### ArrayOfErrors
```golang
- Errors ErrorMessages
- ID     interface{}
```

##### ErrorMessages
```golang
- Messages                                                 []ErrorMessage
- AddError(validator string, target string, err string)
- HaveErrors()                                             bool
```

##### ErrorMessage
```golang
- Message   string
- Validator string
- Target    string
```

#### Go Version
```golang
go 1.22.1
```

## Contributing
1. Fork the repository on GitHub.
2. Create a new branch for your feature or bug fix.
3. Write tests to cover your changes.
4. Update Documentation to include your features/changes.
5. Add yourself to contributers.
6. Send a pull request.

### Contributers
<a href="https://github.com/ashbeelghouri">
  <img src="https://avatars.githubusercontent.com/u/41609537?s=400&u=1b9ea072fc9a11acf32d86c5196a08f2696a458a&v=4" width="50px" height: "50px"/>
</a>

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/ashbeelghouri/jsonschematics/blob/master/LICENSE) file for details.

package rpcshared

import (
    "fmt"
    "os/exec"
	"os"
    "log"
    "bytes"
	"io/ioutil"
	"encoding/json"
)

type OpenNSFW string

type Args struct {
	Data   []byte
	DataID string
}

// Reference: https://github.com/yahoo/open_nsfw
func (t *OpenNSFW) Evaluate(args *Args, reply *string) error {
	pathToTool := "/opennsfw/classify_nsfw.py"
	fmt.Println("Path to Classify_NSFW.py: " , pathToTool)

	ON_Input_Directory := "/tmp/" + args.DataID + "/"
	// ClassifyNSFW works by a asking a trained model whether or not it scores a high value on an 'NSFW' test
	merr := os.MkdirAll(ON_Input_Directory, 0644)
	if merr != nil {
		log.Println("Error creating directory.", merr)
	}

	filepath := ON_Input_Directory + "data.dat"
	fmt.Println("Path to write given data to: ", filepath)
	err := ioutil.WriteFile(filepath, args.Data, 0644)
	if err != nil {
		log.Println("Error writing input to directory.", err)
	}

	//Setup the shell command to launch Bulk Extractor
	opts := []string{"--model_def", "/opennsfw/nsfw_model/deploy.prototxt", "--pretrained_model", "/opennsfw/nsfw_model/resnet_50_1by2_nsfw.caffemodel", filepath}
	//Should look like the following: 
	/* 
		cd open_nsfw
		docker run --volume=$(pwd):/workspace caffe:cpu \
		python ./classify_nsfw.py \
		--model_def nsfw_model/deploy.prototxt \
		--pretrained_model nsfw_model/resnet_50_1by2_nsfw.caffemodel \
		test_image.jpg
	*/
	fmt.Println("Command:\n", pathToTool, " ", opts)
	cmd := exec.Command(pathToTool, opts...)


	//Capture STDOUT
	var out bytes.Buffer
	cmd.Stdout = &out

	// Actually run the command: 
    err = cmd.Run()
	fmt.Println("[-] Output: ", out.String())

	//Post process the BE output
	jsonMapping := make(map[string]string)
	jsonKey := "NSFW-Score"
	jsonMapping[jsonKey] = out.String()

	// Dump everything into JSON in preperation for Elasticsearch upload
	jsonString, err := json.Marshal(jsonMapping)
	if err != nil {
            log.Println(err)
	}
	// Print raw json
	// fmt.Println(string(jsonString))   

	//We want to return the JSON in addition to STDOUT
	*reply = out.String()
	if err != nil {
            log.Println(err)
	}

	// If all goes well, remove temp directory
	remerr := os.RemoveAll(ON_Input_Directory)
	if remerr != nil {
		fmt.Println("Error cleaning up temporary directory: ", remerr)
	}
    return nil
}


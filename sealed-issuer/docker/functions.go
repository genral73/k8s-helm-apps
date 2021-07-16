package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

// Seal takes a sealedrequest object and returns a sealedresponse
func Seal(request SealedRequest) SealedResponse {
	var secretString string
	if request.KVPairs != nil {
		secretString = secretStringFromKV(request.Name, request.Namespace, request.KVPairs)
		// log.Println(secretString)
	} else if request.YAML != "" {
		var err error
		secretString, err = secretStringFromYAML(request.YAML)
		if err != nil {
			return SealedResponse{"", "", err.Error()}
		}
	} else {
		err := errors.New("Cannot determine input format")
		log.Println(err)
		return SealedResponse{"", "", err.Error()}
	}

	// execute kubeseal while piping the secretJSONString
	cmd := exec.Command("kubeseal/bin", "--cert", "kubeseal/cert.pem", "--format", "yaml")
	cmd.Stdin = strings.NewReader(secretString)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if stderrString := stderr.String(); stderrString != "" {
		log.Println(stderrString)
		return SealedResponse{"", "", stderrString}
	} else if err != nil {
		log.Println(err)
		return SealedResponse{"", "", err.Error()}
	} else {
		return SealedResponse{stdout.String(), fmt.Sprintf("%v.yaml", request.Name), ""}
	}
}

// secretString creates a secret and returns its encoded json
func secretStringFromKV(name, namespace string, pairs []KVPairs) string {
	data := make(map[string][]byte)
	for _, pair := range pairs {
		data[pair.K] = []byte(pair.V)
	}
	secret := apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
	return convertSecretToString(secret)
}

func secretStringFromYAML(yaml string) (string, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	secret := apiv1.Secret{}
	_, _, err := decode([]byte(yaml), nil, &secret)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return convertSecretToString(secret), nil
}

func convertSecretToString(secret apiv1.Secret) string {
	// convert to json
	j := kjson.NewSerializer(kjson.DefaultMetaFactory, nil, nil, false)
	buf := new(bytes.Buffer)
	if err := j.Encode(&secret, buf); err != nil {
		log.Println("Error encoding string to JSON")
		panic("Error encoding string to JSON")
	}

	// convert to map then json again to add missing fields
	var obj map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &obj); err != nil {
		log.Println("Error unmarshaling")
		panic("Error unmarshaling")
	}
	obj["apiVersion"] = "v1"
	obj["kind"] = "Secret"
	x, _ := json.Marshal(obj)

	return string(x)
}

// @author       sigu-399
// @description  An implementation of JSON Reference - Go language
// @created      26-02-2013

package gojsonschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"gojsonreference"
	"io/ioutil"
	"net/http"
)

type SchemaPool struct {
	schemaPoolDocuments map[string]*SchemaPoolDocument
}

func NewSchemaPool() *SchemaPool {
	p := &SchemaPool{}
	p.schemaPoolDocuments = make(map[string]*SchemaPoolDocument)
	return p
}

func (p *SchemaPool) GetPoolDocument(reference gojsonreference.JsonReference) (*SchemaPoolDocument, error) {

	var err error

	if !reference.HasFullUrl {
		return nil, errors.New(fmt.Sprintf("Reference must be canonical %s", reference))
	}

	refToUrl := reference
	refToUrl.GetUrl().Fragment = ""

	var spd *SchemaPoolDocument

	for k := range p.schemaPoolDocuments {
		if k == refToUrl.String() {
			spd = p.schemaPoolDocuments[k]
			fmt.Printf("Found in pool %s\n", refToUrl.String())
		}
	}

	if spd != nil {
		return spd, nil
	}

	document, err := getHttpJson(refToUrl.String())
	if err != nil {
		return nil, err
	}

	spd = &SchemaPoolDocument{Document: document}
	p.schemaPoolDocuments[refToUrl.String()] = spd

	fmt.Printf("Added to pool %s\n", refToUrl.String())

	return spd, nil
}

type SchemaPoolDocument struct {
	Document interface{}
}

func getHttpJson(url string) (interface{}, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Could not access schema " + resp.Status)
	}

	bodyBuff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var document interface{}
	err = json.Unmarshal(bodyBuff, &document)
	if err != nil {
		return nil, err
	}

	return document, nil
}

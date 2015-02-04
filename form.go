package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type BaseForm interface {
	FullName() string
	SetSignature(string) error
	IsArtist() bool
	IsModel() bool
}

type Form struct {
	FirstName string
	LastName  string
	Email     string
	Link      string
	Sig       []map[string]int
	EmailSent bool
}

func (f *Form) FullName() string {
	return fmt.Sprintf("%s %s", f.FirstName, f.LastName)
}

func (f *Form) SetSignature(sigString string) error {
	rawSig := []byte(sigString)
	var sig []map[string]int
	err := json.Unmarshal(rawSig, &sig)
	f.Sig = sig
	return err
}

type ArtistForm struct {
	Form
	Works []Work
}

func (a *ArtistForm) IsArtist() bool {
	return true
}

func (a *ArtistForm) IsModel() bool {
	return false
}

type Model struct {
	Form
}

type Work struct {
	Name        string
	Description string
	Photos      []Photo
}

type Photo struct {
	Name   string
	Models []Model
}

func getItemCount(filter string, form map[string][]string) int {
	count := 0
	for k, _ := range form {
		if strings.Contains(k, filter) {
			count++
		}
	}
	return count
}

/**
 * Returns an array of indices for works. This is needed b/c
 * it is possible for user to remove works.
 */
func getIndices(filter string, form map[string][]string) []int {
	numItems := getItemCount(filter, form)
	indices := make([]int, numItems)
	i := 0
	for k, _ := range form {
		if strings.Contains(k, filter) {
			j, err := strconv.Atoi(k[len(k)-1:])
			if err != nil {
				panic(err)
			}

			indices[i] = j
			i++
		}
	}
	return indices
}

func extractWorks(form map[string][]string) []Work {
	numItems := getItemCount("descOfWork", form)
	workIndices := getIndices("descOfWork", form)
	works := make([]Work, numItems)

	for i, e := range workIndices {
		works[i] = Work{
			Name:        form[fmt.Sprintf("nameOfWork%d", e)][0],
			Description: form[fmt.Sprintf("descOfWork%d", e)][0],
		}
	}
	return works
}

package main

import (
	"encoding/json"
	"errors"
	"f"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type GithubProfile struct {
	Login    string `json: login`
	Name     string `json: name`
	company  string `json: company`
	location string `json: location`
}

func githubTask(in *f.Input) (*f.Output, error) {

	client := &http.Client{
		Timeout: time.Second * 2,
	}
	user, ok := in.Data.(string)
	if !ok {
		err := errors.New("Expecting user as a string value ")
		return nil, err
	}
	req, err := http.NewRequest(
		http.MethodGet, "http://api.github.com/users/"+user, nil)
	if err != nil {
		err := errors.New("Error getting github profile for user " + user)
		return nil, err

	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error is ", err)
		//log.Fatal("Error getting github profile for user", user)
		err := errors.New("Error getting github profile for the user " + user)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err := errors.New("Error parsing github profile for user " + user)
		return nil, err
	}
	profile := GithubProfile{}
	jsonErr := json.Unmarshal(body, &profile)
	if jsonErr != nil {
		//log.Fatal(jsonErr)
		return nil, jsonErr
	}
	//fmt.Println(profile.Name)
	out := &f.Output{}
	out.Data = profile
	return out, nil
}

type GoogleProfile struct {
	Email string
	Name  string
}

func getGoogleProfile() *GoogleProfile {
	return &GoogleProfile{
		Name:  "Last, First",
		Email: "name@gmail.com",
	}
}

func toInput(value string) *f.Input {
	in := &f.Input{
		Data: value,
	}
	return in
}

func toOutput(value interface{}) *f.Output {
	return &f.Output{
		Data: value,
	}
}
func main() {
	flow := f.New("sample")
	s1 := f.NewState("state-1")
	s1.SetConcurrent(true)
	s1.AddTasks(
		f.Task{
			Name: "Task-1: Get user profile from Github",
			Execute: func(in *f.Input) (*f.Output, error) {
				time.Sleep(1 * time.Second)
				return githubTask(in)
			},
		},
		f.Task{
			Name: "Task-2: Do someting with user profile2",
			Execute: func(in *f.Input) (*f.Output, error) {
				time.Sleep(3 * time.Second)
				//out := toOutput(getGoogleProfile())
				return githubTask(in)
				//return out, nil
			},
		},
		f.Task{
			Name: "Task-3: Do someting with user profile3",
			Execute: func(in *f.Input) (*f.Output, error) {
				//out := toOutput(getGoogleProfile())
				return githubTask(in)
				//return out, nil
			},
		},
	)
	s2 := f.NewState("state-2")
	s2.AddTasks(
		f.Task{
			Name: "Task-1: Do something else",
			Execute: func(in *f.Input) (*f.Output, error) {
				fmt.Println("This is task 1 of State 2")
				out := &f.Output{}
				return out, nil

			},
		},
	)
	s3 := f.NewState("state-3")
	s4 := f.NewState("state-4")
	flow.AddStates(s1, s2, s3, s4)
	flow.Execute(toInput("deepakjacob"))
}

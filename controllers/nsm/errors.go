package controllers

import "errors"

var ErrNsmAlreadyExists = errors.New("another instance of NSM is already running. Please delete that instance before creating a new one")

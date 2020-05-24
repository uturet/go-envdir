#!/bin/bash
env
while (( "$#" )); do
  echo $1
  shift
done
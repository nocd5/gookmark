#!/bin/bash

# [.bashrc]
# source <THIS_SCRIPT>

function pgm() {
  target=`gookmark list | peco`
  if [ ${#target} != 0 ]; then
    if [ -d ${target} ]; then
      pushd ${target}
    else
      gnome-open ${target}
    fi
  fi
}

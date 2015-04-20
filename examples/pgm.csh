#!/bin/csh

# [.cshrc]
# alias pgm source <THIS_SCRIPT>

set target = `gookmark list | peco`
if (${#target} != 0) then
  if ( -d ${target} ) then
    pushd ${target}
  else
    gnome-open ${target}
  endif
endif

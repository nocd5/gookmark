-- [.nyagos]
-- local chank,err = assert(loadfile("<THIS_SCRIPT>"))
-- if err then
--     print(err)
-- else
--     chank()
-- end

alias {
  pgm=function(args)
    local target = nyagos.eval('gookmark list | peco')
    if #target ~= 0 then
      if is_dir(target) then
        nyagos.exec('pushd ' .. target)
      else
        nyagos.exec('open ' .. target)
      end
    end
  end
}

function is_dir(f)
  return nyagos.eval('file ' .. f):find(': directory')
end

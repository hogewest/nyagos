nyagos.key.C_o = function(this)
    local word,pos = this:lastword()
    word = string.gsub(word,'"','')
    local wildcard = word.."*"
    local list = nyagos.glob(wildcard)
    if #list == 1 and list[1] == wildcard then
        return
    end
    nyagos.write("\n")
    local result=nyagos.box(list)
    this:call("REPAINT_ON_NEWLINE")
    if string.find(result," ",1,true) then
        if string.find(result,"^~[\\/]") then
            result = '~"'..string.sub(result,2)..'"'
        else
            result = '"'..result..'"'
        end
    end
    assert( this:replacefrom(pos,result) )
end

nyagos.alias.__dump_history = function()
    local uniq={}
    for i=nyagos.gethistory()-1,1,-1 do
        local line = nyagos.gethistory(i)
        if line ~= "" and not uniq[line] then
            nyagos.write(line,"\n")
            uniq[line] = true
        end
    end
end

nyagos.key.C_x = function(this)
    nyagos.write("\nC-x: [r]:command-history, [h]:cd-history, [g]:git-revision\n")
    local ch = nyagos.getkey()
    local c = string.lower(string.char(ch))
    local result
    if c == 'r' or ch == (string.byte('r') & 0x1F) then
        result = nyagos.eval('__dump_history | box')
    elseif ch == 'h' or ch == (string.byte('h') & 0x1F) then
        result = nyagos.eval('cd --history | box')
        if string.find(result,' ') then
            result = '"'..result..'"'
        end
    elseif c == 'g' or ch == (string.byte('g') & 0x1F) then
        result = nyagos.eval('git log --pretty="format:%h %s" | box')
        result = string.match(result,"^%S+") or ""
    end
    this:call("REPAINT_ON_NEWLINE")
    return result
end

nyagos.key.M_r = function(this)
    nyagos.write("\n")
    local result = nyagos.eval('__dump_history | box')
    this:call("REPAINT_ON_NEWLINE")
    if string.find(result,' ') then
        result = '"'..result..'"'
    end
    return result
end

nyagos.key.M_h = function(this)
    nyagos.write("\n")
    local result = nyagos.eval('cd --history | box')
    this:call("REPAINT_ON_NEWLINE")
    if string.find(result,' ') then
        result = '"'..result..'"'
    end
    return result
end

nyagos.key.M_g = function(this)
    nyagos.write("\n")
    local result = nyagos.eval('git log --pretty="format:%h %s" | box')
    this:call("REPAINT_ON_NEWLINE")
    return string.match(result,"^%S+") or ""
end

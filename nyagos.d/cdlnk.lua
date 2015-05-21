local wshObj = nyagos.ole.create_object_utf8("WScript.Shell")

local function getfolder(arg1)
    local shortcut1 = wshObj:CreateShortcut(arg1)
    if not shortcut1 then
        return arg1
    end
    local path = shortcut1.TargetPath
    if nyagos.fsObj:FolderExists(path) then
        return path
    end
    path = shortcut1.WorkingDirectory
    if nyagos.fsObj:FolderExists(path) then
        return path
    end
    path = nyagos.fsObj:GetParentFolderName(shortcut1.TargetPath)
    if nyagos.fsObj:FolderExists(path) then
        return path
    end
    return arg1
end

local org_cd = nyagos.alias.cd
nyagos.alias.cd=function(args)
    for i=1,#args do
        local arg1 = args[i]
        if string.match(arg1,"%.[lL][nN][kK]$") then
            arg1 = getfolder(arg1)
        end
        args[i] = arg1
    end
    if org_cd then
        org_cd(args)
    else
        args[0] = "__cd__"
        local status,errmsg = nyagos.exec(args)
        if not status then
            nyagos.writerr(errmsg)
        end
    end
end

NYAGOS 4.0.4\_0
================
* Support Lua 5.3

NYAGOS 4.0.3\_2
===============
* Command-name completion supports alias and built-in commands.

BugFix
------
* `pwd` did not print correct UNC-Path (#37)
* `nyagos.gethistory( large-value )` crashes nyagos (#38 @1)
* `%APPDATA%/NYAOS_ORG/nyagos.history` did not be updated. (#39 @1)
* Crashed when `%APPDATA%/NYAOS_ORG/nyagos.history` did not exists or is empty.x (#40 @1)
* On French keyboard, keys shifted with AltGr could not be input. (#41)

@1 Probably these bugs appeared only on snapshot-build.
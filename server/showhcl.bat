@set SCRIPT_DIR=%~dp0
pushd %~dp1
@java -jar "%SCRIPT_DIR%drawj2d.jar" -T screen -W 297 -H 297 %*
popd

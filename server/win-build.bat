@ECHO OFF
echo builing...
cd main\link 
go build -o ..\..\..\deploy\link.exe 
cd ..\..\
cd main\admin 
go build -o ..\..\..\deploy\admin.exe 
cd ..\..\
cd main\app 
go build -o ..\..\..\deploy\app.exe 
cd ..\..\
cd main\tools 
go build -o ..\..\..\deploy\tools.exe 
cd ..\..\
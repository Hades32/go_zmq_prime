numWorkers=5
if [ -z "$1" ]
then
	echo using default workers
else
	numWorkers=$1
fi

echo using $numWorkers workers
for (( c=1; c<=$numWorkers; c++ ))
do
	bin/taskwork.exe & 
done

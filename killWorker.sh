ps aux | grep -ie taskwork | awk '{print $1}' | xargs kill -9

currDate=`date`
echo "$currDate"

git -C $GOPATH/src/habitica/hackintasks.github.io add .
git -C $GOPATH/src/habitica/hackintasks.github.io commit -m "Update $currDate"
git -C $GOPATH/src/habitica/hackintasks.github.io push origin master

